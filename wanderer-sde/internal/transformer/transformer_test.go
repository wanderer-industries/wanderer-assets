package transformer

import (
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/internal/parser"
)

func TestTransformer_Transform(t *testing.T) {
	cfg := &config.Config{Verbose: false}
	tr := New(cfg)

	sunType6 := int64(6)
	sunType7 := int64(7)
	parseResult := &parser.ParseResult{
		Regions: []models.Region{
			{RegionID: 10000002, RegionName: "The Forge"},
			{RegionID: 10000001, RegionName: "Derelik"},
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
				SunTypeID:       &sunType6,
				Security:        0.9459, // Raw security
				Constellation:   "None",
			},
			{
				SolarSystemID:   30000001,
				RegionID:        10000001,
				ConstellationID: 20000001,
				SolarSystemName: "Tanoo",
				SunTypeID:       &sunType7,
				Security:        0.047, // Edge case: very low positive
				Constellation:   "None",
			},
		},
		Types: map[int64]models.SDEType{
			587:  {GroupID: 25, Name: map[string]string{"en": "Rifter"}, Mass: 1350000, Volume: 27500, Published: true},
			588:  {GroupID: 25, Name: map[string]string{"en": "Slasher"}, Mass: 1200000, Volume: 26000, Published: true},
			2456: {GroupID: 18, Name: map[string]string{"en": "Hobgoblin I"}, Mass: 2500, Volume: 5, Published: true},
		},
		Groups: map[int64]models.SDEGroup{
			25: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Frigate"}, Published: true},
			18: {CategoryID: 7, Name: map[string]string{"en": "Drone"}, Published: true},
		},
		Categories: map[int64]models.SDECategory{
			6: {Name: map[string]string{"en": "Ship"}, Published: true},
			7: {Name: map[string]string{"en": "Drone"}, Published: true},
		},
		WormholeClasses: []models.WormholeClassLocation{
			{LocationID: 31000001, WormholeClassID: 3},
		},
		SystemJumps: []models.SystemJump{
			{FromSolarSystemID: 30000142, ToSolarSystemID: 30000001},
		},
		NPCStations: map[int64]models.SDENPCStation{
			60012736: {OwnerID: 1000137, SolarSystemID: 30000142, TypeID: 2501},
		},
		NPCCorporations: map[int64]models.SDENPCCorporation{
			1000137: {Name: map[string]string{"en": "DED"}},
		},
	}

	result, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify universe data
	if result.Universe == nil {
		t.Fatal("Universe data is nil")
	}

	// Check regions are sorted
	if len(result.Universe.Regions) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(result.Universe.Regions))
	}
	if result.Universe.Regions[0].RegionID != 10000001 {
		t.Error("Regions not sorted by ID")
	}

	// Check solar systems are transformed with correct security
	if len(result.Universe.SolarSystems) != 2 {
		t.Errorf("Expected 2 solar systems, got %d", len(result.Universe.SolarSystems))
	}

	// Find Jita and verify security transformation
	for _, sys := range result.Universe.SolarSystems {
		if sys.SolarSystemName == "Jita" {
			// Raw security value should be preserved (not transformed)
			// Wanderer calculates true security itself from the raw value
			if sys.Security != 0.9459 {
				t.Errorf("Expected Jita raw security 0.9459, got %f", sys.Security)
			}
		}
		if sys.SolarSystemName == "Tanoo" {
			// Raw security value should be preserved
			if sys.Security != 0.047 {
				t.Errorf("Expected Tanoo raw security 0.047, got %f", sys.Security)
			}
		}
	}

	// Check types (ships only - groupID 25 has categoryID 6)
	if len(result.InvTypes) != 2 {
		t.Errorf("Expected 2 ship types, got %d", len(result.InvTypes))
	}

	// Check groups (ships only - categoryID 6)
	if len(result.InvGroups) != 1 {
		t.Errorf("Expected 1 ship group, got %d", len(result.InvGroups))
	}

	// Check wormhole classes
	if len(result.WormholeClasses) != 1 {
		t.Errorf("Expected 1 wormhole class, got %d", len(result.WormholeClasses))
	}

	// Check system jumps
	if len(result.SystemJumps) != 1 {
		t.Errorf("Expected 1 system jump, got %d", len(result.SystemJumps))
	}
}

