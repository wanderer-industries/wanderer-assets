// Package models defines data structures for SDE and Wanderer data.
package models

// SDERegion represents a region from the SDE universe files.
type SDERegion struct {
	RegionID        int64     `yaml:"regionID"`
	Center          []float64 `yaml:"center,omitempty"`
	Max             []float64 `yaml:"max,omitempty"`
	Min             []float64 `yaml:"min,omitempty"`
	NameID          int64     `yaml:"nameID,omitempty"`
	DescriptionID   int64     `yaml:"descriptionID,omitempty"`
	FactionID       int64     `yaml:"factionID,omitempty"`
	WormholeClassID int64     `yaml:"wormholeClassID,omitempty"`
}

// SDEConstellation represents a constellation from the SDE universe files.
type SDEConstellation struct {
	ConstellationID int64     `yaml:"constellationID"`
	Center          []float64 `yaml:"center,omitempty"`
	Max             []float64 `yaml:"max,omitempty"`
	Min             []float64 `yaml:"min,omitempty"`
	NameID          int64     `yaml:"nameID,omitempty"`
	FactionID       int64     `yaml:"factionID,omitempty"`
	Radius          float64   `yaml:"radius,omitempty"`
}

// SDESolarSystem represents a solar system from the SDE universe files.
type SDESolarSystem struct {
	SolarSystemID   int64                 `yaml:"solarSystemID"`
	Security        float64               `yaml:"security"`
	SecurityClass   string                `yaml:"securityClass,omitempty"`
	Center          []float64             `yaml:"center,omitempty"`
	Max             []float64             `yaml:"max,omitempty"`
	Min             []float64             `yaml:"min,omitempty"`
	Border          bool                  `yaml:"border,omitempty"`
	Corridor        bool                  `yaml:"corridor,omitempty"`
	FactionID       int64                 `yaml:"factionID,omitempty"`
	Fringe          bool                  `yaml:"fringe,omitempty"`
	Hub             bool                  `yaml:"hub,omitempty"`
	International   bool                  `yaml:"international,omitempty"`
	Luminosity      float64               `yaml:"luminosity,omitempty"`
	Radius          float64               `yaml:"radius,omitempty"`
	Regional        bool                  `yaml:"regional,omitempty"`
	Star            *SDEStar              `yaml:"star,omitempty"`
	Stargates       map[int64]SDEStargate `yaml:"stargates,omitempty"`
	Planets         map[int64]SDEPlanet   `yaml:"planets,omitempty"`
	SunTypeID       int64                 `yaml:"sunTypeID,omitempty"`
	WormholeClassID int64                 `yaml:"wormholeClassID,omitempty"`
}

// SDEStar represents star data within a solar system.
type SDEStar struct {
	ID         int64              `yaml:"id"`
	Radius     float64            `yaml:"radius,omitempty"`
	TypeID     int64              `yaml:"typeID"`
	Statistics *SDEStarStatistics `yaml:"statistics,omitempty"`
}

// SDEStarStatistics holds statistical data about a star.
type SDEStarStatistics struct {
	Age           float64 `yaml:"age,omitempty"`
	Life          float64 `yaml:"life,omitempty"`
	Locked        bool    `yaml:"locked,omitempty"`
	Luminosity    float64 `yaml:"luminosity,omitempty"`
	Radius        float64 `yaml:"radius,omitempty"`
	SpectralClass string  `yaml:"spectralClass,omitempty"`
	Temperature   float64 `yaml:"temperature,omitempty"`
}

// SDEStargate represents a stargate within a solar system.
type SDEStargate struct {
	Destination int64     `yaml:"destination"`
	Position    []float64 `yaml:"position,omitempty"`
	TypeID      int64     `yaml:"typeID,omitempty"`
}

// SDEPlanet represents a planet within a solar system.
type SDEPlanet struct {
	CelestialIndex   int64                `yaml:"celestialIndex,omitempty"`
	PlanetAttributes *SDEPlanetAttributes `yaml:"planetAttributes,omitempty"`
	Position         []float64            `yaml:"position,omitempty"`
	Radius           float64              `yaml:"radius,omitempty"`
	TypeID           int64                `yaml:"typeID,omitempty"`
	Moons            map[int64]SDEMoon    `yaml:"moons,omitempty"`
}

