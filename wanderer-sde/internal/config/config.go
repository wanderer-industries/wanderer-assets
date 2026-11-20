// Package config provides configuration management for the SDE converter.
package config

// SDEAPIURL lists the Tranquility releases of CCP's SDE API. latest/ below
// it redirects to the current release.
const SDEAPIURL = "https://sde.evetech.net/tranquility/"

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
	// SDEPath is the path to a directory of SDE JSONL tables (<table>.jsonl).
	SDEPath string

	// OutputDir is the directory where output files will be written.
	OutputDir string

	// SDEUrl is the SDE API URL that lists the releases to download from.
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
		SDEUrl:       SDEAPIURL,
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
