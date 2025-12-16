package parser

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// SDEStargateDestination represents the destination of a stargate.
type SDEStargateDestination struct {
	SolarSystemID int64 `yaml:"solarSystemID"`
	StargateID    int64 `yaml:"stargateID"`
}

// SDEMapStargate represents a stargate in the flat SDE format.
type SDEMapStargate struct {
	SolarSystemID int64                  `yaml:"solarSystemID"`
	Destination   SDEStargateDestination `yaml:"destination"`
	TypeID        int64                  `yaml:"typeID,omitempty"`
}

// ParseStargates parses the mapStargates.yaml file and extracts system jumps.
// Both directions of each stargate connection are included (A→B and B→A)
// to match Fuzzwork CSV format.
func (p *Parser) ParseStargates() ([]models.SystemJump, error) {
	path := p.filePath("mapStargates.yaml")

	rawStargates, err := yaml.ParseFileMap[int64, SDEMapStargate](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stargates file: %w", err)
	}

	// Use a map to deduplicate identical A→B entries (but keep A→B and B→A as separate)
	jumpSet := make(map[[2]int64]struct{})

	for _, data := range rawStargates {
		fromSystem := data.SolarSystemID
		toSystem := data.Destination.SolarSystemID

		if fromSystem == 0 || toSystem == 0 {
			// Invalid data, skip
			continue
		}

		// Store the jump as-is (preserving direction from stargate)
		pair := [2]int64{fromSystem, toSystem}
		jumpSet[pair] = struct{}{}
	}

	// Convert to slice
	jumps := make([]models.SystemJump, 0, len(jumpSet))
	for pair := range jumpSet {
		jumps = append(jumps, models.SystemJump{
			FromSolarSystemID: pair[0],
			ToSolarSystemID:   pair[1],
		})
	}

	// Sort by from system ID, then to system ID for consistent output
	sort.Slice(jumps, func(i, j int) bool {
		if jumps[i].FromSolarSystemID != jumps[j].FromSolarSystemID {
			return jumps[i].FromSolarSystemID < jumps[j].FromSolarSystemID
		}
		return jumps[i].ToSolarSystemID < jumps[j].ToSolarSystemID
	})

	return jumps, nil
}
