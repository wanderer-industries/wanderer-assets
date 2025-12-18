package downloader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/guarzo/wanderer-sde/internal/config"
)

// LatestJSONLURL is the URL to check for the latest SDE build number.
const LatestJSONLURL = "https://developers.eveonline.com/static-data/tranquility/latest.jsonl"

// VersionFileName is the name of the file that stores the last processed SDE version.
const VersionFileName = ".last-sde-version"

// VersionChecker handles checking and tracking SDE versions.
type VersionChecker struct {
	config     *config.Config
	httpClient *http.Client
	versionURL string // URL to check for latest version (defaults to LatestJSONLURL)
}

// NewVersionChecker creates a new VersionChecker.
func NewVersionChecker(cfg *config.Config) *VersionChecker {
	return &VersionChecker{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		versionURL: LatestJSONLURL,
	}
}

// VersionInfo contains information about the SDE version.
type VersionInfo struct {
	BuildNumber string `json:"sde"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	ETag        string `json:"-"`
}

// latestRecord represents a record from the latest.jsonl file.
type latestRecord struct {
	Key         string `json:"_key"`
	BuildNumber int64  `json:"buildNumber,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}

// GetLatestVersion fetches the latest SDE version from CCP.
func (vc *VersionChecker) GetLatestVersion(ctx context.Context) (*VersionInfo, error) {
	if vc.config.Verbose {
		fmt.Println("Checking latest SDE version...")
	}

	url := vc.versionURL
	if url == "" {
		url = LatestJSONLURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent header
	version := vc.config.Version
	if version == "" {
		version = "dev"
	}
	req.Header.Set("User-Agent", fmt.Sprintf("wanderer-sde/%s (https://github.com/guarzo/wanderer-sde)", version))

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse JSON Lines format - each line is a separate JSON object
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record latestRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			// Skip malformed lines
			if vc.config.Verbose {
				fmt.Printf("Warning: skipping malformed JSONL line: %v\n", err)
			}
			continue
		}

		// Look for the "sde" key with buildNumber
		if record.Key == "sde" && record.BuildNumber > 0 {
			return &VersionInfo{
				BuildNumber: fmt.Sprintf("%d", record.BuildNumber),
				ReleaseDate: record.ReleaseDate,
				ETag:        resp.Header.Get("ETag"),
			}, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return nil, fmt.Errorf("SDE version not found in latest.jsonl")
}

// GetStoredVersion retrieves the previously stored SDE version.
func (vc *VersionChecker) GetStoredVersion(dir string) (string, error) {
	versionFile := filepath.Join(dir, VersionFileName)

	data, err := os.ReadFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No stored version
		}
		return "", fmt.Errorf("failed to read version file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

// StoreVersion saves the current SDE version.
func (vc *VersionChecker) StoreVersion(dir, version string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	versionFile := filepath.Join(dir, VersionFileName)
	if err := os.WriteFile(versionFile, []byte(version+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}

	return nil
}

// NeedsUpdate checks if the SDE needs to be updated.
func (vc *VersionChecker) NeedsUpdate(ctx context.Context, storageDir string) (bool, *VersionInfo, error) {
	// Get the latest version
	latest, err := vc.GetLatestVersion(ctx)
	if err != nil {
		return false, nil, err
	}

	if vc.config.Verbose {
		fmt.Printf("Latest SDE version: %s\n", latest.BuildNumber)
	}

	// Get the stored version
	stored, err := vc.GetStoredVersion(storageDir)
	if err != nil {
		return false, nil, err
	}

	if stored == "" {
		if vc.config.Verbose {
			fmt.Println("No stored version found, update needed")
		}
		return true, latest, nil
	}

	if vc.config.Verbose {
		fmt.Printf("Stored SDE version: %s\n", stored)
	}

	needsUpdate := stored != latest.BuildNumber
	if vc.config.Verbose {
		if needsUpdate {
			fmt.Println("Update available")
		} else {
			fmt.Println("SDE is up to date")
		}
	}

	return needsUpdate, latest, nil
}

// CheckETag performs an HTTP HEAD request to check if the SDE has been updated
// using ETag headers, which is more efficient than downloading the full version info.
// If url is empty, the default versionURL (or LatestJSONLURL) is used.
func (vc *VersionChecker) CheckETag(ctx context.Context, url string, storedETag string) (bool, string, error) {
	if url == "" {
		url = vc.versionURL
		if url == "" {
			url = LatestJSONLURL
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, "", fmt.Errorf("failed to create request: %w", err)
	}

	if storedETag != "" {
		req.Header.Set("If-None-Match", storedETag)
	}

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("failed to check ETag: %w", err)
	}
	_ = resp.Body.Close()

	currentETag := resp.Header.Get("ETag")

	// 304 Not Modified means no update needed
	if resp.StatusCode == http.StatusNotModified {
		return false, currentETag, nil
	}

	if resp.StatusCode == http.StatusOK {
		// ETag changed or no previous ETag
		return storedETag != currentETag, currentETag, nil
	}

	return false, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
