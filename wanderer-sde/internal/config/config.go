// Package config provides configuration management for the SDE converter.
package config

// SDELatestURL is the download URL for the latest EVE SDE YAML archive.
// This is a shorthand URL that redirects to the latest build number.
const SDELatestURL = "https://developers.eveonline.com/static-data/eve-online-static-data-latest-yaml.zip"

// OutputFormat specifies the output file format.
type OutputFormat string

const (
	// FormatCSV outputs data in CSV format (Fuzzwork-compatible).
	FormatCSV OutputFormat = "csv"
	// FormatJSON outputs data in JSON format.
	FormatJSON OutputFormat = "json"
)

// Config holds all configuration options for the converter.
type Config struct {
	// SDEPath is the path to the SDE directory or ZIP file.
	SDEPath string

	// OutputDir is the directory where output files will be written.
	OutputDir string

	// SDEUrl is the URL to download the SDE from.
	SDEUrl string

	// DownloadSDE indicates whether to download the SDE.
	DownloadSDE bool

	// Verbose enables verbose logging.
	Verbose bool

	// PassthroughDir is the directory containing existing Wanderer JSON files to copy.
	PassthroughDir string

	// PrettyPrint enables indented JSON output (only applies to JSON format).
	PrettyPrint bool

	// OutputFormat specifies the output file format (csv or json).
	OutputFormat OutputFormat

	// Version is the application version for User-Agent headers.
	Version string
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		OutputDir:    "./output",
		SDEUrl:       SDELatestURL,
		PrettyPrint:  true,
		OutputFormat: FormatCSV, // Default to CSV for Fuzzwork compatibility
	}
}

// Validate checks that the configuration is valid.
func (c *Config) Validate() error {
	if !c.DownloadSDE && c.SDEPath == "" {
		return ErrNoSDESource
	}
	if c.OutputDir == "" {
		return ErrNoOutputDir
	}
	return nil
}