func TestTransformer_Validate(t *testing.T) {
	cfg := &config.Config{Verbose: false}
	tr := New(cfg)

	tests := []struct {
		name           string
		data           *models.ConvertedData
		expectErrors   bool
		expectWarnings bool
	}{
		{
			name: "valid data",
			data: &models.ConvertedData{
				Universe: &models.UniverseData{
					Regions:        make([]models.Region, 110),
					Constellations: make([]models.Constellation, 1100),
					SolarSystems:   make([]models.SolarSystem, 8500),
				},
				InvTypes:        make([]models.InvType, 35000),
				InvGroups:       make([]models.InvGroup, 50),
				WormholeClasses: make([]models.WormholeClassLocation, 800), // Regions + constellations + systems, expected ~803
				SystemJumps:     make([]models.SystemJump, 14000),          // Bidirectional jumps, expected ~13,776
				NPCStations:     make([]models.NPCStation, 45),             // DED stations, expected ~45
			},
			expectErrors:   false,
			expectWarnings: false,
		},
		{
			name: "empty solar systems",
			data: &models.ConvertedData{
				Universe: &models.UniverseData{
					Regions:        make([]models.Region, 110),
					Constellations: make([]models.Constellation, 1100),
					SolarSystems:   []models.SolarSystem{},
				},
				InvTypes:    make([]models.InvType, 35000),
				SystemJumps: make([]models.SystemJump, 11000),
			},
			expectErrors:   true,
			expectWarnings: true, // Empty solar systems also triggers minimum count warning
		},
		{
			name: "below minimum counts",
			data: &models.ConvertedData{
				Universe: &models.UniverseData{
					Regions:        make([]models.Region, 10),        // Below minimum
					Constellations: make([]models.Constellation, 10), // Below minimum
					SolarSystems:   make([]models.SolarSystem, 100),  // Below minimum
				},
				InvTypes:    make([]models.InvType, 50), // Below minimum
				SystemJumps: make([]models.SystemJump, 100),
			},
			expectErrors:   false, // These are warnings, not errors
			expectWarnings: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.Validate(tt.data)

			if tt.expectErrors && len(result.Errors) == 0 {
				t.Error("Expected errors but got none")
			}
			if !tt.expectErrors && len(result.Errors) > 0 {
				t.Errorf("Unexpected errors: %v", result.Errors)
			}
			if tt.expectWarnings && len(result.Warnings) == 0 {
				t.Error("Expected warnings but got none")
			}
			if !tt.expectWarnings && len(result.Warnings) > 0 {
				t.Errorf("Unexpected warnings: %v", result.Warnings)
			}
		})
	}
}

