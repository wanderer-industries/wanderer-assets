package transformer

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/internal/parser"
)

// Transformer handles transformation of parsed SDE data to Wanderer format.
type Transformer struct {
	config *config.Config
}

// New creates a new Transformer with the given configuration.
func New(cfg *config.Config) *Transformer {
	return &Transformer{
		config: cfg,
	}
}

// Transform converts parsed SDE data into Wanderer's output format.
func (t *Transformer) Transform(parseResult *parser.ParseResult) (*models.ConvertedData, error) {
	if t.config.Verbose {
		fmt.Println("Transforming SDE data...")
	}

	// Transform solar systems with security calculation
	if t.config.Verbose {
		fmt.Println("  Transforming solar systems...")
	}
	systems := t.transformSolarSystems(parseResult.SolarSystems)

	// Sort regions for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting regions...")
	}
	regions := t.sortRegions(parseResult.Regions)

	// Sort constellations for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting constellations...")
	}
	constellations := t.sortConstellations(parseResult.Constellations)

	// Transform groups (filter to ships only - category 6)
	if t.config.Verbose {
		fmt.Println("  Transforming groups (ships only)...")
	}
	invGroups := t.transformShipGroups(parseResult.Groups)

	// Transform types (filter to ships only - groups with category 6)
	if t.config.Verbose {
		fmt.Println("  Transforming types (ships only)...")
	}
	invTypes := t.transformShipTypes(parseResult.Types, parseResult.Groups)

	// Sort wormhole classes for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting wormhole classes...")
	}
	wormholeClasses := t.sortWormholeClasses(parseResult.WormholeClasses)

	// Transform system jumps with region/constellation lookup
	if t.config.Verbose {
		fmt.Println("  Transforming system jumps...")
	}
	systemJumps := t.transformSystemJumps(parseResult.SystemJumps, systems)

	// Calculate bounds for regions and constellations from constituent systems
	if t.config.Verbose {
		fmt.Println("  Calculating region bounds...")
	}
	CalculateRegionBounds(regions, systems)

	if t.config.Verbose {
		fmt.Println("  Calculating constellation bounds...")
	}
	CalculateConstellationBounds(constellations, systems)

	// Inherit factionID from region for systems that don't have one
	if t.config.Verbose {
		fmt.Println("  Inheriting faction IDs...")
	}
	InheritFactionIDs(systems, regions)

	// Transform NPC stations
	if t.config.Verbose {
		fmt.Println("  Transforming NPC stations...")
	}
	npcStations := t.transformNPCStations(parseResult.NPCStations, parseResult.NPCCorporations)

	result := &models.ConvertedData{
		Universe: &models.UniverseData{
			Regions:        regions,
			Constellations: constellations,
			SolarSystems:   systems,
		},
		InvTypes:        invTypes,
		InvGroups:       invGroups,
		WormholeClasses: wormholeClasses,
		SystemJumps:     systemJumps,
		NPCStations:     npcStations,
	}

	if t.config.Verbose {
		fmt.Printf("Transformation complete:\n")
		fmt.Printf("  Regions:         %d\n", len(result.Universe.Regions))
		fmt.Printf("  Constellations:  %d\n", len(result.Universe.Constellations))
		fmt.Printf("  Solar Systems:   %d\n", len(result.Universe.SolarSystems))
		fmt.Printf("  Types:           %d\n", len(result.InvTypes))
		fmt.Printf("  Groups:          %d\n", len(result.InvGroups))
		fmt.Printf("  Wormhole Classes: %d\n", len(result.WormholeClasses))
		fmt.Printf("  System Jumps:    %d\n", len(result.SystemJumps))
		fmt.Printf("  NPC Stations:    %d\n", len(result.NPCStations))
	}

	return result, nil
}

// transformSolarSystems sorts solar systems while preserving all fields.
// Note: We output raw security values (not rounded) because Wanderer calculates
// true security itself from the raw value.
func (t *Transformer) transformSolarSystems(systems []models.SolarSystem) []models.SolarSystem {
	result := make([]models.SolarSystem, len(systems))
	copy(result, systems)

	// Sort by system ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].SolarSystemID < result[j].SolarSystemID
	})

	return result
}

// sortRegions returns regions sorted by ID.
func (t *Transformer) sortRegions(regions []models.Region) []models.Region {
	result := make([]models.Region, len(regions))
	copy(result, regions)

	sort.Slice(result, func(i, j int) bool {
		return result[i].RegionID < result[j].RegionID
	})

	return result
}

