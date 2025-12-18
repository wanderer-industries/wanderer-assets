// Package main provides the CLI entry point for the SDE converter.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/downloader"
	"github.com/guarzo/wanderer-sde/internal/parser"
	"github.com/guarzo/wanderer-sde/internal/transformer"
	"github.com/guarzo/wanderer-sde/internal/writer"
)

// Version is set at build time via ldflags.
var Version = "dev"

// SDEMetadata contains metadata about the SDE conversion.
type SDEMetadata struct {
	SDEVersion  string `json:"sde_version"`
	ReleaseDate string `json:"release_date"`
	GeneratedBy string `json:"generated_by"`
	GeneratedAt string `json:"generated_at"`
	Source      string `json:"source"`
}

// MetadataFileName is the name of the metadata output file.
const MetadataFileName = "sde_metadata.json"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var cfg = config.NewConfig()

var rootCmd = &cobra.Command{
	Use:   "sdeconvert",
	Short: "Convert EVE SDE to Wanderer data format",
	Long: `Converts EVE Online's Static Data Export (SDE) YAML files
into CSV or JSON format compatible with Wanderer/Fuzzwork.

This tool can download the latest SDE from CCP or use an existing
SDE directory, then parse the YAML files and generate output files
compatible with Wanderer's data format.`,
	Example: `  # Download latest SDE and convert to CSV (default)
  sdeconvert --download --output ./output

  # Convert to JSON format instead
  sdeconvert --download --output ./output --format json

  # Convert an existing SDE directory
  sdeconvert --sde-path ./sde --output ./output

  # Include Wanderer passthrough files (wormholes.json, etc.)
  sdeconvert --sde-path ./sde --output ./output --passthrough ../wanderer/priv/repo/data`,
	RunE: runConversion,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number and build information for sdeconvert.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sdeconvert version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().StringVarP(&cfg.SDEPath, "sde-path", "s", "", "Path to SDE directory or ZIP file")
	rootCmd.Flags().StringVarP(&cfg.OutputDir, "output", "o", "./output", "Output directory for output files")
	rootCmd.Flags().BoolVarP(&cfg.DownloadSDE, "download", "d", false, "Download latest SDE from CCP")
	rootCmd.Flags().StringVarP(&cfg.PassthroughDir, "passthrough", "p", "", "Directory with Wanderer JSON files to copy")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVar(&cfg.PrettyPrint, "pretty", true, "Pretty-print JSON output (only applies to JSON format)")
	rootCmd.Flags().StringVar(&cfg.SDEUrl, "sde-url", config.SDELatestURL, "URL to download SDE from")

	var formatStr string
	rootCmd.Flags().StringVarP(&formatStr, "format", "f", "csv", "Output format: csv or json (default: csv)")
	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		switch formatStr {
		case "csv":
			cfg.OutputFormat = config.FormatCSV
		case "json":
			cfg.OutputFormat = config.FormatJSON
		default:
			return fmt.Errorf("invalid format '%s': must be 'csv' or 'json'", formatStr)
		}
		return nil
	}
}

