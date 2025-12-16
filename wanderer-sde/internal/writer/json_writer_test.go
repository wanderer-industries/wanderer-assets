package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

func TestJSONWriter_WriteAll(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "writer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{
		OutputDir:   tmpDir,
		PrettyPrint: true,
		Verbose:     false,
	}

	w := NewJSONWriter(cfg)

	// Create test data
	sunTypeID := int64(6)
	data := &models.ConvertedData{
		Universe: &models.UniverseData{
			Regions: []models.Region{
				{RegionID: 10000002, RegionName: "The Forge"},
			},
			Constellations: []models.Constellation{
				{ConstellationID: 20000020, ConstellationName: "Kimotoro", RegionID: 10000002},
			},
			SolarSystems: []models.SolarSystem{
				{
					SolarSystemID:   30000142,
					RegionID:        10000002,
					ConstellationID: 20000020,
					SolarSystemName: "Jita",
					SunTypeID:       &sunTypeID,
					Security:        0.9,
					Constellation:   "None",
				},
			},
		},
		InvTypes: []models.InvType{
			{TypeID: 587, GroupID: 25, TypeName: "Rifter", Mass: 1350000, Volume: 27500, Capacity: 125},
		},
		InvGroups: []models.InvGroup{
			{GroupID: 25, CategoryID: 6, GroupName: "Frigate"},
		},
		WormholeClasses: []models.WormholeClassLocation{
			{LocationID: 10000002, WormholeClassID: 7},
		},
		SystemJumps: []models.SystemJump{
			{FromSolarSystemID: 30000142, ToSolarSystemID: 30000144},
		},
		NPCStations: []models.NPCStation{
			{StationID: 60003760, SolarSystemID: 30000142, OwnerID: 1000035, OwnerName: "Caldari Navy", TypeID: 1529},
		},
	}

	// Write all files
	if err := w.WriteAll(data); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Verify files exist and have correct content
	tests := []struct {
		filename string
		expected interface{}
	}{
		{FileSolarSystems, data.Universe.SolarSystems},
		{FileRegions, data.Universe.Regions},
		{FileConstellations, data.Universe.Constellations},
		{FileShipTypes, data.InvTypes},
		{FileItemGroups, data.InvGroups},
		{FileWormholeClasses, data.WormholeClasses},
		{FileSystemJumps, data.SystemJumps},
		{FileNPCStations, data.NPCStations},
	}

	for _, tt := range tests {
		path := filepath.Join(tmpDir, tt.filename)

		// Check file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("file %s was not created", tt.filename)
			continue
		}

		// Read and parse file
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read %s: %v", tt.filename, err)
			continue
		}

		// Verify it's valid JSON
		var parsed interface{}
		if err := json.Unmarshal(content, &parsed); err != nil {
			t.Errorf("file %s contains invalid JSON: %v", tt.filename, err)
		}
	}
}

func TestJSONWriter_CopyPassthroughFiles(t *testing.T) {
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
		OutputDir:   dstDir,
		PrettyPrint: true,
		Verbose:     false,
	}

	w := NewJSONWriter(cfg)

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

func TestJSONWriter_PrettyPrint(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "writer_pretty_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tests := []struct {
		name        string
		prettyPrint bool
		expectMulti bool // expect multiple lines
	}{
		{"pretty", true, true},
		{"compact", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				t.Fatalf("failed to create output directory: %v", err)
			}

			cfg := &config.Config{
				OutputDir:   outDir,
				PrettyPrint: tt.prettyPrint,
				Verbose:     false,
			}

			w := NewJSONWriter(cfg)
			data := []models.Region{{RegionID: 1, RegionName: "Test"}}

			if err := w.WriteRegions(data); err != nil {
				t.Fatalf("WriteRegions failed: %v", err)
			}

			content, err := os.ReadFile(filepath.Join(outDir, FileRegions))
			if err != nil {
				t.Fatalf("failed to read output: %v", err)
			}

			// Count newlines in the output
			newlineCount := 0
			hasIndentedLine := false
			for i, b := range content {
				if b == '\n' {
					newlineCount++
					// Check if next character is whitespace (space or tab) indicating indentation
					if i+1 < len(content) && (content[i+1] == ' ' || content[i+1] == '\t') {
						hasIndentedLine = true
					}
				}
			}

			if tt.expectMulti {
				// Pretty-printed output must have multiple lines
				if newlineCount < 2 {
					t.Errorf("expected pretty-printed output to have multiple lines, got %d newlines", newlineCount)
				}
				// Pretty-printed output must have indentation (newline followed by whitespace)
				if !hasIndentedLine {
					t.Error("expected pretty-printed output to have indented lines (newline followed by whitespace)")
				}
			} else {
				// Compact output should be a single line (at most one trailing newline)
				if newlineCount > 1 {
					t.Errorf("expected compact output to be single-line, got %d newlines", newlineCount)
				}
				if hasIndentedLine {
					t.Error("expected compact output to have no indentation")
				}
			}
		})
	}
}
