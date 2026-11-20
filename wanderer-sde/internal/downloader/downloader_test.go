package downloader

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
)

const testRelease = "20250115T110000Z"

// fakeSDEAPI serves releases like sde.evetech.net: latest/ redirects to the
// release, table URLs redirect to blobs, and every blob is gzipped whatever
// the client asked for. skip names tables left out of the index.
func fakeSDEAPI(t *testing.T, skip ...string) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tranquility/latest/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/tranquility/"+testRelease+"/", http.StatusFound)
	})
	mux.HandleFunc("GET /tranquility/"+testRelease+"/", func(w http.ResponseWriter, r *http.Request) {
		var tables []string
		for _, name := range Tables {
			if contains(skip, name) {
				continue
			}
			tables = append(tables, fmt.Sprintf(`%q: {"version": "v-%s", "jsonl": {"url": "%s/tranquility/%s/%s.jsonl"}}`,
				name, name, srv.URL, testRelease, name))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", `"index-etag"`)
		_, _ = fmt.Fprintf(w, `{"release": %q, "buildNumber": 2025001, "releaseDate": "2025-01-15T11:00:00Z", "tables": {%s}}`,
			testRelease, strings.Join(tables, ", "))
	})
	mux.HandleFunc("GET /tranquility/"+testRelease+"/{file}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/blobs/"+r.PathValue("file"), http.StatusFound)
	})
	mux.HandleFunc("GET /blobs/{file}", func(w http.ResponseWriter, r *http.Request) {
		table := strings.TrimSuffix(r.PathValue("file"), ".jsonl")
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, _ = fmt.Fprintf(zw, `{"id": 1, "table": %q}`+"\n", table)
		_ = zw.Close()
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(buf.Bytes())
	})
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func TestDownloader_Validate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create all expected files
	for _, filename := range ExpectedFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte("{}\n"), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", filename, err)
		}
	}

	cfg := &config.Config{Verbose: false}
	dl := New(cfg)

	// Test valid structure
	if err := dl.Validate(tmpDir); err != nil {
		t.Errorf("Validate failed on valid SDE: %v", err)
	}

	// Test missing file
	_ = os.Remove(filepath.Join(tmpDir, "types.jsonl"))
	if err := dl.Validate(tmpDir); err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestDownloader_Download(t *testing.T) {
	srv := fakeSDEAPI(t)

	for _, tt := range []struct {
		name    string
		release string
		client  *http.Client
	}{
		// Go's transport asks for gzip and decodes it on its own
		{"latest", "", nil},
		{"named release", testRelease, nil},
		// A client that does not ask for gzip still gets it
		{"gzip not requested", "", &http.Client{Transport: &http.Transport{DisableCompression: true}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{SDEUrl: srv.URL + "/tranquility/"}
			dl := New(cfg)
			if tt.client != nil {
				dl.httpClient = tt.client
			}

			dir, rel, err := dl.Download(context.Background(), tt.release)
			if err != nil {
				t.Fatalf("Download failed: %v", err)
			}
			defer func() { _ = os.RemoveAll(dir) }()

			if rel.Release != testRelease || rel.BuildNumber != 2025001 {
				t.Errorf("release = %+v", rel)
			}
			if err := dl.Validate(dir); err != nil {
				t.Errorf("Validate: %v", err)
			}
			for _, name := range Tables {
				got, err := os.ReadFile(filepath.Join(dir, name+".jsonl"))
				if err != nil {
					t.Fatalf("read %s: %v", name, err)
				}
				if want := fmt.Sprintf(`{"id": 1, "table": %q}`+"\n", name); string(got) != want {
					t.Errorf("%s.jsonl = %q, want %q", name, got, want)
				}
			}
		})
	}
}

func TestDownloader_DownloadMissingTable(t *testing.T) {
	srv := fakeSDEAPI(t, "npcStations")
	dl := New(&config.Config{SDEUrl: srv.URL + "/tranquility/"})

	_, _, err := dl.Download(context.Background(), "")
	if err == nil || !strings.Contains(err.Error(), `"npcStations"`) {
		t.Errorf("err = %v, want a missing npcStations table", err)
	}
}

func TestDownloader_DownloadError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := &config.Config{
		SDEUrl:  server.URL,
		Verbose: false,
	}
	dl := New(cfg)

	ctx := context.Background()
	_, _, err := dl.Download(ctx, "")
	if err == nil {
		t.Error("Expected error for failed download")
	}
}

func TestDownloader_NotAReleaseIndex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"releases": []}`))
	}))
	defer server.Close()

	dl := New(&config.Config{SDEUrl: server.URL})
	if _, err := dl.FetchRelease(context.Background(), ""); err == nil {
		t.Error("Expected error for a response that is not a release index")
	}
}

func TestDownloader_DownloadContextCancellation(t *testing.T) {
	srv := fakeSDEAPI(t)
	cfg := &config.Config{
		SDEUrl:  srv.URL + "/tranquility/",
		Verbose: false,
	}
	dl := New(cfg)

	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, _, err := dl.Download(ctx, "")
	if err == nil {
		t.Error("Expected error for canceled context")
	}
}

func TestRelease_VersionInfo(t *testing.T) {
	rel := &Release{Release: testRelease, BuildNumber: 2025001, ReleaseDate: "2025-01-15T11:00:00Z"}
	vi := rel.VersionInfo()
	if vi.Release != testRelease || vi.BuildNumber != "2025001" || vi.ReleaseDate != "2025-01-15T11:00:00Z" {
		t.Errorf("VersionInfo = %+v", vi)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1572864, "1.50 MB"},
		{1073741824, "1.00 GB"},
		{1610612736, "1.50 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestExpectedFiles(t *testing.T) {
	// Every table the parser reads must be downloaded and validated
	expectedFilesList := []string{
		"types.jsonl",
		"groups.jsonl",
		"categories.jsonl",
		"mapRegions.jsonl",
		"mapConstellations.jsonl",
		"mapSolarSystems.jsonl",
		"mapStars.jsonl",
		"mapStargates.jsonl",
		"npcStations.jsonl",
		"npcCorporations.jsonl",
	}

	if len(ExpectedFiles) != len(expectedFilesList) {
		t.Errorf("ExpectedFiles has %d files, want %d", len(ExpectedFiles), len(expectedFilesList))
	}
	for _, expected := range expectedFilesList {
		if !contains(ExpectedFiles, expected) {
			t.Errorf("Expected file %q not in ExpectedFiles", expected)
		}
	}
}
