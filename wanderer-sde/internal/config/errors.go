package config

import "errors"

var (
	// ErrNoSDESource is returned when neither SDE path nor download flag is set.
	ErrNoSDESource = errors.New("either --sde-path or --download must be specified")

	// ErrNoOutputDir is returned when no output directory is specified.
	ErrNoOutputDir = errors.New("output directory must be specified")
)
