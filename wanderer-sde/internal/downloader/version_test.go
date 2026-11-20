package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
)

func TestVersionChecker_GetStoredVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "version_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	vc := NewVersionChecker(cfg)

	// Test no stored version
	version, err := vc.GetStoredVersion(tmpDir)
	if err != nil {
		t.Errorf("GetStoredVersion failed: %v", err)
	}
	if version != "" {
		t.Errorf("Expected empty version, got %q", version)
	}

	// Store a version
	versionFile := filepath.Join(tmpDir, VersionFileName)
	if err := os.WriteFile(versionFile, []byte("123456\n"), 0644); err != nil {
		t.Fatalf("failed to write version file: %v", err)
	}

	// Test stored version
	version, err = vc.GetStoredVersion(tmpDir)
	if err != nil {
		t.Errorf("GetStoredVersion failed: %v", err)
	}
	if version != "123456" {
		t.Errorf("Expected version '123456', got %q", version)
	}
}

func TestVersionChecker_StoreVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "version_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	vc := NewVersionChecker(cfg)

	// Store a version
	subDir := filepath.Join(tmpDir, "subdir")
	if err := vc.StoreVersion(subDir, "789012"); err != nil {
		t.Fatalf("StoreVersion failed: %v", err)
	}

	// Verify the file was created
	versionFile := filepath.Join(subDir, VersionFileName)
	data, err := os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("failed to read version file: %v", err)
	}

	if string(data) != "789012\n" {
		t.Errorf("Version file content mismatch: got %q, want %q", string(data), "789012\n")
	}
}

func TestVersionChecker_GetLatestVersion(t *testing.T) {
	// A release index as served by sde.evetech.net/tranquility/latest/
	indexResponse := `{"release":"20250115T110000Z","buildNumber":2025001,"releaseDate":"2025-01-15T11:00:00Z",
"tables":{"types":{"version":"abc","jsonl":{"url":"https://example.invalid/types.jsonl"}}}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", "\"test-etag\"")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(indexResponse))
	}))
	defer server.Close()

	cfg := &config.Config{Verbose: false}
	vc := &VersionChecker{
		config:     cfg,
		httpClient: http.DefaultClient,
		versionURL: server.URL,
	}

	ctx := context.Background()
	info, err := vc.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion failed: %v", err)
	}

	if info.Release != "20250115T110000Z" {
		t.Errorf("Expected Release '20250115T110000Z', got %q", info.Release)
	}
	if info.BuildNumber != "2025001" {
		t.Errorf("Expected BuildNumber '2025001', got %q", info.BuildNumber)
	}
	if info.ReleaseDate != "2025-01-15T11:00:00Z" {
		t.Errorf("Expected ReleaseDate '2025-01-15T11:00:00Z', got %q", info.ReleaseDate)
	}
	if info.ETag != "\"test-etag\"" {
		t.Errorf("Expected ETag '\"test-etag\"', got %q", info.ETag)
	}
}

func TestVersionChecker_NeedsUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "version_update_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a test server that returns a specific SDE release
	indexResponse := `{"release":"20250116T110000Z","buildNumber":2025002,"releaseDate":"2025-01-16T11:00:00Z","tables":{}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(indexResponse))
	}))
	defer server.Close()

	cfg := &config.Config{Verbose: false}
	vc := &VersionChecker{
		config:     cfg,
		httpClient: http.DefaultClient,
		versionURL: server.URL,
	}

	ctx := context.Background()

	// Test 1: No stored version should need update
	needsUpdate, versionInfo, err := vc.NeedsUpdate(ctx, tmpDir)
	if err != nil {
		t.Fatalf("NeedsUpdate failed: %v", err)
	}
	if !needsUpdate {
		t.Error("Expected needsUpdate=true when no stored version exists")
	}
	if versionInfo == nil {
		t.Fatal("Expected versionInfo to be non-nil")
	}
	if versionInfo.Release != "20250116T110000Z" || versionInfo.BuildNumber != "2025002" {
		t.Errorf("Expected release 20250116T110000Z build 2025002, got %+v", versionInfo)
	}

	// Test 2: Store an old release, should need update
	if err := vc.StoreVersion(tmpDir, "20250115T110000Z"); err != nil {
		t.Fatalf("StoreVersion failed: %v", err)
	}

	needsUpdate, versionInfo, err = vc.NeedsUpdate(ctx, tmpDir)
	if err != nil {
		t.Fatalf("NeedsUpdate failed: %v", err)
	}
	if !needsUpdate {
		t.Error("Expected needsUpdate=true when stored version is older")
	}
	if versionInfo.BuildNumber != "2025002" {
		t.Errorf("Expected BuildNumber '2025002', got %q", versionInfo.BuildNumber)
	}

	// Test 3: A build number stored before the SDE API never matches a
	// release, so the first run after the switch downloads again
	if err := vc.StoreVersion(tmpDir, "2025002"); err != nil {
		t.Fatalf("StoreVersion failed: %v", err)
	}

	needsUpdate, _, err = vc.NeedsUpdate(ctx, tmpDir)
	if err != nil {
		t.Fatalf("NeedsUpdate failed: %v", err)
	}
	if !needsUpdate {
		t.Error("Expected needsUpdate=true when the stored version is a build number")
	}

	// Test 4: Store the same release, should not need update
	if err := vc.StoreVersion(tmpDir, "20250116T110000Z"); err != nil {
		t.Fatalf("StoreVersion failed: %v", err)
	}

	needsUpdate, versionInfo, err = vc.NeedsUpdate(ctx, tmpDir)
	if err != nil {
		t.Fatalf("NeedsUpdate failed: %v", err)
	}
	if needsUpdate {
		t.Error("Expected needsUpdate=false when stored version matches latest")
	}
	if versionInfo.BuildNumber != "2025002" {
		t.Errorf("Expected BuildNumber '2025002', got %q", versionInfo.BuildNumber)
	}
}

