package parser

import (
	"fmt"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/jsonl"
)

// ParseCategories parses the categories.jsonl file.
func (p *Parser) ParseCategories() (map[int64]models.SDECategory, error) {
	path := p.filePath("categories.jsonl")

	categories, err := jsonl.ParseFileMap(path, models.SDECategory.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse categories file: %w", err)
	}

	return categories, nil
}
