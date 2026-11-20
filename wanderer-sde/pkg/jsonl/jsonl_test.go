package jsonl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type record struct {
	TypeID  int64             `json:"typeID"`
	GroupID int64             `json:"groupID"`
	Name    map[string]string `json:"name"`
	Mass    float64           `json:"mass"`
}

func typeID(r record) int64 { return r.TypeID }

func TestParseFileMap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "types.jsonl")
	content := `{"typeID":0,"groupID":0,"name":{"en":"#System"}}
{"typeID":587,"groupID":25,"name":{"en":"Rifter","de":"Rifter"},"mass":1067000}
{"typeID":34,"groupID":18,"name":{"en":"Tritanium"},"mass":0.01}

{"typeID":22456,"groupID":25,"name":{"en":"Sabre"},"mass":2.82E+6}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ParseFileMap(path, typeID)
	if err != nil {
		t.Fatalf("ParseFileMap: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("got %d records, want 4", len(got))
	}
	if r, ok := got[0]; !ok || r.Name["en"] != "#System" {
		t.Errorf("id 0 = %+v, %v; want #System", r, ok)
	}
	if r := got[587]; r.GroupID != 25 || r.Name["en"] != "Rifter" || r.Mass != 1067000 {
		t.Errorf("587 = %+v", r)
	}
	if r := got[22456]; r.Mass != 2820000 {
		t.Errorf("22456 mass = %v, want 2820000", r.Mass)
	}
}

func TestParseMap_Errors(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"malformed", `{"typeID":1}` + "\n" + `{"typeID":`, "line 2: failed to decode JSON"},
		{"missing ids", `{"_key":1}` + "\n" + `{"_key":2}`, "line 2: duplicate id 0"},
		{"duplicate id", `{"typeID":1}` + "\n" + `{"typeID":1}`, "line 2: duplicate id 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseMap(strings.NewReader(tt.content), typeID)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestParseMap_Empty(t *testing.T) {
	got, err := ParseMap(strings.NewReader(""), typeID)
	if err != nil || len(got) != 0 {
		t.Errorf("got %v, %v; want an empty map", got, err)
	}
}

func TestParseFileMap_NotFound(t *testing.T) {
	if _, err := ParseFileMap("/nonexistent/types.jsonl", typeID); err == nil {
		t.Error("expected an error for a missing file")
	}
}
