package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// blueLootBuyerCorpIDs contains corporation IDs that buy blue loot.
// DED (Directive Enforcement Department) is the primary buyer.
var blueLootBuyerCorpIDs = map[int64]bool{
	1000137: true, // DED
}

// IsBlueLootBuyer returns true if the given corporation ID is a blue loot buyer.
func IsBlueLootBuyer(corpID int64) bool {
	return blueLootBuyerCorpIDs[corpID]
}

// ParseNPCStations parses the npcStations.yaml file.
func (p *Parser) ParseNPCStations() (map[int64]models.SDENPCStation, error) {
	path := p.filePath("npcStations.yaml")

	stations, err := yaml.ParseFileMap[int64, models.SDENPCStation](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NPC stations file: %w", err)
	}

	return stations, nil
}

// ParseNPCCorporations parses the npcCorporations.yaml file.
func (p *Parser) ParseNPCCorporations() (map[int64]models.SDENPCCorporation, error) {
	path := p.filePath("npcCorporations.yaml")

	corps, err := yaml.ParseFileMap[int64, models.SDENPCCorporation](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NPC corporations file: %w", err)
	}

	return corps, nil
}

// FilterBlueLootStations filters stations to only those owned by blue loot buyers.
func FilterBlueLootStations(stations map[int64]models.SDENPCStation) map[int64]models.SDENPCStation {
	filtered := make(map[int64]models.SDENPCStation)
	for stationID, station := range stations {
		if IsBlueLootBuyer(station.OwnerID) {
			filtered[stationID] = station
		}
	}
	return filtered
}
