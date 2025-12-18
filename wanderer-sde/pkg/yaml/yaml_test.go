package yaml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	// Create a temporary YAML file
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	yamlContent := `name: "test"
value: 42
nested:
  key: "value"
`
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	// Parse into struct
	type testStruct struct {
		Name   string `yaml:"name"`
		Value  int    `yaml:"value"`
		Nested struct {
			Key string `yaml:"key"`
		} `yaml:"nested"`
	}

	var result testStruct
	if err := ParseFile(yamlPath, &result); err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("Expected name 'test', got %q", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("Expected value 42, got %d", result.Value)
	}
	if result.Nested.Key != "value" {
		t.Errorf("Expected nested key 'value', got %q", result.Nested.Key)
	}
}

func TestParseFile_NotFound(t *testing.T) {
	var result struct{}
	err := ParseFile("/nonexistent/path.yaml", &result)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParseFile_MalformedYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create malformed YAML
	yamlPath := filepath.Join(tmpDir, "malformed.yaml")
	if err := os.WriteFile(yamlPath, []byte("[[invalid yaml:"), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	var result map[string]interface{}
	err = ParseFile(yamlPath, &result)
	if err == nil {
		t.Error("Expected error for malformed YAML")
	}
}

func TestParse(t *testing.T) {
	yamlContent := `key: value
number: 123
`
	reader := strings.NewReader(yamlContent)

	type testStruct struct {
		Key    string `yaml:"key"`
		Number int    `yaml:"number"`
	}

	var result testStruct
	if err := Parse(reader, &result); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Key != "value" {
		t.Errorf("Expected key 'value', got %q", result.Key)
	}
	if result.Number != 123 {
		t.Errorf("Expected number 123, got %d", result.Number)
	}
}

func TestParseFileMap(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create YAML with map structure (like typeIDs.yaml)
	yamlContent := `587:
  name: "Rifter"
  groupID: 25
588:
  name: "Slasher"
  groupID: 25
`
	yamlPath := filepath.Join(tmpDir, "types.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	type typeData struct {
		Name    string `yaml:"name"`
		GroupID int64  `yaml:"groupID"`
	}

	result, err := ParseFileMap[int64, typeData](yamlPath)
	if err != nil {
		t.Fatalf("ParseFileMap failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(result))
	}

	rifter, ok := result[587]
	if !ok {
		t.Error("Expected to find entry 587")
	}
	if rifter.Name != "Rifter" {
		t.Errorf("Expected name 'Rifter', got %q", rifter.Name)
	}
	if rifter.GroupID != 25 {
		t.Errorf("Expected groupID 25, got %d", rifter.GroupID)
	}
}

func TestParseFileMap_NotFound(t *testing.T) {
	type testStruct struct {
		Name string `yaml:"name"`
	}

	_, err := ParseFileMap[int64, testStruct]("/nonexistent/path.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParseFileMap_MalformedYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	yamlPath := filepath.Join(tmpDir, "malformed.yaml")
	if err := os.WriteFile(yamlPath, []byte("[[invalid yaml:"), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	type testStruct struct {
		Name string `yaml:"name"`
	}

	_, err = ParseFileMap[int64, testStruct](yamlPath)
	if err == nil {
		t.Error("Expected error for malformed YAML")
	}
}

func TestParseFileMap_EmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create empty YAML file
	yamlPath := filepath.Join(tmpDir, "empty.yaml")
	if err := os.WriteFile(yamlPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	type testStruct struct {
		Name string `yaml:"name"`
	}

	result, err := ParseFileMap[int64, testStruct](yamlPath)
	// Empty YAML should result in nil/empty map or error
	if err == nil && len(result) != 0 {
		t.Errorf("Empty file should return nil map or error, got map with %d entries", len(result))
	}
}

func TestParseFileMap_StringKeys(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	yamlContent := `region_a:
  name: "Region A"
region_b:
  name: "Region B"
`
	yamlPath := filepath.Join(tmpDir, "regions.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	type regionData struct {
		Name string `yaml:"name"`
	}

	result, err := ParseFileMap[string, regionData](yamlPath)
	if err != nil {
		t.Fatalf("ParseFileMap with string keys failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(result))
	}

	regionA, ok := result["region_a"]
	if !ok {
		t.Error("Expected to find entry 'region_a'")
	}
	if regionA.Name != "Region A" {
		t.Errorf("Expected name 'Region A', got %q", regionA.Name)
	}
}

func TestParseFileMap_NestedStructs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yaml_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	yamlContent := `1:
  name:
    en: "English Name"
    de: "German Name"
  stats:
    value: 100
`
	yamlPath := filepath.Join(tmpDir, "nested.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	type nestedData struct {
		Name  map[string]string `yaml:"name"`
		Stats struct {
			Value int `yaml:"value"`
		} `yaml:"stats"`
	}

	result, err := ParseFileMap[int64, nestedData](yamlPath)
	if err != nil {
		t.Fatalf("ParseFileMap with nested structs failed: %v", err)
	}

	entry, ok := result[1]
	if !ok {
		t.Error("Expected to find entry 1")
	}
	if entry.Name["en"] != "English Name" {
		t.Errorf("Expected English name, got %q", entry.Name["en"])
	}
	if entry.Stats.Value != 100 {
		t.Errorf("Expected stats value 100, got %d", entry.Stats.Value)
	}
}
