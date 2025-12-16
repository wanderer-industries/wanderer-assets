package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
)

// createTestSDE creates a minimal SDE structure for testing
func createTestSDE(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "parser_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create mapRegions.yaml with new SDE format (2025+)
	regionsYAML := `10000002:
  name:
    en: "The Forge"
  position:
    x: -96538397329247680
    y: 68904722523889856
    z: 103886273221498080
  factionID: 500001
  nebulaID: 11806
  wormholeClassID: 7
10000001:
  name:
    en: "Derelik"
  position:
    x: -22292048624051248
    y: 75502104274173520
    z: 165604011665621600
  factionID: 500007
  nebulaID: 11799
  wormholeClassID: 7
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(regionsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapRegions.yaml: %v", err)
	}

	// Create mapConstellations.yaml with new SDE format (2025+)
	constellationsYAML := `20000020:
  regionID: 10000002
  name:
    en: "Kimotoro"
  position:
    x: -107314934797574880
    y: 65893634706137696
    z: 106631148888006560
  factionID: 500001
  radius: 12345678901234.5
  wormholeClassID: 7
20000001:
  regionID: 10000001
  name:
    en: "Joas"
  position:
    x: -26892617099116692
    y: 70892048824177140
    z: 153426871465893984
  radius: 9876543210987.6
  wormholeClassID: 7
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
  securityClass: B
  starID: 40000006
  position:
    x: -129500988494612512
    y: 60552325055663632
    z: 116970681498498304
  luminosity: 0.01575
  border: true
  fringe: false
  corridor: true
  hub: true
  international: false
  regional: true
  radius: 3.35e+12
30000001:
  constellationID: 20000001
  regionID: 10000001
  name:
    en: "Tanoo"
  securityStatus: 0.8576
  securityClass: B1
  starID: 40000007
  position:
    x: -23413628500306264
    y: 68096352988768256
    z: 149988050156945664
  luminosity: 0.02145
  border: false
  fringe: true
  corridor: false
  hub: false
  international: true
  regional: false
  radius: 2.89e+12
31000001:
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J123456"
  securityStatus: -1.0
  starID: 40045041
  wormholeClassID: 3
  position:
    x: 0
    y: 0
    z: 0
  luminosity: 0.0
  radius: 1.0e+12
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapSolarSystems.yaml"), []byte(systemsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapSolarSystems.yaml: %v", err)
	}

	// Create mapStargates.yaml
	stargatesYAML := `50000001:
  solarSystemID: 30000142
  destination:
    solarSystemID: 30000001
    stargateID: 50000002
  typeID: 16
50000002:
  solarSystemID: 30000001
  destination:
    solarSystemID: 30000142
    stargateID: 50000001
  typeID: 16
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapStargates.yaml"), []byte(stargatesYAML), 0644); err != nil {
		t.Fatalf("failed to create mapStargates.yaml: %v", err)
	}

	// Create types.yaml with all CSV-required fields
	typesYAML := `587:
  groupID: 25
  name:
    en: "Rifter"
  description:
    en: "The Rifter is a very powerful combat frigate."
  mass: 1350000.0
  volume: 27500.0
  capacity: 125.0
  portionSize: 1
  raceID: 2
  basePrice: 240000.0
  published: true
  marketGroupID: 64
  iconID: 588
  soundID: 20071
  graphicID: 588
588:
  groupID: 25
  name:
    en: "Slasher"
  description:
    en: "The Slasher is a quick and agile frigate."
  mass: 1200000.0
  volume: 26000.0
  capacity: 115.0
  portionSize: 1
  raceID: 2
  basePrice: 200000.0
  published: true
  marketGroupID: 64
  iconID: 589
  graphicID: 589
2456:
  groupID: 18
  name:
    en: "Hobgoblin I"
  description:
    en: "Light Scout Drone"
  mass: 2500.0
  volume: 5.0
  capacity: 0.0
  portionSize: 1
  published: true
  marketGroupID: 837
  graphicID: 2456
`
	if err := os.WriteFile(filepath.Join(tmpDir, "types.yaml"), []byte(typesYAML), 0644); err != nil {
		t.Fatalf("failed to create types.yaml: %v", err)
	}

	// Create groups.yaml with all CSV-required fields
	groupsYAML := `25:
  categoryID: 6
  name:
    en: "Frigate"
  iconID: 73
  useBasePrice: true
  anchored: false
  anchorable: false
  fittableNonSingleton: false
  published: true
