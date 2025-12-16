package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// ParseTypes parses the types.yaml file.
func (p *Parser) ParseTypes() (map[int64]models.SDEType, error) {
	path := p.filePath("types.yaml")

	types, err := yaml.ParseFileMap[int64, models.SDEType](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse types file: %w", err)
	}

	return types, nil
}
