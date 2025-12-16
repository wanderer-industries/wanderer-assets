package parser

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// ExtractAllWormholeClasses extracts wormhole class information from
// regions, constellations, and solar systems.
// This matches Fuzzwork's mapLocationWormholeClasses.csv which includes all three.
func (p *Parser) ExtractAllWormholeClasses() ([]models.WormholeClassLocation, error) {
	var wormholeClasses []models.WormholeClassLocation

	// 1. Extract from regions
	regionsPath := p.filePath("mapRegions.yaml")
	rawRegions, err := yaml.ParseFileMap[int64, SDEMapRegion](regionsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regions for wormhole classes: %w", err)
	}
	for regionID, data := range rawRegions {
		if data.WormholeClassID != 0 {
			wormholeClasses = append(wormholeClasses, models.WormholeClassLocation{
				LocationID:      regionID,
				WormholeClassID: data.WormholeClassID,
			})
		}
	}

	// 2. Extract from constellations
	constellationsPath := p.filePath("mapConstellations.yaml")
	rawConstellations, err := yaml.ParseFileMap[int64, SDEMapConstellation](constellationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse constellations for wormhole classes: %w", err)
	}
	for constellationID, data := range rawConstellations {
		if data.WormholeClassID != 0 {
			wormholeClasses = append(wormholeClasses, models.WormholeClassLocation{
				LocationID:      constellationID,
				WormholeClassID: data.WormholeClassID,
			})
		}
	}

	// 3. Extract from solar systems
	systemsPath := p.filePath("mapSolarSystems.yaml")
	rawSystems, err := yaml.ParseFileMap[int64, SDEMapSolarSystem](systemsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse solar systems for wormhole classes: %w", err)
	}
	for systemID, data := range rawSystems {
		if data.WormholeClassID != 0 {
			wormholeClasses = append(wormholeClasses, models.WormholeClassLocation{
				LocationID:      systemID,
				WormholeClassID: data.WormholeClassID,
			})
		}
	}

	// Sort by location ID for consistent output
	sort.Slice(wormholeClasses, func(i, j int) bool {
		return wormholeClasses[i].LocationID < wormholeClasses[j].LocationID
	})

	return wormholeClasses, nil
}