18:
  categoryID: 7
  name:
    en: "Drone"
  iconID: 21
  useBasePrice: true
  anchored: false
  anchorable: false
  fittableNonSingleton: true
  published: true
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
40045041:
  solarSystemID: 31000001
  typeID: 45041
  radius: 111222333.0
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
  solarSystemID: 30000001
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

func TestParser_ParseRegions(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	regions, err := p.ParseRegions()
	if err != nil {
		t.Fatalf("ParseRegions failed: %v", err)
	}

	if len(regions) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(regions))
	}

	// Verify regions are sorted by ID
	for i := 1; i < len(regions); i++ {
		if regions[i-1].RegionID >= regions[i].RegionID {
			t.Errorf("Regions not sorted: %d >= %d", regions[i-1].RegionID, regions[i].RegionID)
		}
	}

	// Verify first region (should be Derelik with lowest ID)
	if regions[0].RegionID != 10000001 {
		t.Errorf("Expected first region ID to be 10000001, got %d", regions[0].RegionID)
	}
	if regions[0].RegionName != "Derelik" {
		t.Errorf("Expected first region name to be 'Derelik', got %q", regions[0].RegionName)
	}
}

func TestParser_ParseConstellations(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	constellations, err := p.ParseConstellations()
	if err != nil {
		t.Fatalf("ParseConstellations failed: %v", err)
	}

	if len(constellations) != 2 {
		t.Errorf("Expected 2 constellations, got %d", len(constellations))
	}

	// Verify constellations are sorted by ID
	for i := 1; i < len(constellations); i++ {
		if constellations[i-1].ConstellationID >= constellations[i].ConstellationID {
			t.Errorf("Constellations not sorted: %d >= %d", constellations[i-1].ConstellationID, constellations[i].ConstellationID)
		}
	}

	// Verify constellation has correct region ID reference
	for _, c := range constellations {
		if c.ConstellationID == 20000020 && c.RegionID != 10000002 {
			t.Errorf("Constellation 20000020 should have RegionID 10000002, got %d", c.RegionID)
		}
	}
}

func TestParser_ParseSolarSystems(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	// Parse stars first to get the type map
	starTypeMap, err := p.ParseStars()
	if err != nil {
		t.Fatalf("ParseStars failed: %v", err)
	}

	systems, err := p.ParseSolarSystems(starTypeMap)
	if err != nil {
		t.Fatalf("ParseSolarSystems failed: %v", err)
	}

	if len(systems) != 3 {
		t.Errorf("Expected 3 solar systems, got %d", len(systems))
	}

	// Verify systems are sorted by ID
	for i := 1; i < len(systems); i++ {
		if systems[i-1].SolarSystemID >= systems[i].SolarSystemID {
			t.Errorf("Solar systems not sorted: %d >= %d", systems[i-1].SolarSystemID, systems[i].SolarSystemID)
		}
	}

	// Find Jita and verify its data
	var jita *struct {
		id       int64
		security float64
		sunType  int64
	}
	for _, s := range systems {
		if s.SolarSystemName == "Jita" {
			var sunType int64
			if s.SunTypeID != nil {
				sunType = *s.SunTypeID
			}
			jita = &struct {
				id       int64
				security float64
				sunType  int64
			}{s.SolarSystemID, s.Security, sunType}
			break
		}
	}

	if jita == nil {
		t.Fatal("Jita not found in parsed systems")
	}

	if jita.id != 30000142 {
		t.Errorf("Expected Jita ID to be 30000142, got %d", jita.id)
	}

	// Note: raw security value is 0.9459, transformation happens in transformer
	if jita.security < 0.9 || jita.security > 1.0 {
		t.Errorf("Jita security value unexpected: %f", jita.security)
	}

	// Verify sunTypeID is now the type ID (3796) not the star ID (40000006)
	if jita.sunType != 3796 {
		t.Errorf("Expected Jita sunTypeID to be 3796 (type ID), got %d", jita.sunType)
	}
}

