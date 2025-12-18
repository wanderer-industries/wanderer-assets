# Wanderer Assets

This repository contains assets for [Wanderer](https://github.com/wanderer-industries/wanderer), including images and the SDE converter tool.

## SDE Converter (wanderer-sde)

A Go tool to convert EVE Online's Static Data Export (SDE) YAML files into CSV or JSON format compatible with Wanderer and [Fuzzwork](https://www.fuzzwork.co.uk/dump/).

### Overview

This converter eliminates the dependency on third-party CSV dumps (like Fuzzwork) by processing the official SDE directly from CCP. It parses the YAML files and generates CSV (default) or JSON output that Wanderer can consume.

#### Features

- Downloads the latest SDE directly from CCP
- Parses YAML files with parallel processing for performance
- Generates CSV files (default) matching Fuzzwork format, or JSON
- Supports passthrough of community-maintained data files
- Version tracking to avoid redundant downloads
- Cross-platform support (Linux, macOS, Windows)

### Installation

#### From Source

```bash
git clone https://github.com/wanderer-industries/wanderer-assets.git
cd wanderer-assets
make build
```

The binary will be available at `bin/sdeconvert`.

#### Cross-Platform Builds

```bash
make build-all
```

This creates binaries for:
- Linux (amd64)
- macOS (Intel and ARM)
- Windows (amd64)

### Usage

#### Quick Start

Download and convert the latest SDE in one command:

```bash
./bin/sdeconvert --download --output ./output
```

#### Command Line Options

```
Usage:
  sdeconvert [flags]
  sdeconvert [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number

Flags:
  -d, --download             Download latest SDE from CCP
  -f, --format string        Output format: csv or json (default "csv")
  -h, --help                 help for sdeconvert
  -o, --output string        Output directory for output files (default "./output")
  -p, --passthrough string   Directory with Wanderer JSON files to copy
      --pretty               Pretty-print JSON output (only applies to JSON format) (default true)
  -s, --sde-path string      Path to SDE directory or ZIP file
      --sde-url string       URL to download SDE from
  -v, --verbose              Enable verbose output
  -w, --workers int          Number of parallel workers (default 4)
```

#### Usage Examples

##### Download and Convert Latest SDE (CSV format)

The simplest use case - download the latest SDE from CCP and convert it to CSV:

```bash
sdeconvert --download --output ./output
```

##### Convert to JSON Format

To output JSON instead of CSV:

```bash
sdeconvert --download --output ./output --format json
```

##### Convert an Existing SDE Directory

If you already have the SDE extracted locally:

```bash
sdeconvert --sde-path /path/to/sde --output ./output
```

##### Include Wanderer Passthrough Files

Some JSON files contain community-maintained data (wormhole info, effects, etc.) that should be copied as-is from the Wanderer repository:

```bash
sdeconvert --sde-path ./sde \
  --output ./output \
  --passthrough /path/to/wanderer/priv/repo/data
```

##### Verbose Mode with Custom Worker Count

For debugging or monitoring large conversions:

```bash
sdeconvert --download --output ./output --verbose --workers 8
```

##### Custom SDE URL

Use a specific SDE version or mirror:

```bash
sdeconvert --download \
  --sde-url "https://example.com/custom-sde.zip" \
  --output ./output
```

### Output Files

The converter generates the following files (CSV by default, JSON with `--format json`):

#### Generated from SDE

| File | Description | Source |
|------|-------------|--------|
| `mapSolarSystems.csv` | Solar systems with coordinates, security status, region/constellation IDs | `mapSolarSystems.yaml` |
| `mapRegions.csv` | Region ID, name, and coordinate bounds | `mapRegions.yaml` |
| `mapConstellations.csv` | Constellation ID, name, region, and coordinate bounds | `mapConstellations.yaml` |
| `mapLocationWormholeClasses.csv` | Wormhole class assignments for locations | `mapLocationWormholeClasses.yaml` |
| `invTypes.csv` | All item type definitions | `types.yaml` |
| `invGroups.csv` | All item group definitions | `groups.yaml` |
| `mapSolarSystemJumps.csv` | Stargate connections between systems | `mapStargates.yaml` |

#### Passthrough Files (Community-Maintained)

These files are copied from the Wanderer data directory when `--passthrough` is specified:

| File | Description |
|------|-------------|
| `wormholes.json` | Wormhole type definitions |
| `wormholeClasses.json` | Wormhole class definitions |
| `wormholeClassesInfo.json` | Detailed wormhole class information |
| `wormholeSystems.json` | Known wormhole system data |
| `triglavianSystems.json` | Triglavian invasion system data |
| `effects.json` | System effect definitions |
| `shatteredConstellations.json` | Shattered wormhole constellation data |
| `sunTypes.json` | Sun type definitions |
| `triglavianEffectsByFaction.json` | Triglavian effects by faction |

### Data Formats

The output format matches Fuzzwork's CSV dump format. When using `--format json`, the same data is output as JSON arrays.

#### Solar Systems (`mapSolarSystems.csv`)

CSV columns: `regionID`, `constellationID`, `solarSystemID`, `solarSystemName`, `x`, `y`, `z`, `xMin`, `xMax`, `yMin`, `yMax`, `zMin`, `zMax`, `luminosity`, `border`, `fringe`, `corridor`, `hub`, `international`, `regional`, `constellation`, `security`, `factionID`, `radius`, `sunTypeID`, `securityClass`

| Field | Type | Description |
|-------|------|-------------|
| `solarSystemID` | int64 | Unique system identifier |
| `regionID` | int64 | Parent region ID |
| `constellationID` | int64 | Parent constellation ID |
| `solarSystemName` | string | Display name of the system |
| `security` | float64 | Security status (-1.0 to 1.0) |
| `sunTypeID` | int64 | Type ID of the system's star (optional) |
| `securityClass` | string | Security class (A, B, C, etc.) |

#### Regions (`mapRegions.csv`)

CSV columns: `regionID`, `regionName`, `x`, `y`, `z`, `xMin`, `xMax`, `yMin`, `yMax`, `zMin`, `zMax`, `factionID`, `nebula`, `radius`

#### Constellations (`mapConstellations.csv`)

CSV columns: `regionID`, `constellationID`, `constellationName`, `x`, `y`, `z`, `xMin`, `xMax`, `yMin`, `yMax`, `zMin`, `zMax`, `factionID`, `radius`

#### Wormhole Classes (`mapLocationWormholeClasses.csv`)

CSV columns: `locationID`, `wormholeClassID`

Location IDs can be regions, constellations, or solar systems. Wormhole class IDs:
- 1-6: C1-C6 wormhole space
- 7: High-sec (0.5-1.0 security)
- 8: Low-sec (0.1-0.4 security)
- 9: Null-sec (0.0 and below)
- 12: Thera
- 13: Shattered wormholes
- 14-18: Drifter wormholes
- 25: Pochven (Triglavian space)

#### Item Types (`invTypes.csv`)

CSV columns: `typeID`, `groupID`, `typeName`, `description`, `mass`, `volume`, `capacity`, `portionSize`, `raceID`, `basePrice`, `published`, `marketGroupID`, `iconID`, `soundID`, `graphicID`

Contains all item types from the SDE.

#### Item Groups (`invGroups.csv`)

CSV columns: `groupID`, `categoryID`, `groupName`, `iconID`, `useBasePrice`, `anchored`, `anchorable`, `fittableNonSingleton`, `published`

Contains all item groups from the SDE.

#### System Jumps (`mapSolarSystemJumps.csv`)

CSV columns: `fromRegionID`, `fromConstellationID`, `fromSolarSystemID`, `toSolarSystemID`, `toConstellationID`, `toRegionID`

Represents stargate connections between solar systems.

### Development

#### Prerequisites

- Go 1.22 or later
- Make (optional, for convenience commands)

#### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to $GOPATH/bin
make install
```

#### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run integration tests (requires downloaded SDE)
go test -v ./internal/... -run Integration
```

#### Project Structure

```
wanderer-assets/
├── images/                        # Wanderer image assets
├── cmd/
│   └── sdeconvert/
│       └── main.go                # CLI entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── downloader/
│   │   ├── downloader.go          # SDE download & extraction
│   │   └── version.go             # Version checking
│   ├── models/
│   │   ├── sde.go                 # SDE data structures
│   │   ├── wanderer.go            # Output data structures
│   │   └── csv.go                 # CSV formatting helpers
│   ├── parser/
│   │   ├── parser.go              # Main parser orchestration
│   │   ├── universe.go            # Region/constellation/system parsing
│   │   ├── types.go               # types.yaml parsing
│   │   ├── groups.go              # groups.yaml parsing
│   │   ├── categories.go          # categories.yaml parsing
│   │   ├── jumps.go               # Stargate jump parsing
│   │   ├── stars.go               # Star type parsing
│   │   └── wormhole_classes.go    # Wormhole class parsing
│   ├── transformer/
│   │   ├── transformer.go         # Data transformation logic
│   │   ├── bounds.go              # Coordinate bounds calculation
│   │   ├── security.go            # Security status calculation
│   │   └── filters.go             # Category filtering
│   └── writer/
│       ├── writer.go              # Writer interface
│       ├── csv_writer.go          # CSV output generation
│       └── json_writer.go         # JSON output generation
├── pkg/
│   └── yaml/
│       └── yaml.go                # YAML utilities
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

#### Code Architecture

The converter follows a pipeline architecture:

1. **Downloader**: Downloads and extracts the SDE from CCP
2. **Parser**: Reads YAML files and converts to internal Go structs
3. **Transformer**: Applies business logic (bounds calculation, faction inheritance, sorting)
4. **Writer**: Serializes data to CSV or JSON files

Each component is isolated and testable independently.

### Data Sources

- **SDE**: [EVE Online Static Data Export](https://developers.eveonline.com/docs/services/static-data/)
- **Wanderer**: [wanderer-industries/wanderer](https://github.com/wanderer-industries/wanderer)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- CCP Games for providing the EVE Online Static Data Export
- The Wanderer project for the data format specifications
- Fuzzwork for the original CSV dump service that inspired this tool
