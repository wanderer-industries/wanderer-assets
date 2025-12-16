package internal

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/internal/parser"
	"github.com/guarzo/wanderer-sde/internal/transformer"
	"github.com/guarzo/wanderer-sde/internal/writer"
)

// createTestSDE creates a complete minimal SDE structure for integration testing
func createTestSDE(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create mapRegions.yaml with multiple regions
	regionsYAML := `10000002:
  regionID: 10000002
  name:
    en: "The Forge"
10000001:
  regionID: 10000001
  name:
    en: "Derelik"
11000001:
  regionID: 11000001
  name:
    en: "J-Space Region"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(regionsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapRegions.yaml: %v", err)
	}

	// Create mapConstellations.yaml
	constellationsYAML := `20000020:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Kimotoro"
20000001:
  constellationID: 20000001
  regionID: 10000001
  name:
    en: "Joas"
21000001:
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J-Constellation"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapConstellations.yaml"), []byte(constellationsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapConstellations.yaml: %v", err)
	}

	// Create mapSolarSystems.yaml with new SDE format (2025+)
	systemsYAML := `30000142:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Jita"
  securityStatus: 0.9459
  starID: 40000006
30000001:
  constellationID: 20000001
  regionID: 10000001
  name:
    en: "Tanoo"
  securityStatus: 0.8576
  starID: 40000007
30000144:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Perimeter"
  securityStatus: 0.94
  starID: 40000008
30000145:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "LowSec System"
  securityStatus: 0.4
  starID: 40000009
30000146:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "NullSec System"
  securityStatus: -0.7
  starID: 40000010
31000001:
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J123456"
  securityStatus: -1.0
  starID: 40045041
  wormholeClassID: 3
31000002:
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J654321"
  securityStatus: -1.0
  starID: 40045041
  wormholeClassID: 5
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapSolarSystems.yaml"), []byte(systemsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapSolarSystems.yaml: %v", err)
	}

	// Create mapStargates.yaml
	stargatesYAML := `50000001:
  solarSystemID: 30000142
  destination:
    solarSystemID: 30000144
    stargateID: 50000002
  typeID: 16
50000002:
  solarSystemID: 30000144
  destination:
    solarSystemID: 30000142
    stargateID: 50000001
  typeID: 16
50000003:
  solarSystemID: 30000144
  destination:
    solarSystemID: 30000145
    stargateID: 50000004
  typeID: 16
50000004:
  solarSystemID: 30000145
  destination:
    solarSystemID: 30000144
    stargateID: 50000003
  typeID: 16
50000005:
  solarSystemID: 30000145
  destination:
    solarSystemID: 30000146
    stargateID: 50000006
  typeID: 16
