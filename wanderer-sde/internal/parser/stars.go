package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// SDEMapStar represents a star in the flat SDE format.
type SDEMapStar struct {
	SolarSystemID int64         `yaml:"solarSystemID"`
	TypeID        int64         `yaml:"typeID"`
	Radius        float64       `yaml:"radius,omitempty"`
	Statistics    *SDEStarStats `yaml:"statistics,omitempty"`
}

// SDEStarStats represents star statistics.
type SDEStarStats struct {
	Age           float64 `yaml:"age,omitempty"`
	Life          float64 `yaml:"life,omitempty"`
	Luminosity    float64 `yaml:"luminosity,omitempty"`
	SpectralClass string  `yaml:"spectralClass,omitempty"`
	Temperature   float64 `yaml:"temperature,omitempty"`
}

// ParseStars parses mapStars.yaml and returns a map of starID -> typeID.
func (p *Parser) ParseStars() (map[int64]int64, error) {
	path := p.filePath("mapStars.yaml")

	rawStars, err := yaml.ParseFileMap[int64, SDEMapStar](path)
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