func TestParser_ParseStars(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	starTypeMap, err := p.ParseStars()
	if err != nil {
		t.Fatalf("ParseStars failed: %v", err)
	}

	// Verify correct number of stars parsed
	if len(starTypeMap) != 3 {
		t.Errorf("Expected 3 stars, got %d", len(starTypeMap))
	}

	// Verify Jita's star (40000006) maps to type 3796
	if typeID, ok := starTypeMap[40000006]; !ok {
		t.Error("Star 40000006 not found in map")
	} else if typeID != 3796 {
		t.Errorf("Expected star 40000006 to have typeID 3796, got %d", typeID)
	}

	// Verify Tanoo's star (40000007) maps to type 3797
	if typeID, ok := starTypeMap[40000007]; !ok {
		t.Error("Star 40000007 not found in map")
	} else if typeID != 3797 {
		t.Errorf("Expected star 40000007 to have typeID 3797, got %d", typeID)
	}

	// Verify wormhole star (40045041) maps to type 45041
	if typeID, ok := starTypeMap[40045041]; !ok {
		t.Error("Star 40045041 not found in map")
	} else if typeID != 45041 {
		t.Errorf("Expected star 40045041 to have typeID 45041, got %d", typeID)
	}
}

func TestParser_ParseSolarSystemsWithNilStarMap(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	// Parse with nil star map - should still work but sunTypeID will be nil
	systems, err := p.ParseSolarSystems(nil)
	if err != nil {
		t.Fatalf("ParseSolarSystems failed: %v", err)
	}

	if len(systems) != 3 {
		t.Errorf("Expected 3 solar systems, got %d", len(systems))
	}

	// Find Jita and verify sunTypeID is nil when no star map provided
	for _, s := range systems {
		if s.SolarSystemName == "Jita" {
			if s.SunTypeID != nil {
				t.Errorf("Expected Jita sunTypeID to be nil without star map, got %v", *s.SunTypeID)
			}
			return
		}
	}
	t.Fatal("Jita not found in parsed systems")
}

func TestParser_ParseStargates(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	jumps, err := p.ParseStargates()
	if err != nil {
		t.Fatalf("ParseStargates failed: %v", err)
	}

	// Two stargates in both directions should result in 2 jumps (A→B and B→A)
	// This matches Fuzzwork CSV format where both directions are included
	if len(jumps) != 2 {
		t.Errorf("Expected 2 jumps (bidirectional), got %d", len(jumps))
	}

	// Verify both directions exist
	foundJitaToTanoo := false
	foundTanooToJita := false
	for _, jump := range jumps {
		if jump.FromSolarSystemID == 30000142 && jump.ToSolarSystemID == 30000001 {
			foundJitaToTanoo = true
		}
		if jump.FromSolarSystemID == 30000001 && jump.ToSolarSystemID == 30000142 {
			foundTanooToJita = true
		}
	}

	if !foundJitaToTanoo {
		t.Error("Expected jump from Jita (30000142) to Tanoo (30000001)")
	}
	if !foundTanooToJita {
		t.Error("Expected jump from Tanoo (30000001) to Jita (30000142)")
	}
}

func TestParser_ParseTypes(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	types, err := p.ParseTypes()
	if err != nil {
		t.Fatalf("ParseTypes failed: %v", err)
	}

	if len(types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(types))
	}

	// Verify Rifter data
	rifter, ok := types[587]
	if !ok {
		t.Fatal("Rifter (587) not found in types")
	}

	if rifter.GroupID != 25 {
		t.Errorf("Expected Rifter groupID to be 25, got %d", rifter.GroupID)
	}

	if rifter.Name["en"] != "Rifter" {
		t.Errorf("Expected Rifter name to be 'Rifter', got %q", rifter.Name["en"])
	}

	if rifter.Mass != 1350000.0 {
		t.Errorf("Expected Rifter mass to be 1350000.0, got %f", rifter.Mass)
	}
}