// sortConstellations returns constellations sorted by ID.
func (t *Transformer) sortConstellations(constellations []models.Constellation) []models.Constellation {
	result := make([]models.Constellation, len(constellations))
	copy(result, constellations)

	sort.Slice(result, func(i, j int) bool {
		return result[i].ConstellationID < result[j].ConstellationID
	})

	return result
}

// getLocalizedName retrieves a localized name from a map, logging a warning in verbose mode if missing.
func (t *Transformer) getLocalizedName(names map[string]string, lang string, context string) string {
	if name, ok := names[lang]; ok {
		return name
	}
	if t.config.Verbose {
		fmt.Printf("  Warning: missing '%s' name for %s\n", lang, context)
	}
	return ""
}

// sortWormholeClasses returns wormhole classes sorted by location ID.
func (t *Transformer) sortWormholeClasses(classes []models.WormholeClassLocation) []models.WormholeClassLocation {
	result := make([]models.WormholeClassLocation, len(classes))
	copy(result, classes)

	sort.Slice(result, func(i, j int) bool {
		return result[i].LocationID < result[j].LocationID
	})

	return result
}

// transformShipTypes converts SDE types to InvType format, filtering to ships only.
// Ships are types whose groupID belongs to a group with categoryID == 6.
func (t *Transformer) transformShipTypes(types map[int64]models.SDEType, groups map[int64]models.SDEGroup) []models.InvType {
	// Build set of ship group IDs
	shipGroupIDs := make(map[int64]bool)
	for groupID, group := range groups {
		if group.CategoryID == ShipCategoryID {
			shipGroupIDs[groupID] = true
		}
	}

	result := make([]models.InvType, 0)

	for typeID, sdeType := range types {
		// Only include types that belong to ship groups
		if !shipGroupIDs[sdeType.GroupID] {
			continue
		}

		invType := models.InvType{
			TypeID:    typeID,
			GroupID:   sdeType.GroupID,
			TypeName:  t.getLocalizedName(sdeType.Name, "en", fmt.Sprintf("type %d", typeID)),
			Mass:      sdeType.Mass,
			Volume:    sdeType.Volume,
			Capacity:  sdeType.Capacity,
			Published: sdeType.Published,
		}
		result = append(result, invType)
	}

	// Sort by type ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].TypeID < result[j].TypeID
	})

	return result
}

// transformShipGroups converts SDE groups to InvGroup format, filtering to ships only.
// Ships are groups with categoryID == 6.
func (t *Transformer) transformShipGroups(groups map[int64]models.SDEGroup) []models.InvGroup {
	result := make([]models.InvGroup, 0)

	for groupID, sdeGroup := range groups {
		// Only include ship groups (category 6)
		if sdeGroup.CategoryID != ShipCategoryID {
			continue
		}

		invGroup := models.InvGroup{
			GroupID:    groupID,
			CategoryID: sdeGroup.CategoryID,
			GroupName:  t.getLocalizedName(sdeGroup.Name, "en", fmt.Sprintf("group %d", groupID)),
		}
		result = append(result, invGroup)
	}

	// Sort by group ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].GroupID < result[j].GroupID
	})

	return result
}

// transformNPCStations converts SDE NPC stations to Wanderer format.
func (t *Transformer) transformNPCStations(
	stations map[int64]models.SDENPCStation,
	corps map[int64]models.SDENPCCorporation,
) []models.NPCStation {
	result := make([]models.NPCStation, 0, len(stations))

	for stationID, sdeStation := range stations {
		ownerName := ""
		if corp, ok := corps[sdeStation.OwnerID]; ok {
			ownerName = t.getLocalizedName(corp.Name, "en", fmt.Sprintf("corporation %d", sdeStation.OwnerID))
		}

		station := models.NPCStation{
			StationID:     stationID,
			SolarSystemID: sdeStation.SolarSystemID,
			OwnerID:       sdeStation.OwnerID,
			OwnerName:     ownerName,
			TypeID:        sdeStation.TypeID,
		}
		result = append(result, station)
	}

	// Sort by station ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].StationID < result[j].StationID
	})

	return result
}

