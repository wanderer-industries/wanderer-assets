// Package writer provides JSON output generation for the SDE converter.
package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

// OutputFiles defines the names of generated JSON files.
const (
	FileSolarSystems    = "mapSolarSystems.json"
	FileRegions         = "mapRegions.json"
	FileConstellations  = "mapConstellations.json"
	FileWormholeClasses = "mapLocationWormholeClasses.json"
	FileShipTypes       = "invTypes.json"
	FileItemGroups      = "invGroups.json"
	FileSystemJumps     = "mapSolarSystemJumps.json"
	FileNPCStations     = "npcStations.json"
)

// PassthroughFiles lists the community-maintained JSON files to copy.
var PassthroughFiles = []string{
	"wormholes.json",
	"wormholeClasses.json",
	"wormholeClassesInfo.json",
	"wormholeSystems.json",
	"triglavianSystems.json",
	"effects.json",
	"shatteredConstellations.json",
	"sunTypes.json",
	"triglavianEffectsByFaction.json",
}

// JSONWriter handles writing converted data to JSON files.
type JSONWriter struct {
	config    *config.Config
	outputDir string
	pretty    bool
}

// NewJSONWriter creates a new JSONWriter with the given configuration.
func NewJSONWriter(cfg *config.Config) *JSONWriter {
	return &JSONWriter{
		config:    cfg,
		outputDir: cfg.OutputDir,
		pretty:    cfg.PrettyPrint,
	}
}

// WriteAll writes all converted data to JSON files.
func (w *JSONWriter) WriteAll(data *models.ConvertedData) error {
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
		fmt.Printf("Writing JSON files to: %s\n", w.outputDir)
	}

	// Write universe data files
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

// WriteSolarSystems writes solar system data to JSON.
func (w *JSONWriter) WriteSolarSystems(systems []models.SolarSystem) error {
	return w.writeJSON(FileSolarSystems, systems)
}

// WriteRegions writes region data to JSON.
func (w *JSONWriter) WriteRegions(regions []models.Region) error {
	return w.writeJSON(FileRegions, regions)
}

// WriteConstellations writes constellation data to JSON.
func (w *JSONWriter) WriteConstellations(constellations []models.Constellation) error {
	return w.writeJSON(FileConstellations, constellations)
}

// WriteWormholeClasses writes wormhole class data to JSON.
func (w *JSONWriter) WriteWormholeClasses(classes []models.WormholeClassLocation) error {
	return w.writeJSON(FileWormholeClasses, classes)
}

// WriteTypes writes type data to JSON.
func (w *JSONWriter) WriteTypes(types []models.InvType) error {
	return w.writeJSON(FileShipTypes, types)
}

// WriteGroups writes group data to JSON.
func (w *JSONWriter) WriteGroups(groups []models.InvGroup) error {
	return w.writeJSON(FileItemGroups, groups)
}

// WriteSystemJumps writes system jump data to JSON.
func (w *JSONWriter) WriteSystemJumps(jumps []models.SystemJump) error {
	return w.writeJSON(FileSystemJumps, jumps)
}

// WriteNPCStations writes NPC station data to JSON.
func (w *JSONWriter) WriteNPCStations(stations []models.NPCStation) error {
	return w.writeJSON(FileNPCStations, stations)
}

// CopyPassthroughFiles copies community-maintained JSON files from the source directory.
func (w *JSONWriter) CopyPassthroughFiles(sourceDir string) error {
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

// writeJSON marshals data to JSON and writes it to a file.
func (w *JSONWriter) writeJSON(filename string, data interface{}) (err error) {
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

	encoder := json.NewEncoder(file)
	if w.pretty {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to %s: %w", path, err)
	}

	if w.config.Verbose {
		fmt.Printf("  Wrote %s\n", filename)
	}

	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return dstFile.Sync()
}