func TestParser_ParseGroups(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	groups, err := p.ParseGroups()
	if err != nil {
		t.Fatalf("ParseGroups failed: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	// Verify Frigate group
	frigate, ok := groups[25]
	if !ok {
		t.Fatal("Frigate group (25) not found")
	}

	if frigate.CategoryID != 6 {
		t.Errorf("Expected Frigate categoryID to be 6, got %d", frigate.CategoryID)
	}

	if frigate.Name["en"] != "Frigate" {
		t.Errorf("Expected Frigate name, got %q", frigate.Name["en"])
	}
}

func TestParser_ParseCategories(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	categories, err := p.ParseCategories()
	if err != nil {
		t.Fatalf("ParseCategories failed: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Verify Ship category
	ship, ok := categories[6]
	if !ok {
		t.Fatal("Ship category (6) not found")
	}

	if ship.Name["en"] != "Ship" {
		t.Errorf("Expected Ship name, got %q", ship.Name["en"])
	}
}

func TestParser_ParseAll(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	result, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	// Verify all data was parsed
	if len(result.Regions) == 0 {
		t.Error("No regions parsed")
	}
	if len(result.Constellations) == 0 {
		t.Error("No constellations parsed")
	}
	if len(result.SolarSystems) == 0 {
		t.Error("No solar systems parsed")
	}
	if len(result.Types) == 0 {
		t.Error("No types parsed")
	}
	if len(result.Groups) == 0 {
		t.Error("No groups parsed")
	}
	if len(result.Categories) == 0 {
		t.Error("No categories parsed")
	}
	if len(result.SystemJumps) == 0 {
		t.Error("No system jumps parsed")
	}
}

func TestParser_MissingFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parser_test_empty")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	// Try to parse non-existent file
	_, err = p.ParseRegions()
	if err == nil {
		t.Error("Expected error when parsing missing file")
	}
}

func TestParser_MalformedYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parser_test_malformed")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create malformed YAML
	malformedYAML := `this is not valid yaml: [[[`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(malformedYAML), 0644); err != nil {
		t.Fatalf("failed to create malformed yaml: %v", err)
	}

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	_, err = p.ParseRegions()
	if err == nil {
		t.Error("Expected error when parsing malformed YAML")
	}
}

func TestParser_EmptyName(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parser_test_empty_name")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create region with empty name
	regionsYAML := `10000001:
  regionID: 10000001
  name: {}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(regionsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapRegions.yaml: %v", err)
	}

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	regions, err := p.ParseRegions()
	if err != nil {
		t.Fatalf("ParseRegions failed: %v", err)
	}

	// Should have a fallback name
	if len(regions) == 0 {
		t.Fatal("No regions parsed")
	}

	if regions[0].RegionName == "" {
		t.Error("Expected fallback name for region with empty name")
	}
}

// TestParser_RegionCoordinates verifies that region coordinates are parsed correctly for CSV output.
func TestParser_RegionCoordinates(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	regions, err := p.ParseRegions()
	if err != nil {
		t.Fatalf("ParseRegions failed: %v", err)
	}

	// Find The Forge region
	var theForge *struct {
		x, y, z          float64
		xMin, yMin, zMin float64
		xMax, yMax, zMax float64
		factionID        *int64
	}
	for _, r := range regions {
		if r.RegionName == "The Forge" {
			theForge = &struct {
				x, y, z          float64
				xMin, yMin, zMin float64
				xMax, yMax, zMax float64
				factionID        *int64
			}{
				r.X, r.Y, r.Z,
				r.XMin, r.YMin, r.ZMin,
				r.XMax, r.YMax, r.ZMax,
				r.FactionID,
			}
			break
		}
	}

	if theForge == nil {
		t.Fatal("The Forge region not found")
	}

	// Verify position coordinates (new SDE format uses position object)
	if theForge.x != -96538397329247680 {
		t.Errorf("Expected X coordinate -96538397329247680, got %f", theForge.x)
	}
	if theForge.y != 68904722523889856 {
		t.Errorf("Expected Y coordinate 68904722523889856, got %f", theForge.y)
	}
	if theForge.z != 103886273221498080 {
		t.Errorf("Expected Z coordinate 103886273221498080, got %f", theForge.z)
	}

	// Note: Min/Max coordinates are not available in new SDE format (2025+)

	// Verify faction ID
	if theForge.factionID == nil || *theForge.factionID != 500001 {
		t.Error("Expected factionID 500001 for The Forge")
	}
}

// TestParser_ConstellationCoordinates verifies that constellation coordinates are parsed correctly.
func TestParser_ConstellationCoordinates(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	constellations, err := p.ParseConstellations()
	if err != nil {
		t.Fatalf("ParseConstellations failed: %v", err)
	}

	// Find Kimotoro constellation
	var kimotoro *struct {
		x, y, z          float64
		xMin, yMin, zMin float64
		xMax, yMax, zMax float64
		factionID        *int64
		radius           float64
	}
	for _, c := range constellations {
		if c.ConstellationName == "Kimotoro" {
			kimotoro = &struct {
				x, y, z          float64
				xMin, yMin, zMin float64
				xMax, yMax, zMax float64
				factionID        *int64
				radius           float64
			}{
				c.X, c.Y, c.Z,
				c.XMin, c.YMin, c.ZMin,
				c.XMax, c.YMax, c.ZMax,
				c.FactionID,
				c.Radius,
			}
			break
		}
	}

	if kimotoro == nil {
		t.Fatal("Kimotoro constellation not found")
	}

	// Verify position coordinates (new SDE format)
	if kimotoro.x != -107314934797574880 {
		t.Errorf("Expected X coordinate -107314934797574880, got %f", kimotoro.x)
	}

	// Verify radius
	if kimotoro.radius != 12345678901234.5 {
		t.Errorf("Expected radius 12345678901234.5, got %f", kimotoro.radius)
	}

	// Verify faction ID
	if kimotoro.factionID == nil || *kimotoro.factionID != 500001 {
		t.Error("Expected factionID 500001 for Kimotoro")
	}
}

// TestParser_SolarSystemCSVFields verifies all CSV-required fields are parsed for solar systems.
func TestParser_SolarSystemCSVFields(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	// Parse stars first
	starTypeMap, err := p.ParseStars()
	if err != nil {
		t.Fatalf("ParseStars failed: %v", err)
	}

	systems, err := p.ParseSolarSystems(starTypeMap)
	if err != nil {
		t.Fatalf("ParseSolarSystems failed: %v", err)
	}

	// Find Jita
	var jita *struct {
		solarSystemID   int64
		regionID        int64
		constellationID int64
		solarSystemName string
		x, y, z         float64
		xMin, xMax      float64
		yMin, yMax      float64
		zMin, zMax      float64
		luminosity      float64
		border          bool
		fringe          bool
		corridor        bool
		hub             bool
		international   bool
		regional        bool
		constellation   string
		security        float64
		factionID       *int64
		radius          float64
		sunTypeID       *int64
		securityClass   string
	}

	for _, s := range systems {
		if s.SolarSystemName == "Jita" {
			jita = &struct {
				solarSystemID   int64
				regionID        int64
				constellationID int64
				solarSystemName string
				x, y, z         float64
				xMin, xMax      float64
				yMin, yMax      float64
				zMin, zMax      float64
				luminosity      float64
				border          bool
				fringe          bool
				corridor        bool
				hub             bool
				international   bool
				regional        bool
				constellation   string
				security        float64
				factionID       *int64
				radius          float64
				sunTypeID       *int64
				securityClass   string
			}{
				s.SolarSystemID,
				s.RegionID,
				s.ConstellationID,
				s.SolarSystemName,
				s.X, s.Y, s.Z,
				s.XMin, s.XMax,
				s.YMin, s.YMax,
				s.ZMin, s.ZMax,
				s.Luminosity,
				s.Border,
				s.Fringe,
				s.Corridor,
				s.Hub,
				s.International,
				s.Regional,
				s.Constellation,
				s.Security,
				s.FactionID,
				s.Radius,
				s.SunTypeID,
				s.SecurityClass,
			}
			break
		}
	}

	if jita == nil {
		t.Fatal("Jita not found")
	}

	// Verify IDs
	if jita.solarSystemID != 30000142 {
		t.Errorf("Expected solarSystemID 30000142, got %d", jita.solarSystemID)
	}
	if jita.regionID != 10000002 {
		t.Errorf("Expected regionID 10000002, got %d", jita.regionID)
	}
	if jita.constellationID != 20000020 {
		t.Errorf("Expected constellationID 20000020, got %d", jita.constellationID)
	}

	// Verify position coordinates (new SDE format)
	if jita.x != -129500988494612512 {
		t.Errorf("Expected X -129500988494612512, got %f", jita.x)
	}
	// Note: Min/Max coordinates are not available in new SDE format (2025+)

	// Verify boolean flags
	if !jita.border {
		t.Error("Expected border to be true")
	}
	if jita.fringe {
		t.Error("Expected fringe to be false")
	}
	if !jita.corridor {
		t.Error("Expected corridor to be true")
	}
	if !jita.hub {
		t.Error("Expected hub to be true")
	}
	if jita.international {
		t.Error("Expected international to be false")
	}
	if !jita.regional {
		t.Error("Expected regional to be true")
	}

	// Verify constellation field (always "None")
	if jita.constellation != "None" {
		t.Errorf("Expected constellation 'None', got %q", jita.constellation)
	}

	// Verify other fields
	if jita.luminosity != 0.01575 {
		t.Errorf("Expected luminosity 0.01575, got %f", jita.luminosity)
	}
	if jita.radius != 3.35e+12 {
		t.Errorf("Expected radius 3.35e+12, got %f", jita.radius)
	}
	if jita.securityClass != "B" {
		t.Errorf("Expected securityClass 'B', got %q", jita.securityClass)
	}
	// Note: factionID is not present on systems in new SDE format - propagated from region
	// sunTypeID is now resolved to the actual type ID via star parsing
	if jita.sunTypeID == nil || *jita.sunTypeID != 3796 {
		t.Errorf("Expected sunTypeID 3796 (resolved from starID 40000006), got %v", jita.sunTypeID)
	}
}

// TestParser_TypeCSVFields verifies all CSV-required fields are parsed for types.
func TestParser_TypeCSVFields(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	types, err := p.ParseTypes()
	if err != nil {
		t.Fatalf("ParseTypes failed: %v", err)
	}

	// Check Rifter has all required fields
	rifter, ok := types[587]
	if !ok {
		t.Fatal("Rifter (587) not found")
	}

	if rifter.GroupID != 25 {
		t.Errorf("Expected groupID 25, got %d", rifter.GroupID)
	}
	if rifter.Name["en"] != "Rifter" {
		t.Errorf("Expected name 'Rifter', got %q", rifter.Name["en"])
	}
	if rifter.Description["en"] != "The Rifter is a very powerful combat frigate." {
		t.Errorf("Expected description, got %q", rifter.Description["en"])
	}
	if rifter.Mass != 1350000.0 {
		t.Errorf("Expected mass 1350000.0, got %f", rifter.Mass)
	}
	if rifter.Volume != 27500.0 {
		t.Errorf("Expected volume 27500.0, got %f", rifter.Volume)
	}
	if rifter.Capacity != 125.0 {
		t.Errorf("Expected capacity 125.0, got %f", rifter.Capacity)
	}
	if rifter.PortionSize != 1 {
		t.Errorf("Expected portionSize 1, got %d", rifter.PortionSize)
	}
	if rifter.RaceID != 2 {
		t.Errorf("Expected raceID 2, got %d", rifter.RaceID)
	}
	if rifter.BasePrice != 240000.0 {
		t.Errorf("Expected basePrice 240000.0, got %f", rifter.BasePrice)
	}
	if !rifter.Published {
		t.Error("Expected published to be true")
	}
	if rifter.MarketGroupID != 64 {
		t.Errorf("Expected marketGroupID 64, got %d", rifter.MarketGroupID)
	}
	if rifter.IconID != 588 {
		t.Errorf("Expected iconID 588, got %d", rifter.IconID)
	}
	if rifter.SoundID != 20071 {
		t.Errorf("Expected soundID 20071, got %d", rifter.SoundID)
	}
	if rifter.GraphicID != 588 {
		t.Errorf("Expected graphicID 588, got %d", rifter.GraphicID)
	}
}

// TestParser_GroupCSVFields verifies all CSV-required fields are parsed for groups.
func TestParser_GroupCSVFields(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	groups, err := p.ParseGroups()
	if err != nil {
		t.Fatalf("ParseGroups failed: %v", err)
	}

	// Check Frigate group has all required fields
	frigate, ok := groups[25]
	if !ok {
		t.Fatal("Frigate group (25) not found")
	}

	if frigate.CategoryID != 6 {
		t.Errorf("Expected categoryID 6, got %d", frigate.CategoryID)
	}
	if frigate.Name["en"] != "Frigate" {
		t.Errorf("Expected name 'Frigate', got %q", frigate.Name["en"])
	}
	if frigate.IconID != 73 {
		t.Errorf("Expected iconID 73, got %d", frigate.IconID)
	}
	if !frigate.UseBasePrice {
		t.Error("Expected useBasePrice to be true")
	}
	if frigate.Anchored {
		t.Error("Expected anchored to be false")
	}
	if frigate.Anchorable {
		t.Error("Expected anchorable to be false")
	}
	if frigate.FittableNonSingleton {
		t.Error("Expected fittableNonSingleton to be false")
	}
	if !frigate.Published {
		t.Error("Expected published to be true")
	}

	// Check Drone group for fittableNonSingleton = true
	drone, ok := groups[18]
	if !ok {
		t.Fatal("Drone group (18) not found")
	}
	if !drone.FittableNonSingleton {
		t.Error("Expected drone fittableNonSingleton to be true")
	}
}
