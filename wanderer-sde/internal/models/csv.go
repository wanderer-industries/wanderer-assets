package models

import (
	"strconv"
	"strings"
)

// CSVHeaders defines the exact column headers for each CSV file.
// These are slimmed down to only the fields that Wanderer actually uses.
var CSVHeaders = map[string][]string{
	"mapSolarSystems": {
		"solarSystemID", "solarSystemName", "regionID", "constellationID",
		"security", "sunTypeID",
	},
	"mapRegions": {
		"regionID", "regionName",
	},
	"mapConstellations": {
		"constellationID", "constellationName",
	},
	"invTypes": {
		"typeID", "groupID", "typeName", "mass", "volume", "capacity",
	},
	"invGroups": {
		"groupID", "categoryID", "groupName",
	},
	"mapLocationWormholeClasses": {
		"locationID", "wormholeClassID",
	},
	"mapSolarSystemJumps": {
		"fromSolarSystemID", "toSolarSystemID",
	},
	"npcStations": {
		"stationID", "solarSystemID", "ownerID", "ownerName", "typeID",
	},
}

// FormatNullableInt64 formats an optional int64 for CSV output.
// Returns "None" if nil, otherwise the integer value.
func FormatNullableInt64(v *int64) string {
	if v == nil {
		return "None"
	}
	return strconv.FormatInt(*v, 10)
}

// FormatBool formats a boolean for CSV output.
// Returns "1" for true, "0" for false (Fuzzwork format).
func FormatBool(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// FormatFloat formats a float64 for CSV output.
// Uses full precision without scientific notation for coordinates.
func FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// FormatSecurity formats security status for CSV output.
// Uses 'f' format to avoid scientific notation and ensures decimal point
// for Elixir binary_to_float compatibility.
func FormatSecurity(v float64) string {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	// Ensure decimal point for Elixir binary_to_float compatibility
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

// Int64PtrNonZero returns a pointer to an int64 value, or nil if the value is 0.
func Int64PtrNonZero(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

// Int64PtrAlways returns a pointer to an int64 value, even if 0.
func Int64PtrAlways(v int64) *int64 {
	return &v
}

// ToCSVRow converts a SolarSystem to a CSV row with only fields Wanderer uses.
func (s *SolarSystem) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(s.SolarSystemID, 10),
		s.SolarSystemName,
		strconv.FormatInt(s.RegionID, 10),
		strconv.FormatInt(s.ConstellationID, 10),
		FormatSecurity(s.Security),
		FormatNullableInt64(s.SunTypeID),
	}
}

// ToCSVRow converts a Region to a CSV row with only fields Wanderer uses.
func (r *Region) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(r.RegionID, 10),
		r.RegionName,
	}
}

// ToCSVRow converts a Constellation to a CSV row with only fields Wanderer uses.
func (c *Constellation) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(c.ConstellationID, 10),
		c.ConstellationName,
	}
}

// ToCSVRow converts an InvType to a CSV row with only fields Wanderer uses.
func (t *InvType) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(t.TypeID, 10),
		strconv.FormatInt(t.GroupID, 10),
		t.TypeName,
		FormatFloat(t.Mass),
		FormatFloat(t.Volume),
		FormatFloat(t.Capacity),
	}
}

// ToCSVRow converts an InvGroup to a CSV row with only fields Wanderer uses.
func (g *InvGroup) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(g.GroupID, 10),
		strconv.FormatInt(g.CategoryID, 10),
		g.GroupName,
	}
}

// ToCSVRow converts a WormholeClassLocation to a CSV row.
func (w *WormholeClassLocation) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(w.LocationID, 10),
		strconv.FormatInt(w.WormholeClassID, 10),
	}
}

// ToCSVRow converts a SystemJump to a CSV row with only fields Wanderer uses.
func (j *SystemJump) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(j.FromSolarSystemID, 10),
		strconv.FormatInt(j.ToSolarSystemID, 10),
	}
}

// ToCSVRow converts an NPCStation to a CSV row.
func (s *NPCStation) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(s.StationID, 10),
		strconv.FormatInt(s.SolarSystemID, 10),
		strconv.FormatInt(s.OwnerID, 10),
		s.OwnerName,
		strconv.FormatInt(s.TypeID, 10),
	}
}
