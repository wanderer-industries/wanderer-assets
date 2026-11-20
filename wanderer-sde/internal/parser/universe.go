package parser

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/jsonl"
)

// SDEMapRegion represents a region in the flat SDE format.
type SDEMapRegion struct {
	RegionID        int64               `json:"regionID"`
	Name            map[string]string   `json:"name"`
	NameID          int64               `json:"nameID,omitempty"`
	DescriptionID   int64               `json:"descriptionID,omitempty"`
	FactionID       int64               `json:"factionID,omitempty"`
	Position        *models.SDEPosition `json:"position,omitempty"`
	NebulaID        int64               `json:"nebulaID,omitempty"`
	WormholeClassID int64               `json:"wormholeClassID,omitempty"`
}

// SDEMapConstellation represents a constellation in the flat SDE format.
type SDEMapConstellation struct {
	ConstellationID int64               `json:"constellationID"`
	RegionID        int64               `json:"regionID"`
	Name            map[string]string   `json:"name"`
	NameID          int64               `json:"nameID,omitempty"`
	FactionID       int64               `json:"factionID,omitempty"`
	Position        *models.SDEPosition `json:"position,omitempty"`
	Radius          float64             `json:"radius,omitempty"`
	WormholeClassID int64               `json:"wormholeClassID,omitempty"`
}

// SDEMapSolarSystem represents a solar system in the flat SDE format.
type SDEMapSolarSystem struct {
	SolarSystemID   int64               `json:"solarSystemID"`
	ConstellationID int64               `json:"constellationID"`
	RegionID        int64               `json:"regionID"`
	Name            map[string]string   `json:"name"`
	SecurityStatus  float64             `json:"securityStatus"`
	SecurityClass   string              `json:"securityClass,omitempty"`
	StarID          int64               `json:"starID,omitempty"`
	WormholeClassID int64               `json:"wormholeClassID,omitempty"`
	FactionID       int64               `json:"factionID,omitempty"`
	Border          bool                `json:"border,omitempty"`
	Corridor        bool                `json:"corridor,omitempty"`
	Fringe          bool                `json:"fringe,omitempty"`
	Hub             bool                `json:"hub,omitempty"`
	International   bool                `json:"international,omitempty"`
	Regional        bool                `json:"regional,omitempty"`
	Position        *models.SDEPosition `json:"position,omitempty"`
	Luminosity      float64             `json:"luminosity,omitempty"`
	Radius          float64             `json:"radius,omitempty"`
}

// ParseRegions parses the mapRegions.jsonl file.
func (p *Parser) ParseRegions() ([]models.Region, error) {
	path := p.filePath("mapRegions.jsonl")

	rawRegions, err := jsonl.ParseFileMap(path, SDEMapRegion.ID)
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

// ParseConstellations parses the mapConstellations.jsonl file.
func (p *Parser) ParseConstellations() ([]models.Constellation, error) {
	path := p.filePath("mapConstellations.jsonl")

	rawConstellations, err := jsonl.ParseFileMap(path, SDEMapConstellation.ID)
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

// ParseSolarSystems parses the mapSolarSystems.jsonl file.
// starTypeMap provides starID -> typeID mapping for resolving sun types.
func (p *Parser) ParseSolarSystems(starTypeMap map[int64]int64) ([]models.SolarSystem, error) {
	path := p.filePath("mapSolarSystems.jsonl")

	rawSystems, err := jsonl.ParseFileMap(path, SDEMapSolarSystem.ID)
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

// ID returns the record id, for jsonl.ParseFileMap.
func (r SDEMapRegion) ID() int64 { return r.RegionID }

// ID returns the record id, for jsonl.ParseFileMap.
func (c SDEMapConstellation) ID() int64 { return c.ConstellationID }

// ID returns the record id, for jsonl.ParseFileMap.
func (s SDEMapSolarSystem) ID() int64 { return s.SolarSystemID }