50000006:
  solarSystemID: 30000146
  destination:
    solarSystemID: 30000145
    stargateID: 50000005
  typeID: 16
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapStargates.yaml"), []byte(stargatesYAML), 0644); err != nil {
		t.Fatalf("failed to create mapStargates.yaml: %v", err)
	}

	// Create types.yaml with ships and non-ships
	typesYAML := `587:
  groupID: 25
  name:
    en: "Rifter"
  mass: 1350000.0
  volume: 27500.0
  capacity: 125.0
  published: true
588:
  groupID: 25
  name:
    en: "Slasher"
  mass: 1200000.0
  volume: 26000.0
  capacity: 115.0
  published: true
589:
  groupID: 25
  name:
    en: "Breacher"
  mass: 1100000.0
  volume: 25000.0
  capacity: 130.0
  published: true
625:
  groupID: 26
  name:
    en: "Caracal"
  mass: 11000000.0
  volume: 92000.0
  capacity: 350.0
  published: true
2456:
  groupID: 18
  name:
    en: "Hobgoblin I"
  mass: 2500.0
  volume: 5.0
  published: true
17738:
  groupID: 419
  name:
    en: "Prophecy"
  mass: 15500000.0
  volume: 270000.0
  capacity: 500.0
  published: true
99999:
  groupID: 25
  name:
    en: "Unpublished Ship"
  mass: 1000000.0
  volume: 20000.0
  published: false
`
	if err := os.WriteFile(filepath.Join(tmpDir, "types.yaml"), []byte(typesYAML), 0644); err != nil {
		t.Fatalf("failed to create types.yaml: %v", err)
	}

	// Create groups.yaml
	groupsYAML := `25:
  categoryID: 6
  name:
    en: "Frigate"
  published: true
26:
  categoryID: 6
  name:
    en: "Cruiser"
  published: true
18:
  categoryID: 7
  name:
    en: "Drone"
  published: true
419:
  categoryID: 6
  name:
    en: "Battlecruiser"
  published: true
999:
  categoryID: 6
  name:
    en: "Unpublished Group"
  published: false
`
	if err := os.WriteFile(filepath.Join(tmpDir, "groups.yaml"), []byte(groupsYAML), 0644); err != nil {
		t.Fatalf("failed to create groups.yaml: %v", err)
	}

	// Create categories.yaml
	categoriesYAML := `6:
  name:
    en: "Ship"
  published: true
7:
  name:
    en: "Drone"
  published: true
`
	if err := os.WriteFile(filepath.Join(tmpDir, "categories.yaml"), []byte(categoriesYAML), 0644); err != nil {
		t.Fatalf("failed to create categories.yaml: %v", err)
	}

	// Create mapStars.yaml with star ID -> type ID mapping
	starsYAML := `40000006:
  solarSystemID: 30000142
  typeID: 3796
  radius: 123456789.0
40000007:
  solarSystemID: 30000001
  typeID: 3797
  radius: 987654321.0
40000008:
  solarSystemID: 30000144
  typeID: 3798
  radius: 111111111.0
40000009:
  solarSystemID: 30000145
  typeID: 3799
  radius: 222222222.0
40000010:
  solarSystemID: 30000146
  typeID: 3800
  radius: 333333333.0
40045041:
  solarSystemID: 31000001
  typeID: 45041
  radius: 444444444.0
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapStars.yaml"), []byte(starsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapStars.yaml: %v", err)
	}

	// Create npcStations.yaml with test stations
	npcStationsYAML := `60012736:
  ownerID: 1000137
  solarSystemID: 30000142
  typeID: 2501
  operationID: 14
60012919:
  ownerID: 1000137
  solarSystemID: 30000144
  typeID: 2498
  operationID: 14
60000001:
  ownerID: 1000002
  solarSystemID: 30000001
  typeID: 1531
  operationID: 26
`
	if err := os.WriteFile(filepath.Join(tmpDir, "npcStations.yaml"), []byte(npcStationsYAML), 0644); err != nil {
		t.Fatalf("failed to create npcStations.yaml: %v", err)
	}

	// Create npcCorporations.yaml with test corporations
	npcCorporationsYAML := `1000137:
  name:
    en: "DED"
  stationID: 60012736
  deleted: false
1000002:
  name:
    en: "CBD Corporation"
  stationID: 60000001
  deleted: false
