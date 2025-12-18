package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
)

func TestExtractAllWormholeClasses(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	wormholeClasses, err := p.ExtractAllWormholeClasses()
	if err != nil {
		t.Fatalf("ExtractAllWormholeClasses failed: %v", err)
	}

	// Test data has:
	// - 2 regions with wormholeClassID 7 (The Forge: 10000002, Derelik: 10000001)
	// - 2 constellations with wormholeClassID 7 (Kimotoro: 20000020, Joas: 20000001)
	// - 1 solar system with wormholeClassID 3 (J123456: 31000001)
	// Total: 5 entries
	expectedCount := 5
	if len(wormholeClasses) != expectedCount {
		t.Errorf("Expected %d wormhole class entries (2 regions + 2 constellations + 1 system), got %d", expectedCount, len(wormholeClasses))
	}

	// Build a map for easier verification
	classMap := make(map[int64]int64)
	for _, wc := range wormholeClasses {
		classMap[wc.LocationID] = wc.WormholeClassID
	}

	// Verify region entries
	if classID, ok := classMap[10000001]; !ok || classID != 7 {
		t.Errorf("Expected region 10000001 (Derelik) to have wormholeClassID 7, got %d", classID)
	}
	if classID, ok := classMap[10000002]; !ok || classID != 7 {
		t.Errorf("Expected region 10000002 (The Forge) to have wormholeClassID 7, got %d", classID)
	}

	// Verify constellation entries
	if classID, ok := classMap[20000001]; !ok || classID != 7 {
		t.Errorf("Expected constellation 20000001 (Joas) to have wormholeClassID 7, got %d", classID)
	}
	if classID, ok := classMap[20000020]; !ok || classID != 7 {
		t.Errorf("Expected constellation 20000020 (Kimotoro) to have wormholeClassID 7, got %d", classID)
	}

	// Verify solar system entry (only J123456 has wormholeClassID in test data)
	if classID, ok := classMap[31000001]; !ok || classID != 3 {
		t.Errorf("Expected solar system 31000001 (J123456) to have wormholeClassID 3, got %d", classID)
	}

	// Verify sorting by location ID
	for i := 1; i < len(wormholeClasses); i++ {
		if wormholeClasses[i-1].LocationID >= wormholeClasses[i].LocationID {
			t.Errorf("Wormhole classes not sorted by LocationID: %d >= %d",
				wormholeClasses[i-1].LocationID, wormholeClasses[i].LocationID)
		}
	}
}

func TestExtractAllWormholeClasses_NoWormholes(t *testing.T) {
	tmpDir := createTestSDENoWormholes(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	wormholeClasses, err := p.ExtractAllWormholeClasses()
	if err != nil {
		t.Fatalf("ExtractAllWormholeClasses failed: %v", err)
	}

	// Should return empty slice when no wormhole classes exist
	if len(wormholeClasses) != 0 {
		t.Errorf("Expected 0 wormhole class entries for non-wormhole space, got %d", len(wormholeClasses))
	}
}

// createTestSDENoWormholes creates SDE data without any wormhole class IDs
func createTestSDENoWormholes(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "parser_test_no_wh")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Regions without wormholeClassID
	regionsYAML := `10000002:
  name:
    en: "The Forge"
  position:
    x: -96538397329247680
    y: 68904722523889856
    z: 103886273221498080
  factionID: 500001
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(regionsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapRegions.yaml: %v", err)
	}

	// Constellations without wormholeClassID
	constellationsYAML := `20000020:
  regionID: 10000002
  name:
    en: "Kimotoro"
  position:
    x: -107314934797574880
    y: 65893634706137696
    z: 106631148888006560
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapConstellations.yaml"), []byte(constellationsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapConstellations.yaml: %v", err)
	}

	// Solar systems without wormholeClassID
	systemsYAML := `30000142:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Jita"
  securityStatus: 0.9459
  position:
    x: -129500988494612512
    y: 60552325055663632
    z: 116970681498498304
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapSolarSystems.yaml"), []byte(systemsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapSolarSystems.yaml: %v", err)
	}

	return tmpDir
}
