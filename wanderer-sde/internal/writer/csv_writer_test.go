package writer

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

func TestCSVWriter_WriteAll(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "csv_writer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	// Create test data
	sunTypeID := int64(6)
	factionID := int64(500001)
	data := &models.ConvertedData{
		Universe: &models.UniverseData{
			Regions: []models.Region{
				{
					RegionID:   10000002,
					RegionName: "The Forge",
					X:          -96538765063520384,
					Y:          60376779980673024,
					Z:          112361271557498880,
					XMin:       -119406988656885760,
					XMax:       -73670541470154752,
					YMin:       42907838618230784,
					YMax:       77845721343115264,
					ZMin:       97540122286497792,
					ZMax:       127182420828499968,
					FactionID:  &factionID,
					Nebula:     0,
					Radius:     0,
				},
			},
			Constellations: []models.Constellation{
				{
					ConstellationID:   20000020,
					ConstellationName: "Kimotoro",
					RegionID:          10000002,
					X:                 -90032979768340480,
					Y:                 44587085903659008,
					Z:                 115019533039984640,
					XMin:              -96605073106452480,
					XMax:              -83460886430228480,
					YMin:              42907838618230784,
					YMax:              46266333189087232,
					ZMin:              109413671420215296,
					ZMax:              120625394659753984,
					FactionID:         &factionID,
					Radius:            0,
				},
			},
			SolarSystems: []models.SolarSystem{
				{
					SolarSystemID:   30000142,
					RegionID:        10000002,
					ConstellationID: 20000020,
					SolarSystemName: "Jita",
					X:               -129483428307508224,
					Y:               60654338355159040,
					Z:               -79015917355671552,
					XMin:            -129483428307508224,
					XMax:            -129483428307508224,
					YMin:            60654338355159040,
					YMax:            60654338355159040,
					ZMin:            -79015917355671552,
					ZMax:            -79015917355671552,
					Luminosity:      0.01575,
					Border:          false,
					Fringe:          false,
					Corridor:        false,
					Hub:             true,
					International:   false,
					Regional:        false,
					Constellation:   "None",
					Security:        0.9459131166648389,
					FactionID:       &factionID,
					Radius:          2838370304000,
					SunTypeID:       &sunTypeID,
					SecurityClass:   "B",
				},
			},
		},
		InvTypes: []models.InvType{
			{
				TypeID:      587,
				GroupID:     25,
				TypeName:    "Rifter",
				Description: "The Rifter is a very fast frigate.",
				Mass:        1350000,
				Volume:      27500,
				Capacity:    125,
				PortionSize: 1,
				RaceID:      nil,
				BasePrice:   0,
				Published:   true,
				MarketGroupID: func() *int64 {
					v := int64(64)
					return &v
				}(),
				IconID:    nil,
				SoundID:   nil,
				GraphicID: nil,
			},
		},
		InvGroups: []models.InvGroup{
			{
				GroupID:              25,
				CategoryID:           6,
				GroupName:            "Frigate",
				IconID:               nil,
				UseBasePrice:         false,
				Anchored:             false,
				Anchorable:           false,
				FittableNonSingleton: true,
				Published:            true,
			},
		},
		WormholeClasses: []models.WormholeClassLocation{
			{LocationID: 10000002, WormholeClassID: 7},
		},
		SystemJumps: []models.SystemJump{
			{
				FromRegionID:        10000002,
				FromConstellationID: 20000020,
				FromSolarSystemID:   30000142,
				ToSolarSystemID:     30000144,
				ToConstellationID:   20000020,
				ToRegionID:          10000002,
			},
		},
		NPCStations: []models.NPCStation{
			{
				StationID:     60012736,
				SolarSystemID: 30002187,
				OwnerID:       1000125,
				OwnerName:     "DED",
				TypeID:        1529,
			},
		},
	}

	// Write all files
	if err := w.WriteAll(data); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Verify files exist and have correct headers
	tests := []struct {
		filename    string
		headerKey   string
		expectedLen int
	}{
		{CSVFileSolarSystems, "mapSolarSystems", 1},
		{CSVFileRegions, "mapRegions", 1},
		{CSVFileConstellations, "mapConstellations", 1},
		{CSVFileTypes, "invTypes", 1},
		{CSVFileGroups, "invGroups", 1},
		{CSVFileWormholeClasses, "mapLocationWormholeClasses", 1},
		{CSVFileSystemJumps, "mapSolarSystemJumps", 1},
		{CSVFileNPCStations, "npcStations", 1},
	}

	for _, tt := range tests {
		path := filepath.Join(tmpDir, tt.filename)

		// Check file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("file %s was not created", tt.filename)
			continue
		}

		// Read and parse file
		file, err := os.Open(path)
		if err != nil {
			t.Errorf("failed to open %s: %v", tt.filename, err)
			continue
		}

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		_ = file.Close()
		if err != nil {
			t.Errorf("file %s contains invalid CSV: %v", tt.filename, err)
			continue
		}

		// Check header row
		if len(records) < 1 {
			t.Errorf("file %s is empty", tt.filename)
			continue
		}

		expectedHeaders := models.CSVHeaders[tt.headerKey]
		if len(records[0]) != len(expectedHeaders) {
			t.Errorf("file %s has %d columns, expected %d", tt.filename, len(records[0]), len(expectedHeaders))
			continue
		}

		for i, header := range expectedHeaders {
			if records[0][i] != header {
				t.Errorf("file %s column %d: got %s, expected %s", tt.filename, i, records[0][i], header)
			}
		}

		// Check data row count (header + data)
		if len(records) != tt.expectedLen+1 {
			t.Errorf("file %s has %d rows, expected %d (including header)", tt.filename, len(records), tt.expectedLen+1)
		}
	}
}

