package transformer

import (
	"testing"

	"github.com/guarzo/wanderer-sde/internal/models"
)

func TestCalculateRegionBounds(t *testing.T) {
	regions := []models.Region{
		{RegionID: 1, RegionName: "Region A"},
		{RegionID: 2, RegionName: "Region B"},
		{RegionID: 3, RegionName: "Empty Region"}, // No systems
	}

	systems := []models.SolarSystem{
		// Region 1 systems
		{SolarSystemID: 101, RegionID: 1, X: 10, Y: 20, Z: 30},
		{SolarSystemID: 102, RegionID: 1, X: 50, Y: -10, Z: 100},
		{SolarSystemID: 103, RegionID: 1, X: -5, Y: 40, Z: 15},
		// Region 2 systems
		{SolarSystemID: 201, RegionID: 2, X: 1000, Y: 2000, Z: 3000},
		{SolarSystemID: 202, RegionID: 2, X: 1500, Y: 2500, Z: 3500},
	}

	CalculateRegionBounds(regions, systems)

	// Check Region 1 bounds
	if regions[0].XMin != -5 || regions[0].XMax != 50 {
		t.Errorf("Region 1 X bounds: expected [-5, 50], got [%f, %f]", regions[0].XMin, regions[0].XMax)
	}
	if regions[0].YMin != -10 || regions[0].YMax != 40 {
		t.Errorf("Region 1 Y bounds: expected [-10, 40], got [%f, %f]", regions[0].YMin, regions[0].YMax)
	}
	if regions[0].ZMin != 15 || regions[0].ZMax != 100 {
		t.Errorf("Region 1 Z bounds: expected [15, 100], got [%f, %f]", regions[0].ZMin, regions[0].ZMax)
	}

	// Check Region 2 bounds
	if regions[1].XMin != 1000 || regions[1].XMax != 1500 {
		t.Errorf("Region 2 X bounds: expected [1000, 1500], got [%f, %f]", regions[1].XMin, regions[1].XMax)
	}
	if regions[1].YMin != 2000 || regions[1].YMax != 2500 {
		t.Errorf("Region 2 Y bounds: expected [2000, 2500], got [%f, %f]", regions[1].YMin, regions[1].YMax)
	}
	if regions[1].ZMin != 3000 || regions[1].ZMax != 3500 {
		t.Errorf("Region 2 Z bounds: expected [3000, 3500], got [%f, %f]", regions[1].ZMin, regions[1].ZMax)
	}

	// Check Empty Region (should remain zeros)
	if regions[2].XMin != 0 || regions[2].XMax != 0 {
		t.Errorf("Empty region should have zero bounds, got X [%f, %f]", regions[2].XMin, regions[2].XMax)
	}
}

func TestCalculateRegionBounds_SingleSystem(t *testing.T) {
	regions := []models.Region{
		{RegionID: 1, RegionName: "Single System Region"},
	}

	systems := []models.SolarSystem{
		{SolarSystemID: 101, RegionID: 1, X: 100, Y: 200, Z: 300},
	}

	CalculateRegionBounds(regions, systems)

	// With a single system, min and max should be the same
	if regions[0].XMin != 100 || regions[0].XMax != 100 {
		t.Errorf("Single system region X bounds should be [100, 100], got [%f, %f]", regions[0].XMin, regions[0].XMax)
	}
	if regions[0].YMin != 200 || regions[0].YMax != 200 {
		t.Errorf("Single system region Y bounds should be [200, 200], got [%f, %f]", regions[0].YMin, regions[0].YMax)
	}
	if regions[0].ZMin != 300 || regions[0].ZMax != 300 {
		t.Errorf("Single system region Z bounds should be [300, 300], got [%f, %f]", regions[0].ZMin, regions[0].ZMax)
	}
}

