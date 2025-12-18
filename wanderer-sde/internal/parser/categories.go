package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// ParseCategories parses the categories.yaml file.
func (p *Parser) ParseCategories() (map[int64]models.SDECategory, error) {
	path := p.filePath("categories.yaml")

	categories, err := yaml.ParseFileMap[int64, models.SDECategory](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse categories file: %w", err)
	}

	return categories, nil
}
