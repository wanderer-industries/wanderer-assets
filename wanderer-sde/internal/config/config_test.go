package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}

	// Check default values
	if cfg.OutputDir != "./output" {
		t.Errorf("Expected default OutputDir './output', got %q", cfg.OutputDir)
	}

	if cfg.SDEUrl != SDELatestURL {
		t.Errorf("Expected default SDEUrl %q, got %q", SDELatestURL, cfg.SDEUrl)
	}

	if !cfg.PrettyPrint {
		t.Error("Expected PrettyPrint to default to true")
	}

	// Check unset defaults
	if cfg.SDEPath != "" {
		t.Errorf("Expected empty SDEPath, got %q", cfg.SDEPath)
	}

	if cfg.DownloadSDE {
		t.Error("Expected DownloadSDE to default to false")
	}

	if cfg.Verbose {
		t.Error("Expected Verbose to default to false")
	}

	if cfg.PassthroughDir != "" {
		t.Errorf("Expected empty PassthroughDir, got %q", cfg.PassthroughDir)
	}

	if cfg.OutputFormat != FormatCSV {
		t.Errorf("Expected default OutputFormat FormatCSV, got %v", cfg.OutputFormat)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError error
	}{
		{
			name: "valid with SDE path",
			config: &Config{
				SDEPath:   "/path/to/sde",
				OutputDir: "./output",
			},
			expectError: nil,
		},
		{
			name: "valid with download flag",
			config: &Config{
				DownloadSDE: true,
				OutputDir:   "./output",
			},
			expectError: nil,
		},
		{
			name: "valid with both SDE path and download",
			config: &Config{
				SDEPath:     "/path/to/sde",
				DownloadSDE: true,
				OutputDir:   "./output",
			},
			expectError: nil,
		},
		{
			name: "missing SDE source",
			config: &Config{
				OutputDir: "./output",
			},
			expectError: ErrNoSDESource,
		},
		{
			name: "missing output directory",
			config: &Config{
				SDEPath:   "/path/to/sde",
				OutputDir: "",
			},
			expectError: ErrNoOutputDir,
		},
		{
			name: "missing both SDE source and output",
			config: &Config{
				OutputDir: "",
			},
			expectError: ErrNoSDESource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError == nil {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectError)
				} else if err != tt.expectError {
					t.Errorf("Expected error %v, got %v", tt.expectError, err)
				}
			}
		})
	}
}

func TestSDELatestURL(t *testing.T) {
	// Verify the URL is properly set
	expectedURL := "https://developers.eveonline.com/static-data/eve-online-static-data-latest-yaml.zip"
	if SDELatestURL != expectedURL {
		t.Errorf("SDELatestURL mismatch: got %q, want %q", SDELatestURL, expectedURL)
	}
}

func TestErrors(t *testing.T) {
	// Verify error messages are meaningful
	if ErrNoSDESource.Error() == "" {
		t.Error("ErrNoSDESource has empty message")
	}
	if ErrNoOutputDir.Error() == "" {
		t.Error("ErrNoOutputDir has empty message")
	}
}