func TestCSVWriter_SolarSystemRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_ss_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	sunTypeID := int64(6)
	factionID := int64(500001)
	systems := []models.SolarSystem{
		{
			SolarSystemID:   30000142,
			RegionID:        10000002,
			ConstellationID: 20000020,
			SolarSystemName: "Jita",
			Security:        0.9459131166648389,
			FactionID:       &factionID,
			SunTypeID:       &sunTypeID,
		},
	}

	if err := w.WriteSolarSystems(systems); err != nil {
		t.Fatalf("WriteSolarSystems failed: %v", err)
	}

	// Read and verify content
	file, err := os.Open(filepath.Join(tmpDir, CSVFileSolarSystems))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + data), got %d", len(records))
	}

	row := records[1]

	// Slimmed down headers: solarSystemID(0), solarSystemName(1), regionID(2),
	// constellationID(3), security(4), sunTypeID(5)
	tests := []struct {
		index    int
		name     string
		expected string
	}{
		{0, "solarSystemID", "30000142"},
		{1, "solarSystemName", "Jita"},
		{2, "regionID", "10000002"},
		{3, "constellationID", "20000020"},
		{5, "sunTypeID", "6"},
	}

	for _, tt := range tests {
		if row[tt.index] != tt.expected {
			t.Errorf("%s: got %s, expected %s", tt.name, row[tt.index], tt.expected)
		}
	}
}

func TestCSVWriter_NullHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_null_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	// Create system with nil optional fields
	systems := []models.SolarSystem{
		{
			SolarSystemID:   30000142,
			RegionID:        10000002,
			ConstellationID: 20000020,
			SolarSystemName: "Test",
			Security:        0.5,
			SunTypeID:       nil, // Should output "None"
		},
	}

	if err := w.WriteSolarSystems(systems); err != nil {
		t.Fatalf("WriteSolarSystems failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileSolarSystems))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	row := records[1]

	// sunTypeID at index 5 should be "None"
	if row[5] != "None" {
		t.Errorf("sunTypeID: got %s, expected None", row[5])
	}
}

func TestCSVWriter_BoolFormatting(t *testing.T) {
	// Test that booleans are formatted as 0/1
	tests := []struct {
		value    bool
		expected string
	}{
		{true, "1"},
		{false, "0"},
	}

	for _, tt := range tests {
		result := models.FormatBool(tt.value)
		if result != tt.expected {
			t.Errorf("FormatBool(%v): got %s, expected %s", tt.value, result, tt.expected)
		}
	}
}

func TestCSVWriter_FloatFormatting(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{0.01575, "0.01575"},
		{1350000, "1350000"},
		{0.9459131166648389, "0.9459131166648389"},
		{-1234567890.5, "-1234567890.5"},
		{0, "0"},
		{-0.5, "-0.5"},
	}

	for _, tt := range tests {
		result := models.FormatFloat(tt.value)
		if result != tt.expected {
			t.Errorf("FormatFloat(%v): got %s, expected %s", tt.value, result, tt.expected)
		}
	}
}