// transformSystemJumps enriches system jumps with region and constellation IDs.
func (t *Transformer) transformSystemJumps(jumps []models.SystemJump, systems []models.SolarSystem) []models.SystemJump {
	// Build lookup map for system -> region/constellation
	systemLookup := make(map[int64]models.SolarSystem, len(systems))
	for _, sys := range systems {
		systemLookup[sys.SolarSystemID] = sys
	}

	result := make([]models.SystemJump, 0, len(jumps))
	for _, jump := range jumps {
		fromSys, fromOK := systemLookup[jump.FromSolarSystemID]
		toSys, toOK := systemLookup[jump.ToSolarSystemID]

		enrichedJump := models.SystemJump{
			FromSolarSystemID: jump.FromSolarSystemID,
			ToSolarSystemID:   jump.ToSolarSystemID,
		}

		if fromOK {
			enrichedJump.FromRegionID = fromSys.RegionID
			enrichedJump.FromConstellationID = fromSys.ConstellationID
		} else if t.config.Verbose {
			fmt.Printf("  Warning: unknown from system %d in jump\n", jump.FromSolarSystemID)
		}

		if toOK {
			enrichedJump.ToRegionID = toSys.RegionID
			enrichedJump.ToConstellationID = toSys.ConstellationID
		} else if t.config.Verbose {
			fmt.Printf("  Warning: unknown to system %d in jump\n", jump.ToSolarSystemID)
		}

		result = append(result, enrichedJump)
	}

	// Sort by from system ID, then to system ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		if result[i].FromSolarSystemID != result[j].FromSolarSystemID {
			return result[i].FromSolarSystemID < result[j].FromSolarSystemID
		}
		return result[i].ToSolarSystemID < result[j].ToSolarSystemID
	})

	return result
}

// Validate performs validation checks on the converted data.
func (t *Transformer) Validate(data *models.ConvertedData) *models.ValidationResult {
	result := &models.ValidationResult{
		SolarSystems:    len(data.Universe.SolarSystems),
		Regions:         len(data.Universe.Regions),
		Constellations:  len(data.Universe.Constellations),
		InvTypes:        len(data.InvTypes),
		InvGroups:       len(data.InvGroups),
		SystemJumps:     len(data.SystemJumps),
		WormholeClasses: len(data.WormholeClasses),
		NPCStations:     len(data.NPCStations),
	}

	// Validation thresholds based on known EVE universe size
	const (
		minSolarSystems    = 8000
		minRegions         = 100
		minConstellations  = 1000
		minShipTypes       = 500   // Ships only (category 6), expected ~700+
		minShipGroups      = 30    // Ship groups only (category 6), expected ~50+
		minSystemJumps     = 13000 // Bidirectional jumps (A→B and B→A), expected ~13,776
		minWormholeClasses = 750   // Regions + constellations + systems, expected ~803
		minNPCStations     = 40    // DED stations only, expected ~45
	)

	// Check minimum counts
	if result.SolarSystems < minSolarSystems {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Solar system count (%d) is below expected minimum (%d)",
				result.SolarSystems, minSolarSystems))
	}

	if result.Regions < minRegions {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Region count (%d) is below expected minimum (%d)",
				result.Regions, minRegions))
	}

	if result.Constellations < minConstellations {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Constellation count (%d) is below expected minimum (%d)",
				result.Constellations, minConstellations))
	}

	if result.InvTypes < minShipTypes {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Ship type count (%d) is below expected minimum (%d)",
				result.InvTypes, minShipTypes))
	}

	if result.InvGroups < minShipGroups {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Ship group count (%d) is below expected minimum (%d)",
				result.InvGroups, minShipGroups))
	}

	if result.SystemJumps < minSystemJumps {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("System jump count (%d) is below expected minimum (%d)",
				result.SystemJumps, minSystemJumps))
	}

	if result.WormholeClasses < minWormholeClasses {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Wormhole class count (%d) is below expected minimum (%d)",
				result.WormholeClasses, minWormholeClasses))
	}

	if result.NPCStations < minNPCStations {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("NPC station count (%d) is below expected minimum (%d)",
				result.NPCStations, minNPCStations))
	}

	// Check for empty required data
	if result.SolarSystems == 0 {
		result.Errors = append(result.Errors, "No solar systems found")
	}

	if result.Regions == 0 {
		result.Errors = append(result.Errors, "No regions found")
	}

	if result.Constellations == 0 {
		result.Errors = append(result.Errors, "No constellations found")
	}

	return result
}
