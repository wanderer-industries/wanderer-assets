// Package writer provides output generation for the SDE converter.
package writer

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

// Writer defines the interface for writing converted data to files.
type Writer interface {
	// WriteAll writes all converted data to output files.
	WriteAll(data *models.ConvertedData) error

	// CopyPassthroughFiles copies community-maintained files from the source directory.
	CopyPassthroughFiles(sourceDir string) error
}

// NewWriter creates a new Writer based on the configured output format.
func NewWriter(cfg *config.Config) (Writer, error) {
	switch cfg.OutputFormat {
	case config.FormatCSV:
		return NewCSVWriter(cfg), nil
	case config.FormatJSON:
		return NewJSONWriter(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", cfg.OutputFormat)
	}
}

// GetOutputFiles returns the list of output file names based on format.
func GetOutputFiles(format config.OutputFormat) []string {
	switch format {
	case config.FormatCSV:
		return []string{
			CSVFileSolarSystems,
			CSVFileRegions,
			CSVFileConstellations,
			CSVFileWormholeClasses,
			CSVFileTypes,
			CSVFileGroups,
			CSVFileSystemJumps,
			CSVFileNPCStations,
		}
	case config.FormatJSON:
		return []string{
			FileSolarSystems,
			FileRegions,
			FileConstellations,
			FileWormholeClasses,
			FileShipTypes,
			FileItemGroups,
			FileSystemJumps,
			FileNPCStations,
		}
	default:
		return nil
	}
}
