package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/jsonl"
)

// ParseTypes parses the types.jsonl file.
func (p *Parser) ParseTypes() (map[int64]models.SDEType, error) {
	path := p.filePath("types.jsonl")

	types, err := jsonl.ParseFileMap(path, models.SDEType.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse types file: %w", err)
	}

	return types, nil
}
