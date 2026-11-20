// Package models defines data structures for SDE and Wanderer data.
package models

// SDERegion represents a region from the SDE universe files.
type SDERegion struct {
	RegionID        int64     `json:"regionID"`
	Center          []float64 `json:"center,omitempty"`
	Max             []float64 `json:"max,omitempty"`
	Min             []float64 `json:"min,omitempty"`
	NameID          int64     `json:"nameID,omitempty"`
	DescriptionID   int64     `json:"descriptionID,omitempty"`
	FactionID       int64     `json:"factionID,omitempty"`
	WormholeClassID int64     `json:"wormholeClassID,omitempty"`
}

// SDEConstellation represents a constellation from the SDE universe files.
type SDEConstellation struct {
	ConstellationID int64     `json:"constellationID"`
	Center          []float64 `json:"center,omitempty"`
	Max             []float64 `json:"max,omitempty"`
	Min             []float64 `json:"min,omitempty"`
	NameID          int64     `json:"nameID,omitempty"`
	FactionID       int64     `json:"factionID,omitempty"`
	Radius          float64   `json:"radius,omitempty"`
}

// SDESolarSystem represents a solar system from the SDE universe files.
type SDESolarSystem struct {
	SolarSystemID   int64                 `json:"solarSystemID"`
	Security        float64               `json:"security"`
	SecurityClass   string                `json:"securityClass,omitempty"`
	Center          []float64             `json:"center,omitempty"`
	Max             []float64             `json:"max,omitempty"`
	Min             []float64             `json:"min,omitempty"`
	Border          bool                  `json:"border,omitempty"`
	Corridor        bool                  `json:"corridor,omitempty"`
	FactionID       int64                 `json:"factionID,omitempty"`
	Fringe          bool                  `json:"fringe,omitempty"`
	Hub             bool                  `json:"hub,omitempty"`
	International   bool                  `json:"international,omitempty"`
	Luminosity      float64               `json:"luminosity,omitempty"`
	Radius          float64               `json:"radius,omitempty"`
	Regional        bool                  `json:"regional,omitempty"`
	Star            *SDEStar              `json:"star,omitempty"`
	Stargates       map[int64]SDEStargate `json:"stargates,omitempty"`
	Planets         map[int64]SDEPlanet   `json:"planets,omitempty"`
	SunTypeID       int64                 `json:"sunTypeID,omitempty"`
	WormholeClassID int64                 `json:"wormholeClassID,omitempty"`
}

// SDEStar represents star data within a solar system.
type SDEStar struct {
	ID         int64              `json:"id"`
	Radius     float64            `json:"radius,omitempty"`
	TypeID     int64              `json:"typeID"`
	Statistics *SDEStarStatistics `json:"statistics,omitempty"`
}

// SDEStarStatistics holds statistical data about a star.
type SDEStarStatistics struct {
	Age           float64 `json:"age,omitempty"`
	Life          float64 `json:"life,omitempty"`
	Locked        bool    `json:"locked,omitempty"`
	Luminosity    float64 `json:"luminosity,omitempty"`
	Radius        float64 `json:"radius,omitempty"`
	SpectralClass string  `json:"spectralClass,omitempty"`
	Temperature   float64 `json:"temperature,omitempty"`
}

// SDEStargate represents a stargate within a solar system.
type SDEStargate struct {
	Destination int64     `json:"destination"`
	Position    []float64 `json:"position,omitempty"`
	TypeID      int64     `json:"typeID,omitempty"`
}

// SDEPlanet represents a planet within a solar system.
type SDEPlanet struct {
	CelestialIndex   int64                `json:"celestialIndex,omitempty"`
	PlanetAttributes *SDEPlanetAttributes `json:"planetAttributes,omitempty"`
	Position         []float64            `json:"position,omitempty"`
	Radius           float64              `json:"radius,omitempty"`
	TypeID           int64                `json:"typeID,omitempty"`
	Moons            map[int64]SDEMoon    `json:"moons,omitempty"`
}

// SDEPlanetAttributes holds planet-specific attributes.
type SDEPlanetAttributes struct {
	HeightMap1   int64 `json:"heightMap1,omitempty"`
	HeightMap2   int64 `json:"heightMap2,omitempty"`
	Population   bool  `json:"population,omitempty"`
	ShaderPreset int64 `json:"shaderPreset,omitempty"`
}

// SDEMoon represents a moon orbiting a planet.
type SDEMoon struct {
	Position []float64 `json:"position,omitempty"`
	Radius   float64   `json:"radius,omitempty"`
	TypeID   int64     `json:"typeID,omitempty"`
}

