package transformer

import "github.com/guarzo/wanderer-sde/internal/models"

// bounds holds min/max coordinates for a location.
type bounds struct {
	minX, maxX float64
	minY, maxY float64
	minZ, maxZ float64
}

// calculateBounds computes min/max coordinates from a list of solar systems.
func calculateBounds(systems []models.SolarSystem) bounds {
	if len(systems) == 0 {
		return bounds{}
	}

	b := bounds{
		minX: systems[0].X, maxX: systems[0].X,
		minY: systems[0].Y, maxY: systems[0].Y,
		minZ: systems[0].Z, maxZ: systems[0].Z,
	}

	for _, sys := range systems[1:] {
		if sys.X < b.minX {
			b.minX = sys.X
		}
		if sys.X > b.maxX {
			b.maxX = sys.X
		}
		if sys.Y < b.minY {
			b.minY = sys.Y
		}
		if sys.Y > b.maxY {
			b.maxY = sys.Y
		}
		if sys.Z < b.minZ {
			b.minZ = sys.Z
		}
		if sys.Z > b.maxZ {
			b.maxZ = sys.Z
		}
	}

	return b
}

// CalculateRegionBounds calculates min/max coordinates for regions
// based on their constituent solar systems.
func CalculateRegionBounds(regions []models.Region, systems []models.SolarSystem) {
	// Build map of regionID -> systems
	regionSystems := make(map[int64][]models.SolarSystem)
	for _, sys := range systems {
		regionSystems[sys.RegionID] = append(regionSystems[sys.RegionID], sys)
	}

	// Calculate bounds for each region
	for i := range regions {
		sysList := regionSystems[regions[i].RegionID]
		if len(sysList) == 0 {
			continue
		}

		b := calculateBounds(sysList)
		regions[i].XMin = b.minX
		regions[i].XMax = b.maxX
		regions[i].YMin = b.minY
		regions[i].YMax = b.maxY
		regions[i].ZMin = b.minZ
		regions[i].ZMax = b.maxZ
	}
}

// CalculateConstellationBounds calculates min/max coordinates for constellations
// based on their constituent solar systems.
func CalculateConstellationBounds(constellations []models.Constellation, systems []models.SolarSystem) {
	// Build map of constellationID -> systems
	constellationSystems := make(map[int64][]models.SolarSystem)
	for _, sys := range systems {
		constellationSystems[sys.ConstellationID] = append(constellationSystems[sys.ConstellationID], sys)
	}

	// Calculate bounds for each constellation
	for i := range constellations {
		sysList := constellationSystems[constellations[i].ConstellationID]
		if len(sysList) == 0 {
			continue
		}

		b := calculateBounds(sysList)
		constellations[i].XMin = b.minX
		constellations[i].XMax = b.maxX
		constellations[i].YMin = b.minY
		constellations[i].YMax = b.maxY
		constellations[i].ZMin = b.minZ
		constellations[i].ZMax = b.maxZ
	}
}

// InheritFactionIDs sets factionID on solar systems that don't have one,
// inheriting from their region.
func InheritFactionIDs(systems []models.SolarSystem, regions []models.Region) {
	// Build map of regionID -> factionID
	regionFactions := make(map[int64]*int64)
	for _, region := range regions {
		regionFactions[region.RegionID] = region.FactionID
	}

	// Inherit factionID from region if system doesn't have one
	for i := range systems {
		if systems[i].FactionID == nil {
			if factionID, ok := regionFactions[systems[i].RegionID]; ok && factionID != nil {
				// Copy the value to avoid pointer aliasing
				val := *factionID
				systems[i].FactionID = &val
			}
		}
	}
}
