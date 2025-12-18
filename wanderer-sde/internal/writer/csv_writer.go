// Package writer provides output generation for the SDE converter.
package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

// CSV file names matching Fuzzwork format.
const (
	CSVFileSolarSystems    = "mapSolarSystems.csv"
	CSVFileRegions         = "mapRegions.csv"
	CSVFileConstellations  = "mapConstellations.csv"
	CSVFileWormholeClasses = "mapLocationWormholeClasses.csv"
	CSVFileTypes           = "invTypes.csv"
	CSVFileGroups          = "invGroups.csv"
	CSVFileSystemJumps     = "mapSolarSystemJumps.csv"
	CSVFileNPCStations     = "npcStations.csv"
)

// CSVWriter handles writing converted data to CSV files.
type CSVWriter struct {
	config    *config.Config
	outputDir string
}

// NewCSVWriter creates a new CSVWriter with the given configuration.
func NewCSVWriter(cfg *config.Config) *CSVWriter {
	return &CSVWriter{
		config:    cfg,
		outputDir: cfg.OutputDir,
	}
}

// WriteAll writes all converted data to CSV files.
func (w *CSVWriter) WriteAll(data *models.ConvertedData) error {
	// Validate input data
	if data == nil {
		return fmt.Errorf("converted data is nil")
	}
	if data.Universe == nil {
		return fmt.Errorf("universe data is nil")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if w.config.Verbose {
		fmt.Printf("Writing CSV files to: %s\n", w.outputDir)
	}

	// Write all data files
	if err := w.WriteSolarSystems(data.Universe.SolarSystems); err != nil {
		return fmt.Errorf("failed to write solar systems: %w", err)
	}

	if err := w.WriteRegions(data.Universe.Regions); err != nil {
		return fmt.Errorf("failed to write regions: %w", err)
	}

	if err := w.WriteConstellations(data.Universe.Constellations); err != nil {
		return fmt.Errorf("failed to write constellations: %w", err)
	}

	if err := w.WriteWormholeClasses(data.WormholeClasses); err != nil {
		return fmt.Errorf("failed to write wormhole classes: %w", err)
	}

	if err := w.WriteTypes(data.InvTypes); err != nil {
		return fmt.Errorf("failed to write types: %w", err)
	}

	if err := w.WriteGroups(data.InvGroups); err != nil {
		return fmt.Errorf("failed to write groups: %w", err)
	}

	if err := w.WriteSystemJumps(data.SystemJumps); err != nil {
		return fmt.Errorf("failed to write system jumps: %w", err)
	}

	if err := w.WriteNPCStations(data.NPCStations); err != nil {
		return fmt.Errorf("failed to write NPC stations: %w", err)
	}

	return nil
}

// WriteSolarSystems writes solar system data to CSV.
func (w *CSVWriter) WriteSolarSystems(systems []models.SolarSystem) error {
	rows := make([][]string, len(systems))
	for i, s := range systems {
		rows[i] = s.ToCSVRow()
	}
	return w.writeCSV(CSVFileSolarSystems, "mapSolarSystems", rows)
}

// WriteRegions writes region data to CSV.
func (w *CSVWriter) WriteRegions(regions []models.Region) error {
	rows := make([][]string, len(regions))
	for i, r := range regions {
		rows[i] = r.ToCSVRow()
	}
	return w.writeCSV(CSVFileRegions, "mapRegions", rows)
}

// WriteConstellations writes constellation data to CSV.
func (w *CSVWriter) WriteConstellations(constellations []models.Constellation) error {
	rows := make([][]string, len(constellations))
	for i, c := range constellations {
		rows[i] = c.ToCSVRow()
	}
	return w.writeCSV(CSVFileConstellations, "mapConstellations", rows)
}

// WriteWormholeClasses writes wormhole class data to CSV.
func (w *CSVWriter) WriteWormholeClasses(classes []models.WormholeClassLocation) error {
	rows := make([][]string, len(classes))
	for i, c := range classes {
		rows[i] = c.ToCSVRow()
	}
	return w.writeCSV(CSVFileWormholeClasses, "mapLocationWormholeClasses", rows)
}

// WriteTypes writes type data to CSV.
func (w *CSVWriter) WriteTypes(types []models.InvType) error {
	rows := make([][]string, len(types))
	for i, t := range types {
		rows[i] = t.ToCSVRow()
	}
	return w.writeCSV(CSVFileTypes, "invTypes", rows)
}

// WriteGroups writes group data to CSV.
func (w *CSVWriter) WriteGroups(groups []models.InvGroup) error {
	rows := make([][]string, len(groups))
	for i, g := range groups {
		rows[i] = g.ToCSVRow()
	}
	return w.writeCSV(CSVFileGroups, "invGroups", rows)
}

// WriteSystemJumps writes system jump data to CSV.
func (w *CSVWriter) WriteSystemJumps(jumps []models.SystemJump) error {
	rows := make([][]string, len(jumps))
	for i, j := range jumps {
		rows[i] = j.ToCSVRow()
	}
	return w.writeCSV(CSVFileSystemJumps, "mapSolarSystemJumps", rows)
}

// WriteNPCStations writes NPC station data to CSV.
func (w *CSVWriter) WriteNPCStations(stations []models.NPCStation) error {
	rows := make([][]string, len(stations))
	for i, s := range stations {
		rows[i] = s.ToCSVRow()
	}
	return w.writeCSV(CSVFileNPCStations, "npcStations", rows)
}

// CopyPassthroughFiles copies community-maintained JSON files from the source directory.
// For CSV output, we still want to copy these JSON files as they're used by Wanderer.
func (w *CSVWriter) CopyPassthroughFiles(sourceDir string) error {
	if sourceDir == "" {
		return nil
	}

	if w.config.Verbose {
		fmt.Printf("Copying passthrough files from: %s\n", sourceDir)
	}

	var copied, skipped int
	for _, filename := range PassthroughFiles {
		srcPath := filepath.Join(sourceDir, filename)
		dstPath := filepath.Join(w.outputDir, filename)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			if w.config.Verbose {
				fmt.Printf("  Skipping %s (not found)\n", filename)
			}
			skipped++
			continue
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", filename, err)
		}

		if w.config.Verbose {
			fmt.Printf("  Copied %s\n", filename)
		}
		copied++
	}

	if w.config.Verbose {
		fmt.Printf("Passthrough complete: %d copied, %d skipped\n", copied, skipped)
	}

	return nil
}

// writeCSV writes data rows to a CSV file with the appropriate headers.
func (w *CSVWriter) writeCSV(filename, headerKey string, rows [][]string) (err error) {
	path := filepath.Join(w.outputDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", path, closeErr)
		}
	}()

	csvWriter := csv.NewWriter(file)
	defer func() {
		csvWriter.Flush()
		if flushErr := csvWriter.Error(); flushErr != nil && err == nil {
			err = fmt.Errorf("failed to flush CSV writer for %s: %w", path, flushErr)
		}
	}()

	// Write header row
	headers, ok := models.CSVHeaders[headerKey]
	if !ok {
		return fmt.Errorf("no headers defined for %s", headerKey)
	}
	if err := csvWriter.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers to %s: %w", path, err)
	}

	// Write data rows
	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write row to %s: %w", path, err)
		}
	}

	if w.config.Verbose {
		fmt.Printf("  Wrote %s (%d rows)\n", filename, len(rows))
	}

	return nil
}