func TestCSVWriter_CopyPassthroughFiles(t *testing.T) {
	// Create temp directories
	srcDir, err := os.MkdirTemp("", "passthrough_src")
	if err != nil {
		t.Fatalf("failed to create src temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(srcDir) }()

	dstDir, err := os.MkdirTemp("", "passthrough_dst")
	if err != nil {
		t.Fatalf("failed to create dst temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(dstDir) }()

	// Create some passthrough files in source
	testFiles := []string{"wormholes.json", "effects.json"}
	for _, filename := range testFiles {
		content := []byte(`{"test": true}`)
		if err := os.WriteFile(filepath.Join(srcDir, filename), content, 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	cfg := &config.Config{
		OutputDir:    dstDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	if err := w.CopyPassthroughFiles(srcDir); err != nil {
		t.Fatalf("CopyPassthroughFiles failed: %v", err)
	}

	// Verify copied files
	for _, filename := range testFiles {
		dstPath := filepath.Join(dstDir, filename)
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Errorf("file %s was not copied", filename)
		}
	}
}

func TestNewWriter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "writer_factory_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tests := []struct {
		format      config.OutputFormat
		expectError bool
		writerType  string
	}{
		{config.FormatCSV, false, "*writer.CSVWriter"},
		{config.FormatJSON, false, "*writer.JSONWriter"},
		{config.OutputFormat("invalid"), true, ""},
	}

	for _, tt := range tests {
		cfg := &config.Config{
			OutputDir:    tmpDir,
			OutputFormat: tt.format,
		}

		w, err := NewWriter(cfg)

		if tt.expectError {
			if err == nil {
				t.Errorf("format %s: expected error, got nil", tt.format)
			}
		} else {
			if err != nil {
				t.Errorf("format %s: unexpected error: %v", tt.format, err)
			}
			if w == nil {
				t.Errorf("format %s: writer is nil", tt.format)
			}
		}
	}
}

func TestGetOutputFiles(t *testing.T) {
	csvFiles := GetOutputFiles(config.FormatCSV)
	if len(csvFiles) != 8 {
		t.Errorf("expected 8 CSV files, got %d", len(csvFiles))
	}

	// Check that all CSV files have .csv extension
	for _, f := range csvFiles {
		if filepath.Ext(f) != ".csv" {
			t.Errorf("expected .csv extension, got %s", f)
		}
	}

	jsonFiles := GetOutputFiles(config.FormatJSON)
	if len(jsonFiles) != 8 {
		t.Errorf("expected 8 JSON files, got %d", len(jsonFiles))
	}

	// Check that all JSON files have .json extension
	for _, f := range jsonFiles {
		if filepath.Ext(f) != ".json" {
			t.Errorf("expected .json extension, got %s", f)
		}
	}
}

func TestCSVWriter_CoordinateFormatting(t *testing.T) {
	// Test that large coordinates are formatted without scientific notation
	// Note: Float64 has ~15-16 significant digits precision
	tests := []struct {
		value    float64
		expected string
	}{
		{-129483428307508220, "-129483428307508220"},
		{60654338355159040, "60654338355159040"},
		{-79015917355671550, "-79015917355671550"},
		{0, "0"},
		{-1, "-1"},
		{1.5, "1.5"},
		{0.123456789012345, "0.123456789012345"},
		{1000000000000, "1000000000000"},
		{-1000000000000, "-1000000000000"},
	}

	for _, tt := range tests {
		result := models.FormatFloat(tt.value)
		if result != tt.expected {
			t.Errorf("FormatFloat(%v): got %s, expected %s", tt.value, result, tt.expected)
		}
	}
}

func TestCSVWriter_SecurityFormatting(t *testing.T) {
	// Test security value formatting
	tests := []struct {
		value    float64
		expected string
	}{
		{0.9459131166648389, "0.9459131166648389"},
		{-1.0, "-1"},
		{0.5, "0.5"},
		{0.0, "0"},
		{0.1, "0.1"},
	}

	for _, tt := range tests {
		result := models.FormatSecurity(tt.value)
		if result != tt.expected {
			t.Errorf("FormatSecurity(%v): got %s, expected %s", tt.value, result, tt.expected)
		}
	}
}

func TestCSVWriter_RegionRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_region_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	regions := []models.Region{
		{
			RegionID:   10000002,
			RegionName: "The Forge",
		},
	}

	if err := w.WriteRegions(regions); err != nil {
		t.Fatalf("WriteRegions failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileRegions))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + data), got %d", len(records))
	}

	row := records[1]
	// Slimmed down: regionID(0), regionName(1)
	tests := []struct {
		index    int
		name     string
		expected string
	}{
		{0, "regionID", "10000002"},
		{1, "regionName", "The Forge"},
	}

	for _, tt := range tests {
		if row[tt.index] != tt.expected {
			t.Errorf("%s: got %s, expected %s", tt.name, row[tt.index], tt.expected)
		}
	}
}

func TestCSVWriter_ConstellationRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_constellation_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	constellations := []models.Constellation{
		{
			ConstellationID:   20000020,
			ConstellationName: "Kimotoro",
			RegionID:          10000002,
		},
	}

	if err := w.WriteConstellations(constellations); err != nil {
		t.Fatalf("WriteConstellations failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileConstellations))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + data), got %d", len(records))
	}

	row := records[1]
	// Slimmed down: constellationID(0), constellationName(1)
	tests := []struct {
		index    int
		name     string
		expected string
	}{
		{0, "constellationID", "20000020"},
		{1, "constellationName", "Kimotoro"},
	}

	for _, tt := range tests {
		if row[tt.index] != tt.expected {
			t.Errorf("%s: got %s, expected %s", tt.name, row[tt.index], tt.expected)
		}
	}
}

func TestCSVWriter_TypeRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_type_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	types := []models.InvType{
		{
			TypeID:   587,
			GroupID:  25,
			TypeName: "Rifter",
			Mass:     1350000,
			Volume:   27500,
			Capacity: 125,
		},
	}

	if err := w.WriteTypes(types); err != nil {
		t.Fatalf("WriteTypes failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileTypes))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + data), got %d", len(records))
	}

	row := records[1]
	// Slimmed down: typeID(0), groupID(1), typeName(2), mass(3), volume(4), capacity(5)
	tests := []struct {
		index    int
		name     string
		expected string
	}{
		{0, "typeID", "587"},
		{1, "groupID", "25"},
		{2, "typeName", "Rifter"},
		{3, "mass", "1350000"},
		{4, "volume", "27500"},
		{5, "capacity", "125"},
	}

	for _, tt := range tests {
		if row[tt.index] != tt.expected {
			t.Errorf("%s: got %s, expected %s", tt.name, row[tt.index], tt.expected)
		}
	}
}

func TestCSVWriter_GroupRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_group_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	groups := []models.InvGroup{
		{
			GroupID:    25,
			CategoryID: 6,
			GroupName:  "Frigate",
		},
	}

	if err := w.WriteGroups(groups); err != nil {
		t.Fatalf("WriteGroups failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileGroups))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + data), got %d", len(records))
	}

	row := records[1]
	// Slimmed down: groupID(0), categoryID(1), groupName(2)
	tests := []struct {
		index    int
		name     string
		expected string
	}{
		{0, "groupID", "25"},
		{1, "categoryID", "6"},
		{2, "groupName", "Frigate"},
	}

	for _, tt := range tests {
		if row[tt.index] != tt.expected {
			t.Errorf("%s: got %s, expected %s", tt.name, row[tt.index], tt.expected)
		}
	}
}

