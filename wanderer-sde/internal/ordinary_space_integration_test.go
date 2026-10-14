package internal

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/internal/parser"
	"github.com/guarzo/wanderer-sde/internal/transformer"
	"github.com/guarzo/wanderer-sde/internal/writer"
)

func TestIntegration_OrdinarySpaceClassifications(t *testing.T) {
	sdeDir := t.TempDir()
	files := map[string]string{
		"mapRegions.yaml": `10000041:
  name: {en: Syndicate}
10001004:
  name: {en: Exordium}
10000070:
  name: {en: Pochven}
  wormholeClassID: 25
11000031:
  name: {en: G-R00031}
  wormholeClassID: 12
10000004:
  name: {en: UUA-F4}
`,
		"mapConstellations.yaml": `20000487:
  regionID: 10000041
20010001:
  regionID: 10001004
20000788:
  regionID: 10000070
21000324:
  regionID: 11000031
20000047:
  regionID: 10000004
`,
		"mapSolarSystems.yaml": `30003346:
  name: {en: 5-VKCN}
  regionID: 10000041
  constellationID: 20000487
  securityStatus: -0.016492
  securityClass: J
30100032:
  name: {en: Manifest}
  regionID: 10001004
  constellationID: 20010001
  securityStatus: 1.0
  securityClass: A
30002770:
  name: {en: Raravoss}
  regionID: 10000070
  constellationID: 20000788
  securityStatus: -1.0
  securityClass: C1
31000005:
  name: {en: Thera}
  regionID: 11000031
  constellationID: 21000324
  securityStatus: -0.99
30000326:
  name: {en: WF-1LM}
  regionID: 10000004
  constellationID: 20000047
  securityStatus: -0.0019
`,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(sdeDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	p := parser.New(&config.Config{}, sdeDir)
	systems, err := p.ParseSolarSystems(nil)
	if err != nil {
		t.Fatal(err)
	}
	classes, err := p.ExtractAllWormholeClasses()
	if err != nil {
		t.Fatal(err)
	}
	original := []models.WormholeClassLocation{
		{LocationID: 10000070, WormholeClassID: 25},
		{LocationID: 11000031, WormholeClassID: 12},
	}
	if !reflect.DeepEqual(classes, original) {
		t.Fatalf("parser must retain only explicit SDE mappings: got %v, want %v", classes, original)
	}
	want := append(append([]models.WormholeClassLocation(nil), original...),
		models.WormholeClassLocation{LocationID: 30003346, WormholeClassID: 9},
		models.WormholeClassLocation{LocationID: 30100032, WormholeClassID: 7},
	)

	for _, format := range []config.OutputFormat{config.FormatCSV, config.FormatJSON} {
		t.Run(string(format), func(t *testing.T) {
			cfg := &config.Config{OutputDir: t.TempDir(), OutputFormat: format}
			converted, err := transformer.New(cfg).Transform(&parser.ParseResult{
				SolarSystems: systems, WormholeClasses: classes,
			})
			if err != nil {
				t.Fatal(err)
			}
			w, err := writer.NewWriter(cfg)
			if err != nil {
				t.Fatal(err)
			}
			if err := w.WriteAll(converted); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(cfg.OutputDir, "mapLocationWormholeClasses."+string(format)))
			if err != nil {
				t.Fatal(err)
			}
			var got []models.WormholeClassLocation
			if format == config.FormatJSON {
				if err := json.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
			} else {
				rows, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(rows[0], []string{"locationID", "wormholeClassID"}) {
					t.Fatalf("unexpected class CSV header: %v", rows[0])
				}
				for _, row := range rows[1:] {
					location, err := strconv.ParseInt(row[0], 10, 64)
					if err != nil {
						t.Fatal(err)
					}
					class, err := strconv.ParseInt(row[1], 10, 64)
					if err != nil {
						t.Fatal(err)
					}
					got = append(got, models.WormholeClassLocation{LocationID: location, WormholeClassID: class})
				}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("published class mappings = %v, want %v", got, want)
			}
			if !reflect.DeepEqual(classes, original) {
				t.Fatal("normalization mutated the parser's explicit assignments")
			}
			for _, system := range converted.Universe.SolarSystems {
				if system.SolarSystemID == 30003346 && system.Security != -0.016492 {
					t.Fatalf("5-VKCN's raw security changed: %v", system.Security)
				}
			}
		})
	}
}
