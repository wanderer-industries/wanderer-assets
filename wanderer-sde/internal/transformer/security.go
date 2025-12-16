// Package transformer provides data transformation logic for converting SDE data
// to Wanderer's format.
package transformer

import "math"

// TruncateToTwoDigits truncates a float to 2 decimal places without rounding.
// Example: 0.456 -> 0.45, -0.789 -> -0.78
func TruncateToTwoDigits(value float64) float64 {
	return math.Trunc(value*100) / 100
}

// GetTrueSecurity calculates the display security status matching EVE Online's
// in-game display behavior. This is important because EVE uses non-standard
// rounding rules that players expect to see.
//
// EVE Online Security Display Rules:
//
//  1. Positive near-zero systems (0 < sec < 0.05) always display as 0.1
//     This prevents "0.0" from appearing for any system with positive security,
//     which would be misleading since true 0.0 indicates nullsec.
//     Example: A system with security 0.02 displays as 0.1, not 0.0
//
//  2. All other values use "round half up" at the second decimal place:
//     - Truncate to 2 decimal places first
//     - If the second decimal digit is >= 5, round up the first decimal
//     - If the second decimal digit is < 5, round down (truncate)
//     Examples:
//     0.45 -> 0.5 (second decimal 5, rounds up)
//     0.44 -> 0.4 (second decimal 4, rounds down)
//     0.849 -> 0.84 (truncated) -> 0.8 (second decimal 4, rounds down)
//     -0.45 -> -0.5 (negative values follow same rules)
//
// This matches the security status shown in-game in the system info panel
// and on the starmap.
func GetTrueSecurity(security float64) float64 {
	// Rule 1: Very low positive security always rounds up to 0.1
	// This ensures no positive-security system displays as "0.0"
	if security > 0.0 && security < 0.05 {
		return math.Ceil(security*10) / 10
	}

	// Rule 2: Standard EVE rounding for all other values
	// Step 1: Truncate to 2 decimal places
	truncated := TruncateToTwoDigits(security)

	// Step 2: Get the value at 1 decimal place (truncated toward zero)
	truncatedOneDecimal := math.Trunc(truncated*10) / 10

	// Step 3: Calculate the absolute second decimal digit's contribution
	diff := math.Abs(math.Round((truncated-truncatedOneDecimal)*100)) / 100

	// Step 4: Round based on the second decimal digit
	// If >= 0.05, round away from zero; otherwise keep the truncated value
	if diff < 0.05 {
		return truncatedOneDecimal
	}
	// Round away from zero: positive values round up, negative values round down
	if security >= 0 {
		return math.Ceil(truncated*10) / 10
	}
	return math.Floor(truncated*10) / 10
}

// RoundSecurity rounds security to one decimal place using standard rounding.
// This is a simpler alternative when EVE-specific rounding is not needed.
func RoundSecurity(security float64) float64 {
	return math.Round(security*10) / 10
}
