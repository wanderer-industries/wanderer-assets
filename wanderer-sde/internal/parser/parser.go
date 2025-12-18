// Package parser provides YAML parsing functionality for the EVE SDE.
package parser

import (
	"fmt"
	"path/filepath"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

// Parser orchestrates parsing of all SDE files.
type Parser struct {
	config  *config.Config
	sdePath string
}

// New creates a new Parser with the given configuration.
func New(cfg *config.Config, sdePath string) *Parser {
	return &Parser{
		config:  cfg,
		sdePath: sdePath,
	}
}

// ParseResult contains all parsed data from the SDE.
type ParseResult struct {
	Regions         []models.Region
	Constellations  []models.Constellation
	SolarSystems    []models.SolarSystem
	Types           map[int64]models.SDEType
	Groups          map[int64]models.SDEGroup
	Categories      map[int64]models.SDECategory
	WormholeClasses []models.WormholeClassLocation
	SystemJumps     []models.SystemJump
	NPCStations     map[int64]models.SDENPCStation
	NPCCorporations map[int64]models.SDENPCCorporation
}

// ParseAll parses all SDE files and returns the combined result.
func (p *Parser) ParseAll() (*ParseResult, error) {
	result := &ParseResult{}

	if p.config.Verbose {
		fmt.Println("Parsing SDE files...")
	}

	// Parse categories first (needed for filtering ships)
	if p.config.Verbose {
		fmt.Println("  Parsing categories...")
	}
	categories, err := p.ParseCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to parse categories: %w", err)
	}
	result.Categories = categories

	// Parse groups (needed for filtering ships)
	if p.config.Verbose {
		fmt.Println("  Parsing groups...")
	}
	groups, err := p.ParseGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to parse groups: %w", err)
	}
	result.Groups = groups

	// Parse types
	if p.config.Verbose {
		fmt.Println("  Parsing types...")
	}
	types, err := p.ParseTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to parse types: %w", err)
	}
	result.Types = types

	// Parse regions
	if p.config.Verbose {
		fmt.Println("  Parsing regions...")
	}
	regions, err := p.ParseRegions()
	if err != nil {
		return nil, fmt.Errorf("failed to parse regions: %w", err)
	}
	result.Regions = regions

	// Parse constellations
	if p.config.Verbose {
		fmt.Println("  Parsing constellations...")
	}
	constellations, err := p.ParseConstellations()
	if err != nil {
		return nil, fmt.Errorf("failed to parse constellations: %w", err)
	}
	result.Constellations = constellations

	// Parse stars (needed for solar system sun type resolution)
	if p.config.Verbose {
		fmt.Println("  Parsing stars...")
	}
	starTypeMap, err := p.ParseStars()
	if err != nil {
		return nil, fmt.Errorf("failed to parse stars: %w", err)
	}

	// Parse solar systems with star type lookup
	if p.config.Verbose {
		fmt.Println("  Parsing solar systems...")
	}
	systems, err := p.ParseSolarSystems(starTypeMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse solar systems: %w", err)
	}
	result.SolarSystems = systems

	// Parse stargates (system jumps)
	if p.config.Verbose {
		fmt.Println("  Parsing stargates...")
	}
	jumps, err := p.ParseStargates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse stargates: %w", err)
	}
	result.SystemJumps = jumps

	// Extract wormhole classes from regions, constellations, and systems
	if p.config.Verbose {
		fmt.Println("  Extracting wormhole classes...")
	}
	wormholeClasses, err := p.ExtractAllWormholeClasses()
	if err != nil {
		return nil, fmt.Errorf("failed to extract wormhole classes: %w", err)
	}
	result.WormholeClasses = wormholeClasses

	// Parse NPC stations (for blue loot buyer locations)
	if p.config.Verbose {
		fmt.Println("  Parsing NPC stations...")
	}
	npcStations, err := p.ParseNPCStations()
	if err != nil {
		return nil, fmt.Errorf("failed to parse NPC stations: %w", err)
	}
	// Filter to only blue loot buyer stations
	result.NPCStations = FilterBlueLootStations(npcStations)

	// Parse NPC corporations (for name lookup)
	if p.config.Verbose {
		fmt.Println("  Parsing NPC corporations...")
	}
	npcCorps, err := p.ParseNPCCorporations()
	if err != nil {
		return nil, fmt.Errorf("failed to parse NPC corporations: %w", err)
	}
	result.NPCCorporations = npcCorps

	if p.config.Verbose {
		fmt.Printf("Parsing complete:\n")
		fmt.Printf("  Regions:        %d\n", len(result.Regions))
		fmt.Printf("  Constellations: %d\n", len(result.Constellations))
		fmt.Printf("  Solar Systems:  %d\n", len(result.SolarSystems))
		fmt.Printf("  Types:          %d\n", len(result.Types))
		fmt.Printf("  Groups:         %d\n", len(result.Groups))
		fmt.Printf("  Categories:     %d\n", len(result.Categories))
		fmt.Printf("  Wormhole Classes: %d\n", len(result.WormholeClasses))
		fmt.Printf("  System Jumps:   %d\n", len(result.SystemJumps))
		fmt.Printf("  NPC Stations:   %d\n", len(result.NPCStations))
	}

	return result, nil
}

// filePath returns the full path to an SDE file.
func (p *Parser) filePath(filename string) string {
	return filepath.Join(p.sdePath, filename)
}