func TestCSVWriter_SystemJumpRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_jump_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	jumps := []models.SystemJump{
		{
			FromRegionID:        10000002,
			FromConstellationID: 20000020,
			FromSolarSystemID:   30000142,
			ToSolarSystemID:     30000144,
			ToConstellationID:   20000020,
			ToRegionID:          10000002,
		},
	}

	if err := w.WriteSystemJumps(jumps); err != nil {
		t.Fatalf("WriteSystemJumps failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileSystemJumps))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + data), got %d", len(records))
	}

	row := records[1]
	// Slimmed down: fromSolarSystemID(0), toSolarSystemID(1)
	tests := []struct {
		index    int
		name     string
		expected string
	}{
		{0, "fromSolarSystemID", "30000142"},
		{1, "toSolarSystemID", "30000144"},
	}

	for _, tt := range tests {
		if row[tt.index] != tt.expected {
			t.Errorf("%s: got %s, expected %s", tt.name, row[tt.index], tt.expected)
		}
	}
}

func TestCSVWriter_WormholeClassRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "csv_wh_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:    tmpDir,
		OutputFormat: config.FormatCSV,
		Verbose:      false,
	}

	w := NewCSVWriter(cfg)

	classes := []models.WormholeClassLocation{
		{LocationID: 11000001, WormholeClassID: 7},
		{LocationID: 31000001, WormholeClassID: 3},
	}

	if err := w.WriteWormholeClasses(classes); err != nil {
		t.Fatalf("WriteWormholeClasses failed: %v", err)
	}

	file, err := os.Open(filepath.Join(tmpDir, CSVFileWormholeClasses))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("expected 3 rows (header + 2 data), got %d", len(records))
	}

	// Check header
	if records[0][0] != "locationID" || records[0][1] != "wormholeClassID" {
		t.Errorf("unexpected headers: %v", records[0])
	}

	// Check first data row
	if records[1][0] != "11000001" || records[1][1] != "7" {
		t.Errorf("unexpected first row: %v", records[1])
	}

	// Check second data row
	if records[2][0] != "31000001" || records[2][1] != "3" {
		t.Errorf("unexpected second row: %v", records[2])
	}
}
