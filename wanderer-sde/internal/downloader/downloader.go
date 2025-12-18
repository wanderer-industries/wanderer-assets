// Package downloader handles downloading and extracting the EVE SDE.
package downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/guarzo/wanderer-sde/internal/config"
)

// ExpectedFiles are the files we expect to find in a valid SDE.
// The new SDE format (2025+) uses a flat file structure.
var ExpectedFiles = []string{
	"types.yaml",
	"groups.yaml",
	"categories.yaml",
	"mapSolarSystems.yaml",
	"mapRegions.yaml",
	"mapConstellations.yaml",
	"mapStargates.yaml",
}

// Downloader handles downloading and extracting the SDE.
type Downloader struct {
	config     *config.Config
	httpClient *http.Client
}

// New creates a new Downloader with the given configuration.
func New(cfg *config.Config) *Downloader {
	return &Downloader{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Minute, // SDE is large, allow generous timeout
		},
	}
}

// DownloadResult contains information about a completed download.
type DownloadResult struct {
	ZipPath     string
	ExtractPath string
	SDEPath     string
	BytesRead   int64
}

// Download downloads the SDE from the configured URL.
// Returns the path to the downloaded ZIP file.
func (d *Downloader) Download(ctx context.Context) (*DownloadResult, error) {
	if d.config.Verbose {
		fmt.Printf("Downloading SDE from: %s\n", d.config.SDEUrl)
	}

	// Create a temporary directory for the download
	tempDir, err := os.MkdirTemp("", "wanderer-sde-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clean up temp directory on error; disabled on success
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(tempDir)
		}
	}()

	zipPath := filepath.Join(tempDir, "sde.zip")

	// Create the output file
	out, err := os.Create(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = out.Close() }()

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.config.SDEUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent header (required by some CDNs)
	version := d.config.Version
	if version == "" {
		version = "dev"
	}
	req.Header.Set("User-Agent", fmt.Sprintf("wanderer-sde/%s (https://github.com/guarzo/wanderer-sde)", version))
	req.Header.Set("Accept", "*/*")

	// Perform the request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download SDE: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Get content length for progress tracking
	contentLength := resp.ContentLength

	// Create progress writer
	pw := &progressWriter{
		writer:        out,
		total:         contentLength,
		verbose:       d.config.Verbose,
		lastPrintTime: time.Now(),
	}

	// Copy the response body to file with progress tracking
	bytesWritten, err := io.Copy(pw, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write SDE file: %w", err)
	}

	if d.config.Verbose {
		fmt.Printf("\nDownload complete: %d bytes\n", bytesWritten)
	}

	success = true
	return &DownloadResult{
		ZipPath:   zipPath,
		BytesRead: bytesWritten,
	}, nil
}

// Extract extracts a ZIP archive to the specified destination directory.
// Returns the path to the extracted SDE directory.
func (d *Downloader) Extract(zipPath, destDir string) (string, error) {
	if d.config.Verbose {
		fmt.Printf("Extracting SDE to: %s\n", destDir)
	}

	// Open the ZIP file
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer func() { _ = r.Close() }()

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	totalFiles := len(r.File)
	extractedCount := 0

	// Extract each file
	for _, f := range r.File {
		if err := d.extractFile(f, destDir); err != nil {
			return "", fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}
		extractedCount++

		if d.config.Verbose && extractedCount%1000 == 0 {
			fmt.Printf("Extracted %d/%d files...\n", extractedCount, totalFiles)
		}
	}

	if d.config.Verbose {
		fmt.Printf("Extraction complete: %d files\n", extractedCount)
	}

	// The new SDE format extracts directly to the destination directory
	return destDir, nil
}

// extractFile extracts a single file from the ZIP archive.
func (d *Downloader) extractFile(f *zip.File, destDir string) error {
	// Sanitize the file path to prevent zip slip attacks
	// Using filepath.IsLocal is the recommended approach per Go documentation
	if !filepath.IsLocal(f.Name) {
		return fmt.Errorf("illegal file path: %s", f.Name)
	}
	destPath := filepath.Join(destDir, f.Name)

	// Handle directories
	if f.FileInfo().IsDir() {
		return os.MkdirAll(destPath, f.Mode())
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Create the file
	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()

	// Open the file in the archive
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	// Copy contents
	_, err = io.Copy(outFile, rc)
	return err
}

// Validate checks that the SDE directory has the expected structure.
func (d *Downloader) Validate(sdePath string) error {
	if d.config.Verbose {
		fmt.Printf("Validating SDE structure at: %s\n", sdePath)
	}

	// Check for expected files (new flat SDE format)
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

// DownloadAndExtract is a convenience method that downloads and extracts the SDE.
func (d *Downloader) DownloadAndExtract(ctx context.Context) (string, error) {
	// Download
	result, err := d.Download(ctx)
	if err != nil {
		return "", err
	}

	// Determine extraction directory
	extractDir := filepath.Dir(result.ZipPath)

	// Extract
	sdePath, err := d.Extract(result.ZipPath, extractDir)
	if err != nil {
		// Clean up temp directory on extraction failure
		_ = os.RemoveAll(extractDir)
		return "", err
	}

	// Clean up ZIP file
	if err := os.Remove(result.ZipPath); err != nil && d.config.Verbose {
		fmt.Printf("Warning: failed to remove ZIP file: %v\n", err)
	}

	// Validate
	if err := d.Validate(sdePath); err != nil {
		// Clean up extracted directory on validation failure
		_ = os.RemoveAll(extractDir)
		return "", fmt.Errorf("SDE validation failed: %w", err)
	}

	return sdePath, nil
}

// progressWriter wraps an io.Writer to track and display progress.
type progressWriter struct {
	writer        io.Writer
	total         int64
	written       int64
	verbose       bool
	lastPrintTime time.Time
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.written += int64(n)

	// Print progress every second
	if pw.verbose && time.Since(pw.lastPrintTime) > time.Second {
		pw.printProgress()
		pw.lastPrintTime = time.Now()
	}

	return n, err
}

func (pw *progressWriter) printProgress() {
	if pw.total > 0 {
		percent := float64(pw.written) / float64(pw.total) * 100
		fmt.Printf("\rDownloading: %.1f%% (%s / %s)",
			percent,
			formatBytes(pw.written),
			formatBytes(pw.total))
	} else {
		fmt.Printf("\rDownloading: %s", formatBytes(pw.written))
	}
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