// SDEType represents an item type from the types table.
type SDEType struct {
	TypeID         int64             `json:"typeID"`
	GroupID        int64             `json:"groupID"`
	Name           map[string]string `json:"name"`
	Description    map[string]string `json:"description,omitempty"`
	Mass           float64           `json:"mass,omitempty"`
	Volume         float64           `json:"volume,omitempty"`
	Capacity       float64           `json:"capacity,omitempty"`
	PortionSize    int64             `json:"portionSize,omitempty"`
	Published      bool              `json:"published"`
	MarketGroupID  int64             `json:"marketGroupID,omitempty"`
	GraphicID      int64             `json:"graphicID,omitempty"`
	IconID         int64             `json:"iconID,omitempty"`
	SoundID        int64             `json:"soundID,omitempty"`
	BasePrice      float64           `json:"basePrice,omitempty"`
	RaceID         int64             `json:"raceID,omitempty"`
	SofFactionName string            `json:"sofFactionName,omitempty"`
}

// SDEGroup represents an item group from the groups table.
type SDEGroup struct {
	GroupID              int64             `json:"groupID"`
	CategoryID           int64             `json:"categoryID"`
	Name                 map[string]string `json:"name"`
	Published            bool              `json:"published"`
	Anchorable           bool              `json:"anchorable,omitempty"`
	Anchored             bool              `json:"anchored,omitempty"`
	FittableNonSingleton bool              `json:"fittableNonSingleton,omitempty"`
	UseBasePrice         bool              `json:"useBasePrice,omitempty"`
	IconID               int64             `json:"iconID,omitempty"`
}

// SDECategory represents an item category from the categories table.
type SDECategory struct {
	CategoryID int64             `json:"categoryID"`
	Name       map[string]string `json:"name"`
	Published  bool              `json:"published"`
	IconID     int64             `json:"iconID,omitempty"`
}

// SDEWormholeClassLocation represents a location's wormhole class (not an SDE table: derived from regions, constellations and systems).
type SDEWormholeClassLocation struct {
	LocationID      int64 `json:"locationID"`
	WormholeClassID int64 `json:"wormholeClassID"`
}

// SDESystemJump represents a stargate connection (not an SDE table: derived from mapStargates).
type SDESystemJump struct {
	FromSolarSystemID int64 `json:"fromSolarSystemID"`
	ToSolarSystemID   int64 `json:"toSolarSystemID"`
}

// SDEPosition represents x,y,z coordinates in the SDE.
type SDEPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// SDENPCStation represents an NPC station from the npcStations table.
type SDENPCStation struct {
	StationID                int64       `json:"stationID"`
	CelestialIndex           int64       `json:"celestialIndex,omitempty"`
	OperationID              int64       `json:"operationID,omitempty"`
	OrbitID                  int64       `json:"orbitID,omitempty"`
	OrbitIndex               int64       `json:"orbitIndex,omitempty"`
	OwnerID                  int64       `json:"ownerID"`
	Position                 SDEPosition `json:"position,omitempty"`
	ReprocessingEfficiency   float64     `json:"reprocessingEfficiency,omitempty"`
	ReprocessingHangarFlag   int64       `json:"reprocessingHangarFlag,omitempty"`
	ReprocessingStationsTake float64     `json:"reprocessingStationsTake,omitempty"`
	SolarSystemID            int64       `json:"solarSystemID"`
	TypeID                   int64       `json:"typeID,omitempty"`
	UseOperationName         bool        `json:"useOperationName,omitempty"`
}

// SDENPCCorporation represents an NPC corporation from the npcCorporations table.
// Only includes fields needed for station lookup.
type SDENPCCorporation struct {
	CorporationID int64             `json:"corporationID"`
	Name          map[string]string `json:"name"`
	StationID     int64             `json:"stationID,omitempty"`
	Deleted       bool              `json:"deleted,omitempty"`
}

// ID returns the record id, for jsonl.ParseFileMap.
func (t SDEType) ID() int64 { return t.TypeID }

// ID returns the record id, for jsonl.ParseFileMap.
func (g SDEGroup) ID() int64 { return g.GroupID }

// ID returns the record id, for jsonl.ParseFileMap.
func (c SDECategory) ID() int64 { return c.CategoryID }

// ID returns the record id, for jsonl.ParseFileMap.
func (s SDENPCStation) ID() int64 { return s.StationID }

// ID returns the record id, for jsonl.ParseFileMap.
func (c SDENPCCorporation) ID() int64 { return c.CorporationID }
