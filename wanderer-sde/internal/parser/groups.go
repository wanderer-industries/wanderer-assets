package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/jsonl"
)

// ParseGroups parses the groups.jsonl file.
func (p *Parser) ParseGroups() (map[int64]models.SDEGroup, error) {
	path := p.filePath("groups.jsonl")

	groups, err := jsonl.ParseFileMap(path, models.SDEGroup.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse groups file: %w", err)
	}

	return groups, nil
}
