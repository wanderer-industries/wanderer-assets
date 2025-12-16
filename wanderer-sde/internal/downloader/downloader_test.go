package downloader

import (
	"archive/zip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
)

func TestDownloader_Validate(t *testing.T) {
	// Create a valid SDE structure
	tmpDir, err := os.MkdirTemp("", "downloader_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create all expected files
	for _, filename := range ExpectedFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte("test"), 0644); err != nil {
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
	_ = os.Remove(filepath.Join(tmpDir, "types.yaml"))
	if err := dl.Validate(tmpDir); err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestDownloader_Extract(t *testing.T) {
	// Create a test ZIP file
	tmpDir, err := os.MkdirTemp("", "downloader_extract_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	zipPath := filepath.Join(tmpDir, "test.zip")
	extractDir := filepath.Join(tmpDir, "extracted")

	// Create ZIP file with test content
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create zip file: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	// Add some files to the ZIP
	testFiles := map[string]string{
		"types.yaml":             "type: test",
		"groups.yaml":            "group: test",
		"categories.yaml":        "category: test",
		"mapSolarSystems.yaml":   "system: test",
		"mapRegions.yaml":        "region: test",
		"mapConstellations.yaml": "constellation: test",
		"mapStargates.yaml":      "stargate: test",
		"subdir/nested.yaml":     "nested: test",
	}

	for name, content := range testFiles {
		w, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("failed to create file in zip: %v", err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write file content: %v", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
	_ = zipFile.Close()

	cfg := &config.Config{Verbose: false}
	dl := New(cfg)

	// Extract the ZIP
	sdePath, err := dl.Extract(zipPath, extractDir)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify files were extracted
	for name := range testFiles {
		path := filepath.Join(sdePath, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("File not extracted: %s", name)
		}
	}
}

func TestDownloader_ExtractZipSlipProtection(t *testing.T) {
	// Create a test ZIP file with path traversal attempt
	tmpDir, err := os.MkdirTemp("", "downloader_zipslip_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	zipPath := filepath.Join(tmpDir, "malicious.zip")
	extractDir := filepath.Join(tmpDir, "extracted")

	// Create ZIP file with path traversal
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create zip file: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	// Add a file with path traversal
	w, err := zipWriter.Create("../../../etc/passwd")
	if err != nil {
		t.Fatalf("failed to create file in zip: %v", err)
	}
	if _, err := w.Write([]byte("malicious content")); err != nil {
		t.Fatalf("failed to write file content: %v", err)
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
	_ = zipFile.Close()

	cfg := &config.Config{Verbose: false}
	dl := New(cfg)

	// Extract should fail due to zip slip protection
	_, err = dl.Extract(zipPath, extractDir)
	if err == nil {
		t.Error("Expected error for zip slip attack")
	}
}

func TestDownloader_Download(t *testing.T) {
	// Create a test server
	testContent := []byte("test zip content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(testContent)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testContent)
	}))
	defer server.Close()

	cfg := &config.Config{
		SDEUrl:  server.URL,
		Verbose: false,
	}
	dl := New(cfg)

	ctx := context.Background()
	result, err := dl.Download(ctx)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Clean up downloaded file
	defer func() { _ = os.RemoveAll(filepath.Dir(result.ZipPath)) }()

	// Verify file was downloaded
	if _, err := os.Stat(result.ZipPath); os.IsNotExist(err) {
		t.Error("Downloaded file does not exist")
	}

	// Verify content
	content, err := os.ReadFile(result.ZipPath)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("Downloaded content mismatch: got %q, want %q", string(content), string(testContent))
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
	_, err := dl.Download(ctx)
	if err == nil {
		t.Error("Expected error for failed download")
	}
}

func TestDownloader_DownloadContextCancellation(t *testing.T) {
	// Create a test server that delays
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if context was canceled
		select {
		case <-r.Context().Done():
			return
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("test"))
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		SDEUrl:  server.URL,
		Verbose: false,
	}
	dl := New(cfg)

	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := dl.Download(ctx)
	if err == nil {
		t.Error("Expected error for canceled context")
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
	// Verify all expected files are defined
	expectedFilesList := []string{
		"types.yaml",
		"groups.yaml",
		"categories.yaml",
		"mapSolarSystems.yaml",
		"mapRegions.yaml",
		"mapConstellations.yaml",
		"mapStargates.yaml",
	}

	for _, expected := range expectedFilesList {
		found := false
		for _, actual := range ExpectedFiles {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %q not in ExpectedFiles", expected)
		}
	}
}