// SDEPlanetAttributes holds planet-specific attributes.
type SDEPlanetAttributes struct {
	HeightMap1   int64 `yaml:"heightMap1,omitempty"`
	HeightMap2   int64 `yaml:"heightMap2,omitempty"`
	Population   bool  `yaml:"population,omitempty"`
	ShaderPreset int64 `yaml:"shaderPreset,omitempty"`
}

// SDEMoon represents a moon orbiting a planet.
type SDEMoon struct {
	Position []float64 `yaml:"position,omitempty"`
	Radius   float64   `yaml:"radius,omitempty"`
	TypeID   int64     `yaml:"typeID,omitempty"`
}

// SDEType represents an item type from typeIDs.yaml.
type SDEType struct {
	GroupID        int64             `yaml:"groupID"`
	Name           map[string]string `yaml:"name"`
	Description    map[string]string `yaml:"description,omitempty"`
	Mass           float64           `yaml:"mass,omitempty"`
	Volume         float64           `yaml:"volume,omitempty"`
	Capacity       float64           `yaml:"capacity,omitempty"`
	PortionSize    int64             `yaml:"portionSize,omitempty"`
	Published      bool              `yaml:"published"`
	MarketGroupID  int64             `yaml:"marketGroupID,omitempty"`
	GraphicID      int64             `yaml:"graphicID,omitempty"`
	IconID         int64             `yaml:"iconID,omitempty"`
	SoundID        int64             `yaml:"soundID,omitempty"`
	BasePrice      float64           `yaml:"basePrice,omitempty"`
	RaceID         int64             `yaml:"raceID,omitempty"`
	SofFactionName string            `yaml:"sofFactionName,omitempty"`
}

// SDEGroup represents an item group from groupIDs.yaml.
type SDEGroup struct {
	CategoryID           int64             `yaml:"categoryID"`
	Name                 map[string]string `yaml:"name"`
	Published            bool              `yaml:"published"`
	Anchorable           bool              `yaml:"anchorable,omitempty"`
	Anchored             bool              `yaml:"anchored,omitempty"`
	FittableNonSingleton bool              `yaml:"fittableNonSingleton,omitempty"`
	UseBasePrice         bool              `yaml:"useBasePrice,omitempty"`
	IconID               int64             `yaml:"iconID,omitempty"`
}

// SDECategory represents an item category from categoryIDs.yaml.
type SDECategory struct {
	Name      map[string]string `yaml:"name"`
	Published bool              `yaml:"published"`
	IconID    int64             `yaml:"iconID,omitempty"`
}

// SDEWormholeClassLocation represents a location's wormhole class from mapLocationWormholeClasses.yaml.
type SDEWormholeClassLocation struct {
	LocationID      int64 `yaml:"locationID"`
	WormholeClassID int64 `yaml:"wormholeClassID"`
}

// SDESystemJump represents a stargate connection from mapSolarSystemJumps.yaml.
type SDESystemJump struct {
	FromSolarSystemID int64 `yaml:"fromSolarSystemID"`
	ToSolarSystemID   int64 `yaml:"toSolarSystemID"`
}

// SDEPosition represents x,y,z coordinates in the SDE.
type SDEPosition struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
	Z float64 `yaml:"z"`
}

// SDENPCStation represents an NPC station from npcStations.yaml.
type SDENPCStation struct {
	CelestialIndex           int64       `yaml:"celestialIndex,omitempty"`
	OperationID              int64       `yaml:"operationID,omitempty"`
	OrbitID                  int64       `yaml:"orbitID,omitempty"`
	OrbitIndex               int64       `yaml:"orbitIndex,omitempty"`
	OwnerID                  int64       `yaml:"ownerID"`
	Position                 SDEPosition `yaml:"position,omitempty"`
	ReprocessingEfficiency   float64     `yaml:"reprocessingEfficiency,omitempty"`
	ReprocessingHangarFlag   int64       `yaml:"reprocessingHangarFlag,omitempty"`
	ReprocessingStationsTake float64     `yaml:"reprocessingStationsTake,omitempty"`
	SolarSystemID            int64       `yaml:"solarSystemID"`
	TypeID                   int64       `yaml:"typeID,omitempty"`
	UseOperationName         bool        `yaml:"useOperationName,omitempty"`
}

// SDENPCCorporation represents an NPC corporation from npcCorporations.yaml.
// Only includes fields needed for station lookup.
type SDENPCCorporation struct {
	Name      map[string]string `yaml:"name"`
	StationID int64             `yaml:"stationID,omitempty"`
	Deleted   bool              `yaml:"deleted,omitempty"`
}
