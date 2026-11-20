// Package jsonl reads the JSON Lines tables of the SDE API.
//
// Each line of a table is one record, and the record carries its own id in
// a named field (typeID, solarSystemID, ...), which the table's JSON Schema
// lists as x-sde-key.
package jsonl

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// ParseFileMap reads a JSONL table into a map keyed by the id that key
// returns for each record. An id that appears twice is an error. That also
// catches records without the id field, which all read as 0 (0 itself is a
// valid id: category, group and type 0 are "#System").
func ParseFileMap[V any](path string, key func(V) int64) (map[int64]V, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	result, err := ParseMap(f, key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return result, nil
}

// ParseMap is ParseFileMap for any reader.
func ParseMap[V any](r io.Reader, key func(V) int64) (map[int64]V, error) {
	result := make(map[int64]V)
	br := bufio.NewReaderSize(r, 1<<20)
	for n := 1; ; n++ {
		line, err := br.ReadBytes('\n')
		if len(line) > 0 && !isBlank(line) {
			var v V
			if uerr := json.Unmarshal(line, &v); uerr != nil {
				return nil, fmt.Errorf("line %d: failed to decode JSON: %w", n, uerr)
			}
			id := key(v)
			if _, dup := result[id]; dup {
				return nil, fmt.Errorf("line %d: duplicate id %d", n, id)
			}
			result[id] = v
		}
		if errors.Is(err, io.EOF) {
			return result, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func isBlank(line []byte) bool {
	for _, c := range line {
		if c != ' ' && c != '\t' && c != '\r' && c != '\n' {
			return false
		}
	}
	return true
}