func TestCalculateConstellationBounds(t *testing.T) {
	constellations := []models.Constellation{
		{ConstellationID: 1, ConstellationName: "Constellation A"},
		{ConstellationID: 2, ConstellationName: "Constellation B"},
		{ConstellationID: 3, ConstellationName: "Empty Constellation"},
	}

	systems := []models.SolarSystem{
		// Constellation 1 systems
		{SolarSystemID: 101, ConstellationID: 1, X: -100, Y: -200, Z: -300},
		{SolarSystemID: 102, ConstellationID: 1, X: 100, Y: 200, Z: 300},
		// Constellation 2 systems
		{SolarSystemID: 201, ConstellationID: 2, X: 500, Y: 500, Z: 500},
	}

	CalculateConstellationBounds(constellations, systems)

	// Check Constellation 1 bounds
	if constellations[0].XMin != -100 || constellations[0].XMax != 100 {
		t.Errorf("Constellation 1 X bounds: expected [-100, 100], got [%f, %f]", constellations[0].XMin, constellations[0].XMax)
	}
	if constellations[0].YMin != -200 || constellations[0].YMax != 200 {
		t.Errorf("Constellation 1 Y bounds: expected [-200, 200], got [%f, %f]", constellations[0].YMin, constellations[0].YMax)
	}
	if constellations[0].ZMin != -300 || constellations[0].ZMax != 300 {
		t.Errorf("Constellation 1 Z bounds: expected [-300, 300], got [%f, %f]", constellations[0].ZMin, constellations[0].ZMax)
	}

	// Check Constellation 2 bounds (single system)
	if constellations[1].XMin != 500 || constellations[1].XMax != 500 {
		t.Errorf("Constellation 2 X bounds: expected [500, 500], got [%f, %f]", constellations[1].XMin, constellations[1].XMax)
	}
	if constellations[1].YMin != 500 || constellations[1].YMax != 500 {
		t.Errorf("Constellation 2 Y bounds: expected [500, 500], got [%f, %f]", constellations[1].YMin, constellations[1].YMax)
	}
	if constellations[1].ZMin != 500 || constellations[1].ZMax != 500 {
		t.Errorf("Constellation 2 Z bounds: expected [500, 500], got [%f, %f]", constellations[1].ZMin, constellations[1].ZMax)
	}

	// Check Empty Constellation (should remain zeros)
	if constellations[2].XMin != 0 || constellations[2].XMax != 0 {
		t.Errorf("Empty constellation should have zero X bounds, expected [0, 0], got [%f, %f]", constellations[2].XMin, constellations[2].XMax)
	}
	if constellations[2].YMin != 0 || constellations[2].YMax != 0 {
		t.Errorf("Empty constellation should have zero Y bounds, expected [0, 0], got [%f, %f]", constellations[2].YMin, constellations[2].YMax)
	}
	if constellations[2].ZMin != 0 || constellations[2].ZMax != 0 {
		t.Errorf("Empty constellation should have zero Z bounds, expected [0, 0], got [%f, %f]", constellations[2].ZMin, constellations[2].ZMax)
	}
}

func TestInheritFactionIDs(t *testing.T) {
	faction1 := int64(500001)
	faction2 := int64(500002)
	systemFaction := int64(500003)

	regions := []models.Region{
		{RegionID: 1, RegionName: "Caldari Space", FactionID: &faction1},
		{RegionID: 2, RegionName: "Gallente Space", FactionID: &faction2},
		{RegionID: 3, RegionName: "No Faction Region", FactionID: nil},
	}

	systems := []models.SolarSystem{
		// System without faction, should inherit from region 1
		{SolarSystemID: 101, RegionID: 1, FactionID: nil},
		// System with its own faction, should NOT inherit
		{SolarSystemID: 102, RegionID: 1, FactionID: &systemFaction},
		// System in region 2 without faction, should inherit
		{SolarSystemID: 201, RegionID: 2, FactionID: nil},
		// System in region without faction, should remain nil
		{SolarSystemID: 301, RegionID: 3, FactionID: nil},
	}

	InheritFactionIDs(systems, regions)

	// Check system 101 inherited faction 1
	if systems[0].FactionID == nil || *systems[0].FactionID != faction1 {
		t.Errorf("System 101 should inherit factionID %d from region", faction1)
	}

	// Check system 102 kept its own faction
	if systems[1].FactionID == nil || *systems[1].FactionID != systemFaction {
		t.Errorf("System 102 should keep its own factionID %d", systemFaction)
	}

	// Check system 201 inherited faction 2
	if systems[2].FactionID == nil || *systems[2].FactionID != faction2 {
		t.Errorf("System 201 should inherit factionID %d from region", faction2)
	}

	// Check system 301 remains nil (region has no faction)
	if systems[3].FactionID != nil {
		t.Errorf("System 301 should have nil factionID since region has none")
	}
}

func TestInheritFactionIDs_NoPointerAliasing(t *testing.T) {
	faction1 := int64(500001)

	regions := []models.Region{
		{RegionID: 1, RegionName: "Region", FactionID: &faction1},
	}

	systems := []models.SolarSystem{
		{SolarSystemID: 101, RegionID: 1, FactionID: nil},
		{SolarSystemID: 102, RegionID: 1, FactionID: nil},
	}

	InheritFactionIDs(systems, regions)

	// Verify both systems have the correct value
	if systems[0].FactionID == nil || *systems[0].FactionID != faction1 {
		t.Errorf("System 101 should have factionID %d", faction1)
	}
	if systems[1].FactionID == nil || *systems[1].FactionID != faction1 {
		t.Errorf("System 102 should have factionID %d", faction1)
	}

	// Verify they don't share the same pointer (modifying one shouldn't affect the other)
	*systems[0].FactionID = 999999
	if *systems[1].FactionID == 999999 {
		t.Error("Systems should not share the same pointer - modifying one affected the other")
	}
}
