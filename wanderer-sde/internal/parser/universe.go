package parser

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// SDEMapRegion represents a region in the flat SDE format.
type SDEMapRegion struct {
	RegionID        int64               `yaml:"regionID"`
	Name            map[string]string   `yaml:"name"`
	NameID          int64               `yaml:"nameID,omitempty"`
	DescriptionID   int64               `yaml:"descriptionID,omitempty"`
	FactionID       int64               `yaml:"factionID,omitempty"`
	Position        *models.SDEPosition `yaml:"position,omitempty"`
	NebulaID        int64               `yaml:"nebulaID,omitempty"`
	WormholeClassID int64               `yaml:"wormholeClassID,omitempty"`
}

// SDEMapConstellation represents a constellation in the flat SDE format.
type SDEMapConstellation struct {
	ConstellationID int64               `yaml:"constellationID"`
	RegionID        int64               `yaml:"regionID"`
	Name            map[string]string   `yaml:"name"`
	NameID          int64               `yaml:"nameID,omitempty"`
	FactionID       int64               `yaml:"factionID,omitempty"`
	Position        *models.SDEPosition `yaml:"position,omitempty"`
	Radius          float64             `yaml:"radius,omitempty"`
	WormholeClassID int64               `yaml:"wormholeClassID,omitempty"`
}

// SDEMapSolarSystem represents a solar system in the flat SDE format.
type SDEMapSolarSystem struct {
	SolarSystemID   int64               `yaml:"solarSystemID"`
	ConstellationID int64               `yaml:"constellationID"`
	RegionID        int64               `yaml:"regionID"`
	Name            map[string]string   `yaml:"name"`
	SecurityStatus  float64             `yaml:"securityStatus"`
	SecurityClass   string              `yaml:"securityClass,omitempty"`
	StarID          int64               `yaml:"starID,omitempty"`
	WormholeClassID int64               `yaml:"wormholeClassID,omitempty"`
	FactionID       int64               `yaml:"factionID,omitempty"`
	Border          bool                `yaml:"border,omitempty"`
	Corridor        bool                `yaml:"corridor,omitempty"`
	Fringe          bool                `yaml:"fringe,omitempty"`
	Hub             bool                `yaml:"hub,omitempty"`
	International   bool                `yaml:"international,omitempty"`
	Regional        bool                `yaml:"regional,omitempty"`
	Position        *models.SDEPosition `yaml:"position,omitempty"`
	Luminosity      float64             `yaml:"luminosity,omitempty"`
	Radius          float64             `yaml:"radius,omitempty"`
}

// ParseRegions parses the mapRegions.yaml file.
func (p *Parser) ParseRegions() ([]models.Region, error) {
	path := p.filePath("mapRegions.yaml")

	rawRegions, err := yaml.ParseFileMap[int64, SDEMapRegion](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regions file: %w", err)
	}

	regions := make([]models.Region, 0, len(rawRegions))
	for id, data := range rawRegions {
		name := data.Name["en"]
		if name == "" {
			// Fall back to using ID-based name if no English name
			name = fmt.Sprintf("Region %d", id)
		}

		region := models.Region{
			RegionID:   id,
			RegionName: name,
			FactionID:  models.Int64PtrNonZero(data.FactionID),
			Nebula:     data.NebulaID,
			Radius:     0, // Not directly available in new SDE format
		}

		// Extract coordinates from position object
		if data.Position != nil {
			region.X = data.Position.X
			region.Y = data.Position.Y
			region.Z = data.Position.Z
		}

		regions = append(regions, region)
	}

	// Sort by region ID for consistent output
	sort.Slice(regions, func(i, j int) bool {
		return regions[i].RegionID < regions[j].RegionID
	})

	return regions, nil
}

// ParseConstellations parses the mapConstellations.yaml file.
func (p *Parser) ParseConstellations() ([]models.Constellation, error) {
	path := p.filePath("mapConstellations.yaml")

	rawConstellations, err := yaml.ParseFileMap[int64, SDEMapConstellation](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse constellations file: %w", err)
	}

	constellations := make([]models.Constellation, 0, len(rawConstellations))
	for id, data := range rawConstellations {
		name := data.Name["en"]
		if name == "" {
			name = fmt.Sprintf("Constellation %d", id)
		}

		constellation := models.Constellation{
			RegionID:          data.RegionID,
			ConstellationID:   id,
			ConstellationName: name,
			FactionID:         models.Int64PtrNonZero(data.FactionID),
			Radius:            data.Radius,
		}

		// Extract coordinates from position object
		if data.Position != nil {
			constellation.X = data.Position.X
			constellation.Y = data.Position.Y
			constellation.Z = data.Position.Z
		}

		constellations = append(constellations, constellation)
	}

	// Sort by constellation ID for consistent output
	sort.Slice(constellations, func(i, j int) bool {
		return constellations[i].ConstellationID < constellations[j].ConstellationID
	})

	return constellations, nil
}

// ParseSolarSystems parses the mapSolarSystems.yaml file.
// starTypeMap provides starID -> typeID mapping for resolving sun types.
func (p *Parser) ParseSolarSystems(starTypeMap map[int64]int64) ([]models.SolarSystem, error) {
	path := p.filePath("mapSolarSystems.yaml")

	rawSystems, err := yaml.ParseFileMap[int64, SDEMapSolarSystem](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse solar systems file: %w", err)
	}

	systems := make([]models.SolarSystem, 0, len(rawSystems))
	for id, data := range rawSystems {
		name := data.Name["en"]
		if name == "" {
			name = fmt.Sprintf("System %d", id)
		}

		// Resolve sun type ID from star ID using the provided map
		var sunTypeID *int64
		if data.StarID != 0 && starTypeMap != nil {
			if typeID, ok := starTypeMap[data.StarID]; ok {
				sunTypeID = models.Int64PtrNonZero(typeID)
			}
		}

		system := models.SolarSystem{
			RegionID:        data.RegionID,
			ConstellationID: data.ConstellationID,
			SolarSystemID:   id,
			SolarSystemName: name,
			Luminosity:      data.Luminosity,
			Border:          data.Border,
			Fringe:          data.Fringe,
			Corridor:        data.Corridor,
			Hub:             data.Hub,
			International:   data.International,
			Regional:        data.Regional,
			Constellation:   "None", // Always "None" - legacy field
			Security:        data.SecurityStatus,
			FactionID:       models.Int64PtrNonZero(data.FactionID),
			Radius:          data.Radius,
			SunTypeID:       sunTypeID,
			SecurityClass:   data.SecurityClass,
		}

		// Extract coordinates from position object
		if data.Position != nil {
			system.X = data.Position.X
			system.Y = data.Position.Y
			system.Z = data.Position.Z
		}

		systems = append(systems, system)
	}

	// Sort by solar system ID for consistent output
	sort.Slice(systems, func(i, j int) bool {
		return systems[i].SolarSystemID < systems[j].SolarSystemID
	})

	return systems, nil
}