func TestVersionChecker_CheckETag(t *testing.T) {
	tests := []struct {
		name           string
		storedETag     string
		serverETag     string
		serverStatus   int
		expectedUpdate bool
		expectError    bool
	}{
		{
			name:           "no previous etag",
			storedETag:     "",
			serverETag:     "\"new-etag\"",
			serverStatus:   http.StatusOK,
			expectedUpdate: true,
			expectError:    false,
		},
		{
			name:           "etag unchanged (304)",
			storedETag:     "\"same-etag\"",
			serverETag:     "\"same-etag\"",
			serverStatus:   http.StatusNotModified,
			expectedUpdate: false,
			expectError:    false,
		},
		{
			name:           "etag changed",
			storedETag:     "\"old-etag\"",
			serverETag:     "\"new-etag\"",
			serverStatus:   http.StatusOK,
			expectedUpdate: true,
			expectError:    false,
		},
		{
			name:           "server error",
			storedETag:     "\"etag\"",
			serverETag:     "",
			serverStatus:   http.StatusInternalServerError,
			expectedUpdate: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.serverETag != "" {
					w.Header().Set("ETag", tt.serverETag)
				}
				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			cfg := &config.Config{Verbose: false}
			vc := &VersionChecker{
				config:     cfg,
				httpClient: server.Client(),
			}

			ctx := context.Background()
			needsUpdate, _, err := vc.CheckETag(ctx, server.URL, tt.storedETag)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if needsUpdate != tt.expectedUpdate {
				t.Errorf("needsUpdate = %v, want %v", needsUpdate, tt.expectedUpdate)
			}
		})
	}
}

func TestNewVersionChecker(t *testing.T) {
	cfg := &config.Config{Verbose: true}
	vc := NewVersionChecker(cfg)

	if vc == nil {
		t.Fatal("NewVersionChecker returned nil")
	}
	if vc.config != cfg {
		t.Error("Config not set correctly")
	}
	if vc.httpClient == nil {
		t.Error("HTTP client not initialized")
	}
}

func TestVersionInfo_Fields(t *testing.T) {
	vi := VersionInfo{
		Release:     "20250115T110000Z",
		BuildNumber: "2025001",
		ETag:        "\"test-etag\"",
	}

	if vi.Release != "20250115T110000Z" {
		t.Errorf("Release mismatch: got %q, want %q", vi.Release, "20250115T110000Z")
	}

	if vi.BuildNumber != "2025001" {
		t.Errorf("BuildNumber mismatch: got %q, want %q", vi.BuildNumber, "2025001")
	}
	if vi.ETag != "\"test-etag\"" {
		t.Errorf("ETag mismatch: got %q, want %q", vi.ETag, "\"test-etag\"")
	}
}