func TestTransformer_SortFunctions(t *testing.T) {
	cfg := &config.Config{Verbose: false}
	tr := New(cfg)

	// Test sortRegions
	t.Run("sortRegions", func(t *testing.T) {
		regions := []models.Region{
			{RegionID: 3},
			{RegionID: 1},
			{RegionID: 2},
		}
		sorted := tr.sortRegions(regions)

		for i := 1; i < len(sorted); i++ {
			if sorted[i-1].RegionID >= sorted[i].RegionID {
				t.Errorf("Regions not sorted: %d >= %d", sorted[i-1].RegionID, sorted[i].RegionID)
			}
		}
	})

	// Test sortConstellations
	t.Run("sortConstellations", func(t *testing.T) {
		constellations := []models.Constellation{
			{ConstellationID: 3},
			{ConstellationID: 1},
			{ConstellationID: 2},
		}
		sorted := tr.sortConstellations(constellations)

		for i := 1; i < len(sorted); i++ {
			if sorted[i-1].ConstellationID >= sorted[i].ConstellationID {
				t.Errorf("Constellations not sorted: %d >= %d", sorted[i-1].ConstellationID, sorted[i].ConstellationID)
			}
		}
	})

	// Test sortWormholeClasses
	t.Run("sortWormholeClasses", func(t *testing.T) {
		classes := []models.WormholeClassLocation{
			{LocationID: 3},
			{LocationID: 1},
			{LocationID: 2},
		}
		sorted := tr.sortWormholeClasses(classes)

		for i := 1; i < len(sorted); i++ {
			if sorted[i-1].LocationID >= sorted[i].LocationID {
				t.Errorf("Wormhole classes not sorted: %d >= %d", sorted[i-1].LocationID, sorted[i].LocationID)
			}
		}
	})

	// Test transformSystemJumps (renamed from sortSystemJumps)
	t.Run("transformSystemJumps", func(t *testing.T) {
		systems := []models.SolarSystem{
			{SolarSystemID: 1, RegionID: 100, ConstellationID: 1000},
			{SolarSystemID: 2, RegionID: 100, ConstellationID: 1000},
			{SolarSystemID: 3, RegionID: 200, ConstellationID: 2000},
		}
		jumps := []models.SystemJump{
			{FromSolarSystemID: 2, ToSolarSystemID: 3},
			{FromSolarSystemID: 1, ToSolarSystemID: 2},
			{FromSolarSystemID: 1, ToSolarSystemID: 1},
		}
		sorted := tr.transformSystemJumps(jumps, systems)

		// Verify sorting by FromSolarSystemID then ToSolarSystemID
		for i := 1; i < len(sorted); i++ {
			if sorted[i-1].FromSolarSystemID > sorted[i].FromSolarSystemID {
				t.Errorf("System jumps not sorted by FromSolarSystemID")
			}
			if sorted[i-1].FromSolarSystemID == sorted[i].FromSolarSystemID &&
				sorted[i-1].ToSolarSystemID >= sorted[i].ToSolarSystemID {
				t.Errorf("System jumps not sorted by ToSolarSystemID within same FromSolarSystemID")
			}
		}

		// Verify region/constellation enrichment
		for _, jump := range sorted {
			if jump.FromSolarSystemID == 1 {
				if jump.FromRegionID != 100 || jump.FromConstellationID != 1000 {
					t.Errorf("Jump from system 1 should have RegionID 100 and ConstellationID 1000")
				}
			}
		}
	})
}

func TestTransformer_TransformSolarSystems(t *testing.T) {
	cfg := &config.Config{Verbose: false}
	tr := New(cfg)

	systems := []models.SolarSystem{
		{SolarSystemID: 3, SolarSystemName: "C", Security: 0.9459},
		{SolarSystemID: 1, SolarSystemName: "A", Security: -0.5},
		{SolarSystemID: 2, SolarSystemName: "B", Security: 0.047},
	}

	result := tr.transformSolarSystems(systems)

	// Verify sorting
	for i := 1; i < len(result); i++ {
		if result[i-1].SolarSystemID >= result[i].SolarSystemID {
			t.Errorf("Systems not sorted: %d >= %d", result[i-1].SolarSystemID, result[i].SolarSystemID)
		}
	}

	// Verify security values are preserved (not transformed)
	// Wanderer calculates true security itself from raw values
	for _, sys := range result {
		switch sys.SolarSystemID {
		case 1:
			if sys.Security != -0.5 {
				t.Errorf("System 1 security should be -0.5, got %f", sys.Security)
			}
		case 2:
			// Raw value should be preserved
			if sys.Security != 0.047 {
				t.Errorf("System 2 security should be 0.047 (raw), got %f", sys.Security)
			}
		case 3:
			// Raw value should be preserved
			if sys.Security != 0.9459 {
				t.Errorf("System 3 security should be 0.9459 (raw), got %f", sys.Security)
			}
		}
	}
}

func TestValidationResult_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		result   models.ValidationResult
		expected bool
	}{
		{
			name:     "valid - no errors",
			result:   models.ValidationResult{Errors: []string{}},
			expected: true,
		},
		{
			name:     "valid - nil errors",
			result:   models.ValidationResult{Errors: nil},
			expected: true,
		},
		{
			name:     "invalid - has errors",
			result:   models.ValidationResult{Errors: []string{"error1"}},
			expected: false,
		},
		{
			name: "valid - warnings only",
			result: models.ValidationResult{
				Errors:   []string{},
				Warnings: []string{"warning1"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}
