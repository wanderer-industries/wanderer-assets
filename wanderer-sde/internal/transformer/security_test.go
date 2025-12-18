package transformer

import (
	"math"
	"testing"
)

func TestGetTrueSecurity(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		// High sec systems
		{"Jita-like high sec", 0.9459, 0.9},
		{"High sec 0.5", 0.5, 0.5},
		{"High sec rounds down", 0.54, 0.5},
		{"High sec rounds up", 0.55, 0.6},

		// Low sec systems
		{"Low sec 0.4", 0.4, 0.4},
		{"Low sec 0.1", 0.1, 0.1},

		// Edge case: very low positive security
		{"Low positive rounds up", 0.047, 0.1},
		{"Low positive rounds up 2", 0.01, 0.1},
		{"Low positive rounds up 3", 0.049, 0.1},

		// Zero and negative (null sec)
		{"Zero security", 0.0, 0.0},
		{"Negative security", -0.5, -0.5},
		{"Negative rounds to -0.5", -0.45, -0.5},
		{"Negative rounds to -0.4", -0.449, -0.4},
		{"Negative rounds to -0.6", -0.55, -0.6},
		{"Deep null", -0.99, -1.0},

		// Boundary cases
		{"Exactly 0.05", 0.05, 0.1},
		{"Just below 0.05", 0.044, 0.1}, // Still > 0 and < 0.05, so rounds up
		{"Just above 0.05", 0.051, 0.1},

		// More edge cases
		{"0.95 stays 0.9", 0.94, 0.9},
		{"0.95 becomes 1.0", 0.95, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTrueSecurity(tt.input)
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("GetTrueSecurity(%f) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncateToTwoDigits(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.9459, 0.94},
		{0.9999, 0.99},
		{0.5555, 0.55},
		{0.1, 0.1},
		{0.0, 0.0},
		{-0.5555, -0.55}, // Truncation toward zero
	}

	for _, tt := range tests {
		result := TruncateToTwoDigits(tt.input)
		if math.Abs(result-tt.expected) > 0.001 {
			t.Errorf("TruncateToTwoDigits(%f) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}

func TestRoundSecurity(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.94, 0.9},
		{0.95, 1.0},
		{0.44, 0.4},
		{0.45, 0.5},
		{-0.44, -0.4},
		{-0.45, -0.5},
	}

	for _, tt := range tests {
		result := RoundSecurity(tt.input)
		if math.Abs(result-tt.expected) > 0.001 {
			t.Errorf("RoundSecurity(%f) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}
