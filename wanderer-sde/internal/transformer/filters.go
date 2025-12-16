package transformer

import (
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
)

// ShipCategoryID is the category ID for ships in EVE Online.
const ShipCategoryID = 6

// FilterShipTypes filters the type data to only include ships.
// Ships are identified by having a groupID that belongs to a group
// with categoryID == ShipCategoryID (6).
func FilterShipTypes(types map[int64]models.SDEType, groups map[int64]models.SDEGroup) []models.ShipType {
	// Build set of ship group IDs (groups that belong to the ship category)
	shipGroupIDs := make(map[int64]bool)
	for groupID, group := range groups {
		if group.CategoryID == ShipCategoryID {
			shipGroupIDs[groupID] = true
		}
	}

	var ships []models.ShipType
	for typeID, typeData := range types {
		if shipGroupIDs[typeData.GroupID] {
			typeName := typeData.Name["en"]

			ships = append(ships, models.ShipType{
				TypeID:   typeID,
				GroupID:  typeData.GroupID,
				TypeName: typeName,
				Mass:     typeData.Mass,
				Volume:   typeData.Volume,
				Capacity: typeData.Capacity,
			})
		}
	}

	// Sort by typeID for consistent output
	sort.Slice(ships, func(i, j int) bool {
		return ships[i].TypeID < ships[j].TypeID
	})

	return ships
}

// FilterShipGroups filters the group data to only include ship groups.
func FilterShipGroups(groups map[int64]models.SDEGroup) []models.ItemGroup {
	var shipGroups []models.ItemGroup

	for groupID, group := range groups {
		if group.CategoryID == ShipCategoryID {
			groupName := group.Name["en"]

			shipGroups = append(shipGroups, models.ItemGroup{
				GroupID:    groupID,
				CategoryID: group.CategoryID,
				GroupName:  groupName,
			})
		}
	}

	// Sort by groupID for consistent output
	sort.Slice(shipGroups, func(i, j int) bool {
		return shipGroups[i].GroupID < shipGroups[j].GroupID
	})

	return shipGroups
}

// FilterPublishedTypes filters to only include published types.
func FilterPublishedTypes(types map[int64]models.SDEType) map[int64]models.SDEType {
	published := make(map[int64]models.SDEType)
	for typeID, typeData := range types {
		if typeData.Published {
			published[typeID] = typeData
		}
	}
	return published
}

// FilterPublishedGroups filters to only include published groups.
func FilterPublishedGroups(groups map[int64]models.SDEGroup) map[int64]models.SDEGroup {
	published := make(map[int64]models.SDEGroup)
	for groupID, groupData := range groups {
		if groupData.Published {
			published[groupID] = groupData
		}
	}
	return published
}
