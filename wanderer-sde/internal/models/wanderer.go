package models

// SolarSystem represents a solar system in Wanderer's format.
// Fields match Fuzzwork CSV column order for mapSolarSystems.csv.
type SolarSystem struct {
	RegionID        int64   `json:"regionID"`
	ConstellationID int64   `json:"constellationID"`
	SolarSystemID   int64   `json:"solarSystemID"`
	SolarSystemName string  `json:"solarSystemName"`
	X               float64 `json:"x"`
	Y               float64 `json:"y"`
	Z               float64 `json:"z"`
	XMin            float64 `json:"xMin"`
	XMax            float64 `json:"xMax"`
	YMin            float64 `json:"yMin"`
	YMax            float64 `json:"yMax"`
	ZMin            float64 `json:"zMin"`
	ZMax            float64 `json:"zMax"`
	Luminosity      float64 `json:"luminosity"`
	Border          bool    `json:"border"`
	Fringe          bool    `json:"fringe"`
	Corridor        bool    `json:"corridor"`
	Hub             bool    `json:"hub"`
	International   bool    `json:"international"`
	Regional        bool    `json:"regional"`
	Constellation   string  `json:"constellation"` // Always "None" - legacy field
	Security        float64 `json:"security"`
	FactionID       *int64  `json:"factionID,omitempty"` // Pointer to allow "None" in CSV
	Radius          float64 `json:"radius"`
	SunTypeID       *int64  `json:"sunTypeID,omitempty"` // Pointer to allow "None" in CSV
	SecurityClass   string  `json:"securityClass,omitempty"`
}

// Region represents a region in Wanderer's format.
// Fields match Fuzzwork CSV column order for mapRegions.csv.
type Region struct {
	RegionID   int64   `json:"regionID"`
	RegionName string  `json:"regionName"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Z          float64 `json:"z"`
	XMin       float64 `json:"xMin"`
	XMax       float64 `json:"xMax"`
	YMin       float64 `json:"yMin"`
	YMax       float64 `json:"yMax"`
	ZMin       float64 `json:"zMin"`
	ZMax       float64 `json:"zMax"`
	FactionID  *int64  `json:"factionID,omitempty"` // Pointer to allow "None" in CSV
	Nebula     int64   `json:"nebula"`              // Not in SDE, use 0
	Radius     float64 `json:"radius"`
}

// Constellation represents a constellation in Wanderer's format.
// Fields match Fuzzwork CSV column order for mapConstellations.csv.
type Constellation struct {
	RegionID          int64   `json:"regionID"`
	ConstellationID   int64   `json:"constellationID"`
	ConstellationName string  `json:"constellationName"`
	X                 float64 `json:"x"`
	Y                 float64 `json:"y"`
	Z                 float64 `json:"z"`
	XMin              float64 `json:"xMin"`
	XMax              float64 `json:"xMax"`
	YMin              float64 `json:"yMin"`
	YMax              float64 `json:"yMax"`
	ZMin              float64 `json:"zMin"`
	ZMax              float64 `json:"zMax"`
	FactionID         *int64  `json:"factionID,omitempty"` // Pointer to allow "None" in CSV
	Radius            float64 `json:"radius"`
}

// WormholeClassLocation represents a wormhole class assignment in Wanderer's format.
type WormholeClassLocation struct {
	LocationID      int64 `json:"locationID"`
	WormholeClassID int64 `json:"wormholeClassID"`
}

// InvType represents an item type in Wanderer's format.
// Fields match Fuzzwork CSV column order for invTypes.csv.
type InvType struct {
	TypeID        int64   `json:"typeID"`
	GroupID       int64   `json:"groupID"`
	TypeName      string  `json:"typeName"`
	Description   string  `json:"description"`
	Mass          float64 `json:"mass"`
	Volume        float64 `json:"volume"`
	Capacity      float64 `json:"capacity"`
	PortionSize   int64   `json:"portionSize"`
	RaceID        *int64  `json:"raceID,omitempty"` // Pointer to allow "None" in CSV
	BasePrice     float64 `json:"basePrice"`
	Published     bool    `json:"published"`
	MarketGroupID *int64  `json:"marketGroupID,omitempty"` // Pointer to allow "None" in CSV
	IconID        *int64  `json:"iconID,omitempty"`        // Pointer to allow "None" in CSV
	SoundID       *int64  `json:"soundID,omitempty"`       // Pointer to allow "None" in CSV
	GraphicID     *int64  `json:"graphicID,omitempty"`     // Pointer to allow "None" in CSV
}

// ShipType is an alias for backward compatibility.
// Deprecated: Use InvType instead.
type ShipType = InvType

// InvGroup represents an item group in Wanderer's format.
// Fields match Fuzzwork CSV column order for invGroups.csv.
type InvGroup struct {
	GroupID              int64  `json:"groupID"`
	CategoryID           int64  `json:"categoryID"`
	GroupName            string `json:"groupName"`
	IconID               *int64 `json:"iconID,omitempty"` // Pointer to allow "None" in CSV
	UseBasePrice         bool   `json:"useBasePrice"`
	Anchored             bool   `json:"anchored"`
	Anchorable           bool   `json:"anchorable"`
	FittableNonSingleton bool   `json:"fittableNonSingleton"`
	Published            bool   `json:"published"`
}

// ItemGroup is an alias for backward compatibility.
// Deprecated: Use InvGroup instead.
type ItemGroup = InvGroup

// SystemJump represents a stargate connection in Wanderer's format.
// Fields match Fuzzwork CSV column order for mapSolarSystemJumps.csv.
type SystemJump struct {
	FromRegionID        int64 `json:"fromRegionID"`
	FromConstellationID int64 `json:"fromConstellationID"`
	FromSolarSystemID   int64 `json:"fromSolarSystemID"`
	ToSolarSystemID     int64 `json:"toSolarSystemID"`
	ToConstellationID   int64 `json:"toConstellationID"`
	ToRegionID          int64 `json:"toRegionID"`
}

// NPCStation represents an NPC station in Wanderer's format.
// Used to identify stations where blue loot can be sold.
type NPCStation struct {
	StationID     int64  `json:"stationID"`
	SolarSystemID int64  `json:"solarSystemID"`
	OwnerID       int64  `json:"ownerID"`
	OwnerName     string `json:"ownerName"`
	TypeID        int64  `json:"typeID"`
}

// UniverseData holds all parsed universe data.
type UniverseData struct {
	Regions        []Region
	Constellations []Constellation
	SolarSystems   []SolarSystem
}

// ConvertedData holds all data ready for output.
type ConvertedData struct {
	Universe        *UniverseData
	InvTypes        []InvType
	InvGroups       []InvGroup
	WormholeClasses []WormholeClassLocation
	SystemJumps     []SystemJump
	NPCStations     []NPCStation
}

// ShipTypes returns InvTypes for backward compatibility.
// Deprecated: Use InvTypes instead.
func (c *ConvertedData) ShipTypes() []ShipType {
	return c.InvTypes
}

// ItemGroups returns InvGroups for backward compatibility.
// Deprecated: Use InvGroups instead.
func (c *ConvertedData) ItemGroups() []ItemGroup {
	return c.InvGroups
}

// ValidationResult holds the results of data validation.
type ValidationResult struct {
	SolarSystems    int
	Regions         int
	Constellations  int
	InvTypes        int
	InvGroups       int
	SystemJumps     int
	WormholeClasses int
	NPCStations     int
	Errors          []string
	Warnings        []string
}

// IsValid returns true if validation found no errors.
func (v *ValidationResult) IsValid() bool {
	return len(v.Errors) == 0
}