`
	if err := os.WriteFile(filepath.Join(tmpDir, "npcCorporations.yaml"), []byte(npcCorporationsYAML), 0644); err != nil {
		t.Fatalf("failed to create npcCorporations.yaml: %v", err)
	}

	return tmpDir
}

func TestIntegration_FullPipeline(t *testing.T) {
	// Create test SDE
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	// Create output directory
	outputDir, err := os.MkdirTemp("", "integration_output")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	cfg := &config.Config{
		SDEPath:      sdeDir,
		OutputDir:    outputDir,
		OutputFormat: config.FormatJSON,
		Verbose:      false,
		PrettyPrint:  true,
	}

	// Step 1: Parse
	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	// Verify parsing results
	t.Run("parsing", func(t *testing.T) {
		if len(parseResult.Regions) != 3 {
			t.Errorf("Expected 3 regions, got %d", len(parseResult.Regions))
		}
		if len(parseResult.Constellations) != 3 {
			t.Errorf("Expected 3 constellations, got %d", len(parseResult.Constellations))
		}
		if len(parseResult.SolarSystems) != 7 {
			t.Errorf("Expected 7 solar systems, got %d", len(parseResult.SolarSystems))
		}
		if len(parseResult.Types) != 7 {
			t.Errorf("Expected 7 types, got %d", len(parseResult.Types))
		}
		if len(parseResult.Groups) != 5 {
			t.Errorf("Expected 5 groups, got %d", len(parseResult.Groups))
		}
		if len(parseResult.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(parseResult.Categories))
		}
		// 6 system jumps (3 stargate connections × 2 directions each)
		if len(parseResult.SystemJumps) != 6 {
			t.Errorf("Expected 6 system jumps (bidirectional), got %d", len(parseResult.SystemJumps))
		}
	})

	// Step 2: Transform
	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify transformation results
	t.Run("transformation", func(t *testing.T) {
		// Should include only ship types (groupID in ship groups with categoryID 6)
		// Test data has 6 ship types (types with groupID 25, 26, or 419)
		if len(convertedData.InvTypes) != 6 {
			t.Errorf("Expected 6 ship types, got %d", len(convertedData.InvTypes))
		}

		// Should include only ship groups (categoryID 6)
		// Test data has 4 ship groups (25, 26, 419, 999)
		if len(convertedData.InvGroups) != 4 {
			t.Errorf("Expected 4 ship groups, got %d", len(convertedData.InvGroups))
		}

		// Should extract wormhole classes (2 wormhole systems)
		if len(convertedData.WormholeClasses) != 2 {
			t.Errorf("Expected 2 wormhole classes, got %d", len(convertedData.WormholeClasses))
		}

		// Verify raw security values are preserved (Wanderer calculates true security itself)
		for _, sys := range convertedData.Universe.SolarSystems {
			switch sys.SolarSystemName {
			case "Jita":
				// Raw security value should be preserved
				if sys.Security != 0.9459 {
					t.Errorf("Jita security should be 0.9459 (raw), got %f", sys.Security)
				}
			case "J123456":
				// Wormhole security stays at -1.0
				if sys.Security != -1.0 {
					t.Errorf("J123456 security should be -1.0, got %f", sys.Security)
				}
			}
		}
	})

	// Step 3: Validate
	t.Run("validation", func(t *testing.T) {
		validationResult := tr.Validate(convertedData)

		// Should have warnings about low counts (test data is minimal)
		if len(validationResult.Warnings) == 0 {
			t.Log("No warnings generated (expected for test data)")
		}

		// Should have no errors (data is valid, just small)
		if len(validationResult.Errors) > 0 {
			t.Errorf("Unexpected validation errors: %v", validationResult.Errors)
		}
	})

	// Step 4: Write output
	w, err := writer.NewWriter(cfg)
	if err != nil {
		t.Fatalf("NewWriter failed: %v", err)
	}
	if err := w.WriteAll(convertedData); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Verify output files
	t.Run("output", func(t *testing.T) {
		expectedFiles := []string{
			writer.FileSolarSystems,
			writer.FileRegions,
			writer.FileConstellations,
			writer.FileWormholeClasses,
			writer.FileShipTypes,
			writer.FileItemGroups,
			writer.FileSystemJumps,
		}

		for _, filename := range expectedFiles {
			path := filepath.Join(outputDir, filename)
			info, err := os.Stat(path)
			if err != nil {
				t.Errorf("Output file %s not created: %v", filename, err)
				continue
			}
			if info.Size() == 0 {
				t.Errorf("Output file %s is empty", filename)
			}
		}
	})
}

func TestIntegration_PassthroughFiles(t *testing.T) {
	// Create source directory with passthrough files
	srcDir, err := os.MkdirTemp("", "passthrough_src")
	if err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(srcDir) }()

	// Create output directory
	outputDir, err := os.MkdirTemp("", "passthrough_dst")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	// Create passthrough files
	passthroughFiles := []string{
		"wormholes.json",
		"wormholeClasses.json",
		"effects.json",
		"triglavianSystems.json",
	}

	for _, filename := range passthroughFiles {
		content := []byte(`{"test": "` + filename + `"}`)
		if err := os.WriteFile(filepath.Join(srcDir, filename), content, 0644); err != nil {
			t.Fatalf("failed to create %s: %v", filename, err)
		}
	}

	cfg := &config.Config{
		OutputDir:    outputDir,
		OutputFormat: config.FormatJSON,
		PrettyPrint:  true,
		Verbose:      false,
	}

	w := writer.NewJSONWriter(cfg)
	if err := w.CopyPassthroughFiles(srcDir); err != nil {
		t.Fatalf("CopyPassthroughFiles failed: %v", err)
	}

	// Verify files were copied
	for _, filename := range passthroughFiles {
		dstPath := filepath.Join(outputDir, filename)
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Errorf("Passthrough file %s not copied", filename)
		}
	}
}

func TestIntegration_SortingConsistency(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	cfg := &config.Config{
		SDEPath: sdeDir,
		Verbose: false,
	}

	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify all data is sorted correctly
	t.Run("regions sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.Universe.Regions); i++ {
			if convertedData.Universe.Regions[i-1].RegionID >= convertedData.Universe.Regions[i].RegionID {
				t.Errorf("Regions not sorted: %d >= %d",
					convertedData.Universe.Regions[i-1].RegionID,
					convertedData.Universe.Regions[i].RegionID)
			}
		}
	})

	t.Run("constellations sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.Universe.Constellations); i++ {
			if convertedData.Universe.Constellations[i-1].ConstellationID >= convertedData.Universe.Constellations[i].ConstellationID {
				t.Errorf("Constellations not sorted: %d >= %d",
					convertedData.Universe.Constellations[i-1].ConstellationID,
					convertedData.Universe.Constellations[i].ConstellationID)
			}
		}
	})

	t.Run("solar systems sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.Universe.SolarSystems); i++ {
			if convertedData.Universe.SolarSystems[i-1].SolarSystemID >= convertedData.Universe.SolarSystems[i].SolarSystemID {
				t.Errorf("Solar systems not sorted: %d >= %d",
					convertedData.Universe.SolarSystems[i-1].SolarSystemID,
					convertedData.Universe.SolarSystems[i].SolarSystemID)
			}
		}
	})

	t.Run("types sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.InvTypes); i++ {
			if convertedData.InvTypes[i-1].TypeID >= convertedData.InvTypes[i].TypeID {
				t.Errorf("Types not sorted: %d >= %d",
					convertedData.InvTypes[i-1].TypeID,
					convertedData.InvTypes[i].TypeID)
			}
		}
	})

	t.Run("groups sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.InvGroups); i++ {
			if convertedData.InvGroups[i-1].GroupID >= convertedData.InvGroups[i].GroupID {
				t.Errorf("Groups not sorted: %d >= %d",
					convertedData.InvGroups[i-1].GroupID,
					convertedData.InvGroups[i].GroupID)
			}
		}
	})

	t.Run("system jumps sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.SystemJumps); i++ {
			prev := convertedData.SystemJumps[i-1]
			curr := convertedData.SystemJumps[i]
			if prev.FromSolarSystemID > curr.FromSolarSystemID {
				t.Errorf("System jumps not sorted by FromSolarSystemID: %d > %d",
					prev.FromSolarSystemID, curr.FromSolarSystemID)
			}
			if prev.FromSolarSystemID == curr.FromSolarSystemID &&
				prev.ToSolarSystemID >= curr.ToSolarSystemID {
				t.Errorf("System jumps not sorted by ToSolarSystemID: %d >= %d",
					prev.ToSolarSystemID, curr.ToSolarSystemID)
			}
		}
	})
}

func TestIntegration_DataIntegrity(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	cfg := &config.Config{
		SDEPath: sdeDir,
		Verbose: false,
	}

	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify referential integrity
	t.Run("constellation region references", func(t *testing.T) {
		regionIDs := make(map[int64]bool)
		for _, r := range convertedData.Universe.Regions {
			regionIDs[r.RegionID] = true
		}

		for _, c := range convertedData.Universe.Constellations {
			if !regionIDs[c.RegionID] {
				t.Errorf("Constellation %d references non-existent region %d",
					c.ConstellationID, c.RegionID)
			}
		}
	})

	t.Run("solar system references", func(t *testing.T) {
		regionIDs := make(map[int64]bool)
		for _, r := range convertedData.Universe.Regions {
			regionIDs[r.RegionID] = true
		}

		constellationIDs := make(map[int64]bool)
		for _, c := range convertedData.Universe.Constellations {
			constellationIDs[c.ConstellationID] = true
		}

		for _, s := range convertedData.Universe.SolarSystems {
			if !regionIDs[s.RegionID] {
				t.Errorf("Solar system %d references non-existent region %d",
					s.SolarSystemID, s.RegionID)
			}
			if !constellationIDs[s.ConstellationID] {
				t.Errorf("Solar system %d references non-existent constellation %d",
					s.SolarSystemID, s.ConstellationID)
			}
		}
	})

	t.Run("system jump references", func(t *testing.T) {
		systemIDs := make(map[int64]bool)
		for _, s := range convertedData.Universe.SolarSystems {
			systemIDs[s.SolarSystemID] = true
		}

		for _, j := range convertedData.SystemJumps {
			if !systemIDs[j.FromSolarSystemID] {
				t.Errorf("System jump references non-existent from system %d",
					j.FromSolarSystemID)
			}
			if !systemIDs[j.ToSolarSystemID] {
				t.Errorf("System jump references non-existent to system %d",
					j.ToSolarSystemID)
			}
		}
	})
}

// TestIntegration_CSVOutput tests the full CSV output pipeline end-to-end.
func TestIntegration_CSVOutput(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	outputDir, err := os.MkdirTemp("", "csv_integration_test")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	cfg := &config.Config{
		SDEPath:      sdeDir,
		OutputDir:    outputDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	// Parse
	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	// Transform
	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Write CSV
	w := writer.NewCSVWriter(cfg)
	if err := w.WriteAll(convertedData); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Verify all expected CSV files exist
	expectedFiles := writer.GetOutputFiles(config.FormatCSV)
	for _, filename := range expectedFiles {
		path := filepath.Join(outputDir, filename)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("CSV file %s not created: %v", filename, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("CSV file %s is empty", filename)
		}
	}
}

// TestIntegration_CSVFuzzworkFormat validates that CSV output matches Fuzzwork format expectations.
func TestIntegration_CSVFuzzworkFormat(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	outputDir, err := os.MkdirTemp("", "csv_fuzzwork_test")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	cfg := &config.Config{
		SDEPath:      sdeDir,
		OutputDir:    outputDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	// Parse and transform
	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Write CSV
	w := writer.NewCSVWriter(cfg)
	if err := w.WriteAll(convertedData); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Test each CSV file for Wanderer format compliance
	t.Run("mapSolarSystems.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileSolarSystems, "mapSolarSystems")
		validateCSVRowCount(t, outputDir, writer.CSVFileSolarSystems, len(convertedData.Universe.SolarSystems))

		// Validate specific row content
		// Slimmed headers: solarSystemID(0), solarSystemName(1), regionID(2), constellationID(3), security(4), sunTypeID(5)
		records := readCSVFile(t, outputDir, writer.CSVFileSolarSystems)
		if len(records) > 1 {
			row := records[1] // First data row
			// Check that we have the expected number of columns
			if len(row) != 6 {
				t.Errorf("expected 6 columns, got %d", len(row))
			}
		}
	})

	t.Run("mapRegions.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileRegions, "mapRegions")
		validateCSVRowCount(t, outputDir, writer.CSVFileRegions, len(convertedData.Universe.Regions))
	})

	t.Run("mapConstellations.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileConstellations, "mapConstellations")
		validateCSVRowCount(t, outputDir, writer.CSVFileConstellations, len(convertedData.Universe.Constellations))
	})

	t.Run("invTypes.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileTypes, "invTypes")
		validateCSVRowCount(t, outputDir, writer.CSVFileTypes, len(convertedData.InvTypes))

		// Slimmed headers: typeID(0), groupID(1), typeName(2), mass(3), volume(4), capacity(5)
		records := readCSVFile(t, outputDir, writer.CSVFileTypes)
		if len(records) > 1 {
			row := records[1]
			if len(row) != 6 {
				t.Errorf("expected 6 columns, got %d", len(row))
			}
		}
	})

	t.Run("invGroups.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileGroups, "invGroups")
		validateCSVRowCount(t, outputDir, writer.CSVFileGroups, len(convertedData.InvGroups))

		// Slimmed headers: groupID(0), categoryID(1), groupName(2)
		records := readCSVFile(t, outputDir, writer.CSVFileGroups)
		if len(records) > 1 {
			row := records[1]
			if len(row) != 3 {
				t.Errorf("expected 3 columns, got %d", len(row))
			}
		}
	})

	t.Run("mapLocationWormholeClasses.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileWormholeClasses, "mapLocationWormholeClasses")
		validateCSVRowCount(t, outputDir, writer.CSVFileWormholeClasses, len(convertedData.WormholeClasses))
	})

	t.Run("mapSolarSystemJumps.csv", func(t *testing.T) {
		validateCSVHeaders(t, outputDir, writer.CSVFileSystemJumps, "mapSolarSystemJumps")
		// SystemJumps are bidirectional, so count is 2x the number of connections
		expectedJumps := len(convertedData.SystemJumps)
		validateCSVRowCount(t, outputDir, writer.CSVFileSystemJumps, expectedJumps)

		// Validate jump data contains all required fields
		records := readCSVFile(t, outputDir, writer.CSVFileSystemJumps)
		if len(records) > 1 {
			row := records[1]
			// All fields should be valid integers
			for i, colName := range models.CSVHeaders["mapSolarSystemJumps"] {
				if _, err := strconv.ParseInt(row[i], 10, 64); err != nil {
					t.Errorf("column %s should be integer, got '%s'", colName, row[i])
				}
			}
		}
	})
}

// TestIntegration_CSVNullHandling tests that null/None values are formatted correctly.
func TestIntegration_CSVNullHandling(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	outputDir, err := os.MkdirTemp("", "csv_null_test")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	cfg := &config.Config{
		SDEPath:      sdeDir,
		OutputDir:    outputDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	w := writer.NewCSVWriter(cfg)
	if err := w.WriteAll(convertedData); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Check that optional fields are formatted as "None" when nil
	t.Run("solar system nullable fields", func(t *testing.T) {
		records := readCSVFile(t, outputDir, writer.CSVFileSolarSystems)
		for i, row := range records {
			if i == 0 {
				continue // Skip header
			}
			// sunTypeID (index 5) can be "None" in the slimmed format
			sunTypeID := row[5]

			// If the value is not "None", it should be a valid integer
			if sunTypeID != "None" {
				if _, err := strconv.ParseInt(sunTypeID, 10, 64); err != nil {
					t.Errorf("row %d: sunTypeID should be 'None' or integer, got '%s'", i, sunTypeID)
				}
			}
		}
	})

	t.Run("type fields", func(t *testing.T) {
		records := readCSVFile(t, outputDir, writer.CSVFileTypes)
		for i, row := range records {
			if i == 0 {
				continue // Skip header
			}
			// Slimmed format has no nullable fields - all are required:
			// typeID(0), groupID(1), typeName(2), mass(3), volume(4), capacity(5)
			if len(row) != 6 {
				t.Errorf("row %d: expected 6 columns, got %d", i, len(row))
			}
		}
	})
}

// TestIntegration_CSVColumnOrdering tests that columns are in the exact Fuzzwork order.
func TestIntegration_CSVColumnOrdering(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	outputDir, err := os.MkdirTemp("", "csv_ordering_test")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	cfg := &config.Config{
		SDEPath:      sdeDir,
		OutputDir:    outputDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	w := writer.NewCSVWriter(cfg)
	if err := w.WriteAll(convertedData); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Verify exact column order for each file
	testCases := []struct {
		filename  string
		headerKey string
	}{
		{writer.CSVFileSolarSystems, "mapSolarSystems"},
		{writer.CSVFileRegions, "mapRegions"},
		{writer.CSVFileConstellations, "mapConstellations"},
		{writer.CSVFileTypes, "invTypes"},
		{writer.CSVFileGroups, "invGroups"},
		{writer.CSVFileWormholeClasses, "mapLocationWormholeClasses"},
		{writer.CSVFileSystemJumps, "mapSolarSystemJumps"},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			records := readCSVFile(t, outputDir, tc.filename)
			if len(records) == 0 {
				t.Fatal("empty CSV file")
			}

			headers := records[0]
			expectedHeaders := models.CSVHeaders[tc.headerKey]

			if len(headers) != len(expectedHeaders) {
				t.Errorf("expected %d columns, got %d", len(expectedHeaders), len(headers))
				return
			}

			for i, expected := range expectedHeaders {
				if headers[i] != expected {
					t.Errorf("column %d: expected '%s', got '%s'", i, expected, headers[i])
				}
			}
		})
	}
}

// Helper functions for CSV validation

func readCSVFile(t *testing.T, dir, filename string) [][]string {
	t.Helper()
	path := filepath.Join(dir, filename)
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open %s: %v", filename, err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV %s: %v", filename, err)
	}
	return records
}

func validateCSVHeaders(t *testing.T, dir, filename, headerKey string) {
	t.Helper()
	records := readCSVFile(t, dir, filename)
	if len(records) == 0 {
		t.Fatalf("file %s is empty", filename)
	}

	expectedHeaders := models.CSVHeaders[headerKey]
	if len(records[0]) != len(expectedHeaders) {
		t.Errorf("file %s: expected %d columns, got %d", filename, len(expectedHeaders), len(records[0]))
	}

	for i, expected := range expectedHeaders {
		if i >= len(records[0]) {
			break
		}
		if records[0][i] != expected {
			t.Errorf("file %s column %d: expected '%s', got '%s'", filename, i, expected, records[0][i])
		}
	}
}

func validateCSVRowCount(t *testing.T, dir, filename string, expectedDataRows int) {
	t.Helper()
	records := readCSVFile(t, dir, filename)
	// +1 for header row
	expectedTotal := expectedDataRows + 1
	if len(records) != expectedTotal {
		t.Errorf("file %s: expected %d rows (including header), got %d", filename, expectedTotal, len(records))
	}
}