func runConversion(cmd *cobra.Command, args []string) error {
	// Set the version for User-Agent headers
	cfg.Version = Version

	if err := cfg.Validate(); err != nil {
		return err
	}

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)
	go func() {
		select {
		case <-sigChan:
			fmt.Println("\nInterrupt received, shutting down...")
			cancel()
		case <-ctx.Done():
			// Context already cancelled, exit goroutine
		}
	}()

	if cfg.Verbose {
		fmt.Println("Configuration:")
		fmt.Printf("  SDE Path:     %s\n", cfg.SDEPath)
		fmt.Printf("  Output Dir:   %s\n", cfg.OutputDir)
		fmt.Printf("  Output Format: %s\n", cfg.OutputFormat)
		fmt.Printf("  Download:     %v\n", cfg.DownloadSDE)
		fmt.Printf("  Passthrough:  %s\n", cfg.PassthroughDir)
	}

	sdePath := cfg.SDEPath
	var versionInfo *downloader.VersionInfo

	// Step 1: Download SDE if requested
	if cfg.DownloadSDE {
		dl := downloader.New(cfg)
		vc := downloader.NewVersionChecker(cfg)

		// Default SDE path when downloading is <output-dir>/sde
		if sdePath == "" {
			sdePath = filepath.Join(cfg.OutputDir, "sde")
		}

		// Check if update is needed
		var needsUpdate bool
		var err error
		needsUpdate, versionInfo, err = vc.NeedsUpdate(ctx, cfg.OutputDir)
		if err != nil {
			fmt.Printf("Warning: could not check SDE version: %v\n", err)
			needsUpdate = true // Proceed with download anyway
		}

		// Also check if the SDE directory exists
		if _, err := os.Stat(sdePath); os.IsNotExist(err) {
			needsUpdate = true
		}

		if needsUpdate {
			fmt.Println("Downloading latest SDE...")

			// Download to a temp location first
			downloadedPath, err := dl.DownloadAndExtract(ctx)
			if err != nil {
				return fmt.Errorf("failed to download SDE: %w", err)
			}

			// Remove old SDE directory if it exists
			if err := os.RemoveAll(sdePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Warning: could not remove old SDE: %v\n", err)
			}

			// Move downloaded SDE to the persistent location
			if err := os.Rename(downloadedPath, sdePath); err != nil {
				// If rename fails (cross-device), fall back to copy
				if err := copyDir(downloadedPath, sdePath); err != nil {
					// Clean up orphaned extracted directory before returning error
					_ = os.RemoveAll(downloadedPath)
					return fmt.Errorf("failed to move SDE to %s: %w", sdePath, err)
				}
				_ = os.RemoveAll(downloadedPath)
			}

			fmt.Printf("SDE downloaded and extracted to: %s\n", sdePath)

			// Store the version
			if versionInfo != nil {
				if err := vc.StoreVersion(cfg.OutputDir, versionInfo.BuildNumber); err != nil {
					fmt.Printf("Warning: could not store version: %v\n", err)
				}
			}
		} else {
			fmt.Println("SDE is up to date, using cached version")
		}
	}

	if sdePath == "" {
		return fmt.Errorf("no SDE path available")
	}

	// Validate the SDE structure
	dl := downloader.New(cfg)
	if err := dl.Validate(sdePath); err != nil {
		return fmt.Errorf("SDE validation failed: %w", err)
	}

	fmt.Printf("Using SDE at: %s\n", sdePath)

	// Step 2: Parse SDE YAML files
	p := parser.New(cfg, sdePath)
	parseResult, err := p.ParseAll()
	if err != nil {
		return fmt.Errorf("failed to parse SDE: %w", err)
	}

	fmt.Printf("\nParsing complete:\n")
	fmt.Printf("  Regions:         %d\n", len(parseResult.Regions))
	fmt.Printf("  Constellations:  %d\n", len(parseResult.Constellations))
	fmt.Printf("  Solar Systems:   %d\n", len(parseResult.SolarSystems))
	fmt.Printf("  Types:           %d\n", len(parseResult.Types))
	fmt.Printf("  Groups:          %d\n", len(parseResult.Groups))
	fmt.Printf("  Categories:      %d\n", len(parseResult.Categories))
	fmt.Printf("  Wormhole Classes: %d\n", len(parseResult.WormholeClasses))
	fmt.Printf("  System Jumps:    %d\n", len(parseResult.SystemJumps))

	// Step 3: Transform data
	t := transformer.New(cfg)
	convertedData, err := t.Transform(parseResult)
	if err != nil {
		return fmt.Errorf("failed to transform data: %w", err)
	}

	// Validate the converted data
	validationResult := t.Validate(convertedData)
	fmt.Printf("\nValidation results:\n")
	fmt.Printf("  Regions:         %d\n", validationResult.Regions)
	fmt.Printf("  Constellations:  %d\n", validationResult.Constellations)
	fmt.Printf("  Solar Systems:   %d\n", validationResult.SolarSystems)
	fmt.Printf("  Types:           %d\n", validationResult.InvTypes)
	fmt.Printf("  Groups:          %d\n", validationResult.InvGroups)
	fmt.Printf("  Wormhole Classes: %d\n", validationResult.WormholeClasses)
	fmt.Printf("  System Jumps:    %d\n", validationResult.SystemJumps)

	if len(validationResult.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range validationResult.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(validationResult.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range validationResult.Errors {
			fmt.Printf("  - %s\n", err)
		}
		return fmt.Errorf("validation failed with %d errors", len(validationResult.Errors))
	}

	// Step 4: Write output files
	w, err := writer.NewWriter(cfg)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}
	if err := w.WriteAll(convertedData); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Step 5: Write metadata file
	if versionInfo != nil {
		if err := writeMetadata(cfg.OutputDir, versionInfo); err != nil {
			fmt.Printf("Warning: could not write metadata file: %v\n", err)
		} else if cfg.Verbose {
			fmt.Printf("  Wrote %s\n", MetadataFileName)
		}
	}

	// Step 6: Copy passthrough files
	if cfg.PassthroughDir != "" {
		if err := w.CopyPassthroughFiles(cfg.PassthroughDir); err != nil {
			return fmt.Errorf("failed to copy passthrough files: %w", err)
		}
	}

	fmt.Printf("\nConversion complete! Output written to: %s\n", cfg.OutputDir)
	fmt.Printf("Generated files (%s format):\n", cfg.OutputFormat)

	outputFiles := writer.GetOutputFiles(cfg.OutputFormat)
	counts := []int{
		len(convertedData.Universe.SolarSystems),
		len(convertedData.Universe.Regions),
		len(convertedData.Universe.Constellations),
		len(convertedData.WormholeClasses),
		len(convertedData.InvTypes),
		len(convertedData.InvGroups),
		len(convertedData.SystemJumps),
		len(convertedData.NPCStations),
	}
	labels := []string{"systems", "regions", "constellations", "classes", "types", "groups", "jumps", "stations"}

	// Defensive length check to prevent index out of bounds
	minLen := len(outputFiles)
	if len(counts) < minLen {
		minLen = len(counts)
	}
	if len(labels) < minLen {
		minLen = len(labels)
	}

	for i := 0; i < minLen; i++ {
		fmt.Printf("  - %s (%d %s)\n", outputFiles[i], counts[i], labels[i])
	}

	return nil
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file, preserving permissions.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		_ = dstFile.Close()
		return err
	}
	return dstFile.Close()
}

// writeMetadata writes the SDE metadata file to the output directory.
func writeMetadata(outputDir string, versionInfo *downloader.VersionInfo) error {
	metadata := SDEMetadata{
		SDEVersion:  versionInfo.BuildNumber,
		ReleaseDate: versionInfo.ReleaseDate,
		GeneratedBy: "wanderer-sde",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Source:      "https://developers.eveonline.com/static-data",
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataPath := filepath.Join(outputDir, MetadataFileName)
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}
