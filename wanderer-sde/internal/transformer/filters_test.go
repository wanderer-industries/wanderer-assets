package transformer

import (
	"testing"

	"github.com/guarzo/wanderer-sde/internal/models"
)

func TestFilterShipTypes(t *testing.T) {
	// Create test data
	groups := map[int64]models.SDEGroup{
		25: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Frigate"}},
		26: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Cruiser"}},
		18: {CategoryID: 7, Name: map[string]string{"en": "Drone"}}, // Not a ship
	}

	types := map[int64]models.SDEType{
		587:  {GroupID: 25, Name: map[string]string{"en": "Rifter"}, Mass: 1350000, Volume: 27500},
		588:  {GroupID: 25, Name: map[string]string{"en": "Slasher"}, Mass: 1200000, Volume: 26000},
		625:  {GroupID: 26, Name: map[string]string{"en": "Caracal"}, Mass: 11000000, Volume: 92000},
		2456: {GroupID: 18, Name: map[string]string{"en": "Hobgoblin I"}, Mass: 2500, Volume: 5}, // Drone, not a ship
	}

	ships := FilterShipTypes(types, groups)

	// Should have 3 ships (Rifter, Slasher, Caracal) but not the drone
	if len(ships) != 3 {
		t.Errorf("Expected 3 ships, got %d", len(ships))
	}

	// Verify ships are sorted by typeID
	for i := 1; i < len(ships); i++ {
		if ships[i-1].TypeID >= ships[i].TypeID {
			t.Errorf("Ships not sorted by typeID: %d >= %d", ships[i-1].TypeID, ships[i].TypeID)
		}
	}

	// Verify first ship is Rifter (lowest typeID)
	if len(ships) > 0 && ships[0].TypeID != 587 {
		t.Errorf("Expected first ship to be Rifter (587), got %d", ships[0].TypeID)
	}
}

func TestFilterShipGroups(t *testing.T) {
	groups := map[int64]models.SDEGroup{
		25: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Frigate"}},
		26: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Cruiser"}},
		27: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Battleship"}},
		18: {CategoryID: 7, Name: map[string]string{"en": "Drone"}}, // Not a ship group
	}

	shipGroups := FilterShipGroups(groups)

	// Should have 3 ship groups
	if len(shipGroups) != 3 {
		t.Errorf("Expected 3 ship groups, got %d", len(shipGroups))
	}

	// Verify groups are sorted by groupID
	for i := 1; i < len(shipGroups); i++ {
		if shipGroups[i-1].GroupID >= shipGroups[i].GroupID {
			t.Errorf("Groups not sorted by groupID: %d >= %d", shipGroups[i-1].GroupID, shipGroups[i].GroupID)
		}
	}

	// All groups should have categoryID == ShipCategoryID
	for _, g := range shipGroups {
		if g.CategoryID != ShipCategoryID {
			t.Errorf("Non-ship group in results: %d has categoryID %d", g.GroupID, g.CategoryID)
		}
	}
}

func TestFilterPublishedTypes(t *testing.T) {
	types := map[int64]models.SDEType{
		1: {GroupID: 1, Published: true},
		2: {GroupID: 1, Published: false},
		3: {GroupID: 1, Published: true},
	}

	published := FilterPublishedTypes(types)

	if len(published) != 2 {
		t.Errorf("Expected 2 published types, got %d", len(published))
	}

	if _, ok := published[2]; ok {
		t.Error("Unpublished type 2 should not be in results")
	}
}

func TestFilterPublishedGroups(t *testing.T) {
	groups := map[int64]models.SDEGroup{
		1: {CategoryID: 1, Published: true},
		2: {CategoryID: 1, Published: false},
		3: {CategoryID: 1, Published: true},
	}

	published := FilterPublishedGroups(groups)

	if len(published) != 2 {
		t.Errorf("Expected 2 published groups, got %d", len(published))
	}

	if _, ok := published[2]; ok {
		t.Error("Unpublished group 2 should not be in results")
	}
}

func TestFilterShipTypes_EmptyName(t *testing.T) {
	groups := map[int64]models.SDEGroup{
		25: {CategoryID: ShipCategoryID, Name: map[string]string{"en": "Frigate"}},
	}

	types := map[int64]models.SDEType{
		587: {GroupID: 25, Name: map[string]string{}}, // No English name
	}

	ships := FilterShipTypes(types, groups)

	if len(ships) != 1 {
		t.Errorf("Expected 1 ship, got %d", len(ships))
	}

	// TypeName should be empty string when no English name available
	if ships[0].TypeName != "" {
		t.Errorf("Expected empty TypeName, got %q", ships[0].TypeName)
	}
}
