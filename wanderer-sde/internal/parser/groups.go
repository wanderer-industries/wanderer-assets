package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// ParseGroups parses the groups.yaml file.
func (p *Parser) ParseGroups() (map[int64]models.SDEGroup, error) {
	path := p.filePath("groups.yaml")

	groups, err := yaml.ParseFileMap[int64, models.SDEGroup](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse groups file: %w", err)
	}

	return groups, nil
}
