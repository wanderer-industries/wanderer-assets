package downloader

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/guarzo/wanderer-sde/internal/config"
)

// VersionFileName is the name of the file that stores the last processed SDE release.
const VersionFileName = ".last-sde-version"

// VersionChecker handles checking and tracking SDE versions.
type VersionChecker struct {
	config     *config.Config
	httpClient *http.Client
	versionURL string // index of the latest release (defaults to <SDEUrl>/latest/)
}

// NewVersionChecker creates a new VersionChecker.
func NewVersionChecker(cfg *config.Config) *VersionChecker {
	return &VersionChecker{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		versionURL: releaseURL(cfg.SDEUrl, ""),
	}
}

// VersionInfo identifies an SDE release. CCP can publish several releases
// with the same build number, so Release is what tells them apart.
type VersionInfo struct {
	Release     string `json:"release"`
	BuildNumber string `json:"sde"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	ETag        string `json:"-"`
}

// GetLatestVersion reads the index of the latest SDE release.
func (vc *VersionChecker) GetLatestVersion(ctx context.Context) (*VersionInfo, error) {
	if vc.config.Verbose {
		fmt.Println("Checking latest SDE version...")
	}

	url := vc.versionURL
	if url == "" {
		url = releaseURL(vc.config.SDEUrl, "")
	}

	var rel Release
	header, err := getJSON(ctx, vc.httpClient, vc.config.Version, url, &rel)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}
	if rel.Release == "" || rel.BuildNumber <= 0 {
		return nil, fmt.Errorf("SDE release not found in %s", url)
	}

	info := rel.VersionInfo()
	info.ETag = header.Get("ETag")
	return info, nil
}

// GetStoredVersion retrieves the previously stored SDE release.
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

// StoreVersion saves the current SDE release.
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

// NeedsUpdate checks if the SDE needs to be updated. It compares releases:
// a version file written before the SDE API holds a build number, which
// never matches, so the first run after the switch downloads again.
func (vc *VersionChecker) NeedsUpdate(ctx context.Context, storageDir string) (bool, *VersionInfo, error) {
	latest, err := vc.GetLatestVersion(ctx)
	if err != nil {
		return false, nil, err
	}

	if vc.config.Verbose {
		fmt.Printf("Latest SDE release: %s (build %s)\n", latest.Release, latest.BuildNumber)
	}

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
		fmt.Printf("Stored SDE release: %s\n", stored)
	}

	needsUpdate := stored != latest.Release
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
// If url is empty, the index of the latest release is used.
func (vc *VersionChecker) CheckETag(ctx context.Context, url string, storedETag string) (bool, string, error) {
	if url == "" {
		url = vc.versionURL
		if url == "" {
			url = releaseURL(vc.config.SDEUrl, "")
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
