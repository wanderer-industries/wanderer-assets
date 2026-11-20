package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/pkg/jsonl"
)

// SDEMapStar represents a star in the flat SDE format.
type SDEMapStar struct {
	StarID        int64         `json:"starID"`
	SolarSystemID int64         `json:"solarSystemID"`
	TypeID        int64         `json:"typeID"`
	Radius        float64       `json:"radius,omitempty"`
	Statistics    *SDEStarStats `json:"statistics,omitempty"`
}

// SDEStarStats represents star statistics.
type SDEStarStats struct {
	Age           float64 `json:"age,omitempty"`
	Life          float64 `json:"life,omitempty"`
	Luminosity    float64 `json:"luminosity,omitempty"`
	SpectralClass string  `json:"spectralClass,omitempty"`
	Temperature   float64 `json:"temperature,omitempty"`
}

// ParseStars parses mapStars.jsonl and returns a map of starID -> typeID.
func (p *Parser) ParseStars() (map[int64]int64, error) {
	path := p.filePath("mapStars.jsonl")

	rawStars, err := jsonl.ParseFileMap(path, SDEMapStar.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stars file: %w", err)
	}

	starTypeMap := make(map[int64]int64, len(rawStars))
	for starID, data := range rawStars {
		if data.TypeID != 0 {
			starTypeMap[starID] = data.TypeID
		}
	}

	return starTypeMap, nil
}

// ID returns the record id, for jsonl.ParseFileMap.
func (s SDEMapStar) ID() int64 { return s.StarID }
