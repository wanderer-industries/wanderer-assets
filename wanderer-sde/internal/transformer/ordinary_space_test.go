package transformer

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/internal/parser"
)

func TestTransformer_InferOrdinarySpace(t *testing.T) {
	tests := []struct {
		name     string
		security float64
		evidence string
		want     int64
	}{
		{"5-VKCN", -0.016492, "J", 9},
		{"Exordium", 1.0, "A", 7},
		{"minimum security", -1.0, "J", 9},
		{"zero security", 0.0, "J", 9},
		{"positive near zero", 0.000001, "D", 8},
		{"below high sec boundary", 0.449, "D", 8},
		{"high sec boundary", 0.45, "D", 7},
		{"no ordinary space evidence", -0.5, "", 0},
		{"whitespace is not evidence", -0.5, " \t\n", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := models.SolarSystem{
				SolarSystemID: 30003346, ConstellationID: 20000487, RegionID: 10000041,
				Security: tt.security, SecurityClass: tt.evidence,
			}
			result, err := New(&config.Config{}).Transform(&parser.ParseResult{
				SolarSystems: []models.SolarSystem{system},
			})
			if err != nil {
				t.Fatalf("Transform failed: %v", err)
			}
			if tt.want == 0 {
				if len(result.WormholeClasses) != 0 {
					t.Fatalf("expected no mapping without ordinary-space evidence, got %v", result.WormholeClasses)
				}
			} else {
				want := []models.WormholeClassLocation{{LocationID: 30003346, WormholeClassID: tt.want}}
				if !reflect.DeepEqual(result.WormholeClasses, want) {
					t.Fatalf("class mappings = %v, want %v", result.WormholeClasses, want)
				}
			}
			if got := result.Universe.SolarSystems[0].Security; got != tt.security {
				t.Errorf("raw security changed: got %v, want %v", got, tt.security)
			}
		})
	}
}

func TestTransformer_PreserveOrdinarySpaceAssignments(t *testing.T) {
	for _, locationID := range []int64{30003346, 20000487, 10000041} {
		for _, classID := range []int64{-1, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 10100} {
			t.Run(fmt.Sprintf("location_%d_class_%d", locationID, classID), func(t *testing.T) {
				classes := []models.WormholeClassLocation{{LocationID: locationID, WormholeClassID: classID}}
				result, err := New(&config.Config{}).Transform(&parser.ParseResult{
					WormholeClasses: classes,
					SolarSystems: []models.SolarSystem{{
						SolarSystemID: 30003346, ConstellationID: 20000487, RegionID: 10000041,
						Security: 1.0, SecurityClass: "A",
					}},
				})
				if err != nil {
					t.Fatalf("Transform failed: %v", err)
				}
				if !reflect.DeepEqual(result.WormholeClasses, classes) {
					t.Fatalf("explicit mapping was changed or overridden: got %v, want %v", result.WormholeClasses, classes)
				}
			})
		}
	}
}

func TestTransformer_DeterministicOrdinarySpaceAssignments(t *testing.T) {
	classes := make([]models.WormholeClassLocation, 2, 4)
	classes[0] = models.WormholeClassLocation{LocationID: 11000031, WormholeClassID: 12}
	classes[1] = models.WormholeClassLocation{LocationID: 10000002, WormholeClassID: 7}
	before := append([]models.WormholeClassLocation(nil), classes...)
	input := &parser.ParseResult{
		WormholeClasses: classes,
		SolarSystems: []models.SolarSystem{
			{SolarSystemID: 30100032, ConstellationID: 20001004, RegionID: 10001004, Security: 1.0, SecurityClass: "A"},
			{SolarSystemID: 30003346, ConstellationID: 20000487, RegionID: 10000041, Security: -0.016492, SecurityClass: "J"},
		},
	}
	want := []models.WormholeClassLocation{
		{LocationID: 10000002, WormholeClassID: 7},
		{LocationID: 11000031, WormholeClassID: 12},
		{LocationID: 30003346, WormholeClassID: 9},
		{LocationID: 30100032, WormholeClassID: 7},
	}
	tr := New(&config.Config{})
	for run := 0; run < 2; run++ {
		result, err := tr.Transform(input)
		if err != nil {
			t.Fatalf("Transform failed: %v", err)
		}
		if !reflect.DeepEqual(result.WormholeClasses, want) {
			t.Fatalf("run %d: mappings = %v, want %v", run, result.WormholeClasses, want)
		}
		if !reflect.DeepEqual(input.WormholeClasses, before) {
			t.Fatal("Transform mutated the parser's class assignments")
		}
	}
}

func TestTransformer_RejectInvalidOrdinarySpaceSecurity(t *testing.T) {
	for _, security := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1.01, 1.01} {
		t.Run(fmt.Sprint(security), func(t *testing.T) {
			result, err := New(&config.Config{}).Transform(&parser.ParseResult{
				SolarSystems: []models.SolarSystem{{
					SolarSystemID: 30003346, Security: security, SecurityClass: "J",
				}},
			})
			if err == nil {
				t.Fatalf("invalid security %v was accepted: %v", security, result.WormholeClasses)
			}
			if !strings.Contains(err.Error(), "30003346") || !strings.Contains(err.Error(), "security") {
				t.Errorf("error lacks system/security context: %v", err)
			}
			if result != nil {
				t.Error("invalid inference returned partially converted data")
			}
		})
	}
}
