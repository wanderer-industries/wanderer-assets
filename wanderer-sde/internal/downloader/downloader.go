// Package downloader fetches the EVE SDE tables the converter needs from
// CCP's SDE API (https://sde.evetech.net).
//
// The API lists releases per tenant. Each release has an index with its
// release name, build number, release date and one JSON Lines file per
// table. Only the current and the previous release are kept, so table URLs
// are always read from a fresh index and never stored.
package downloader

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/guarzo/wanderer-sde/internal/config"
)

// Tables are the SDE tables the parser reads. Each is saved as <table>.jsonl.
var Tables = []string{
	"types",
	"groups",
	"categories",
	"mapRegions",
	"mapConstellations",
	"mapSolarSystems",
	"mapStars",
	"mapStargates",
	"npcStations",
	"npcCorporations",
}

// ExpectedFiles are the files a valid SDE directory holds.
var ExpectedFiles = func() []string {
	files := make([]string, len(Tables))
	for i, t := range Tables {
		files[i] = t + ".jsonl"
	}
	return files
}()

// Release is the index of one SDE release.
type Release struct {
	Release     string           `json:"release"`
	BuildNumber int64            `json:"buildNumber"`
	ReleaseDate string           `json:"releaseDate"`
	Tables      map[string]Table `json:"tables"`
}

// Table is one table of a release.
type Table struct {
	Version string `json:"version"`
	JSONL   struct {
		URL string `json:"url"`
	} `json:"jsonl"`
}

// VersionInfo returns the version of the release.
func (r *Release) VersionInfo() *VersionInfo {
	return &VersionInfo{
		Release:     r.Release,
		BuildNumber: strconv.FormatInt(r.BuildNumber, 10),
		ReleaseDate: r.ReleaseDate,
	}
}

// Downloader downloads SDE tables.
type Downloader struct {
	config     *config.Config
	httpClient *http.Client
}

// New creates a new Downloader with the given configuration.
func New(cfg *config.Config) *Downloader {
	return &Downloader{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Minute, // per table; types is the largest at ~25 MB gzipped
		},
	}
}

// releaseURL returns the index URL of a release, or of the latest one when
// release is empty.
func releaseURL(base, release string) string {
	if release == "" {
		release = "latest"
	}
	return strings.TrimSuffix(base, "/") + "/" + release + "/"
}

// FetchRelease reads the index of a release, such as "20261119T110708Z",
// or of the latest one when release is empty.
func (d *Downloader) FetchRelease(ctx context.Context, release string) (*Release, error) {
	var rel Release
	if _, err := getJSON(ctx, d.httpClient, d.config.Version, releaseURL(d.config.SDEUrl, release), &rel); err != nil {
		return nil, err
	}
	if rel.Release == "" || len(rel.Tables) == 0 {
		return nil, fmt.Errorf("not an SDE release index: %s", releaseURL(d.config.SDEUrl, release))
	}
	return &rel, nil
}

// Download saves every table in Tables from a release (the latest when
// release is empty) into a new temporary directory, and returns that
// directory with the release it came from.
func (d *Downloader) Download(ctx context.Context, release string) (string, *Release, error) {
	rel, err := d.FetchRelease(ctx, release)
	if err != nil {
		return "", nil, err
	}
	if d.config.Verbose {
		fmt.Printf("Downloading SDE release %s (build %d)\n", rel.Release, rel.BuildNumber)
	}

	dir, err := os.MkdirTemp("", "wanderer-sde-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	for _, name := range Tables {
		t, ok := rel.Tables[name]
		if !ok || t.JSONL.URL == "" {
			_ = os.RemoveAll(dir)
			return "", nil, fmt.Errorf("release %s has no JSONL table %q", rel.Release, name)
		}
		n, err := d.downloadFile(ctx, t.JSONL.URL, filepath.Join(dir, name+".jsonl"))
		if err != nil {
			_ = os.RemoveAll(dir)
			return "", nil, fmt.Errorf("failed to download %s: %w", name, err)
		}
		if d.config.Verbose {
			fmt.Printf("  %s.jsonl (%s)\n", name, formatBytes(n))
		}
	}
	return dir, rel, nil
}

func (d *Downloader) downloadFile(ctx context.Context, url, path string) (int64, error) {
	body, _, err := get(ctx, d.httpClient, d.config.Version, url)
	if err != nil {
		return 0, err
	}
	defer func() { _ = body.Close() }()

	out, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(out, body)
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	return n, err
}

// Validate checks that the SDE directory holds every expected table.
func (d *Downloader) Validate(sdePath string) error {
	if d.config.Verbose {
		fmt.Printf("Validating SDE structure at: %s\n", sdePath)
	}

	for _, file := range ExpectedFiles {
		checkPath := filepath.Join(sdePath, file)
		if _, err := os.Stat(checkPath); os.IsNotExist(err) {
			return fmt.Errorf("missing expected file: %s", file)
		}
	}

	if d.config.Verbose {
		fmt.Println("SDE validation successful")
	}
	return nil
}

func userAgent(version string) string {
	if version == "" {
		version = "dev"
	}
	return fmt.Sprintf("wanderer-sde/%s (https://github.com/guarzo/wanderer-sde)", version)
}

// get follows redirects (table URLs point at content-addressed blobs) and
// undoes the gzip the API applies to every file. Go's transport only does
// that on its own when it asked for gzip itself.
func get(ctx context.Context, client *http.Client, version, url string) (io.ReadCloser, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent(version))

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("GET %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, nil, fmt.Errorf("GET %s: unexpected status code: %d", url, resp.StatusCode)
	}
	if !resp.Uncompressed && strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			return nil, nil, fmt.Errorf("GET %s: %w", url, err)
		}
		return gzipBody{zr, resp.Body}, resp.Header, nil
	}
	return resp.Body, resp.Header, nil
}

type gzipBody struct {
	*gzip.Reader
	body io.Closer
}

func (g gzipBody) Close() error { return g.body.Close() }

func getJSON(ctx context.Context, client *http.Client, version, url string, v any) (http.Header, error) {
	body, header, err := get(ctx, client, version, url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = body.Close() }()
	if err := json.NewDecoder(body).Decode(v); err != nil {
		return nil, fmt.Errorf("GET %s: failed to decode JSON: %w", url, err)
	}
	return header, nil
}

// formatBytes formats bytes into a human-readable string.
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
