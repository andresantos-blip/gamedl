# GameDL

A command-line tool for downloading and analyzing game payloads for betgenius and sportradar.

## Running gamedl

To run gamedl, you can either download a pre-built binary or build from source.
Using a pre-built binary is the easiest way to get started as it does not require you to have Go installed on your machine.

### Option 1: Automatic download and install

The easiest way to download/update and install the script is by running:

```bash
curl -sSfL https://raw.githubusercontent.com/andresantos-blip/gamedl/main/scripts/autoinstall.sh | bash
```

This command will automatically fetch the binary from the latest release and try to put it into you PATH so you can use it from any directory!

### Option 2: Manual Download

**Download v0.2.0:** [Mac (Intel)](https://github.com/andresantos-blip/gamedl/releases/download/v0.2.0/gamedl_Darwin_x86_64.zip) | [Mac (Apple Silicon)](https://github.com/andresantos-blip/gamedl/releases/download/v0.2.0/gamedl_Darwin_arm64.zip)

Use the download links above if you want more control of the installation process.

After downloading, extract the zip folder and run the binary inside it from a terminal:

```shell
unzip path/to/downloaded/gamedl_Darwin_x86_64.zip
cd path/to/extracted/folder
./gamedl --help
```

To be able to run `gamedl` from any directory, move the binary to a directory included in your system's PATH.

For more operating systems other than macOS, check the [releases page](https://github.com/andresantos-blip/gamedl/releases) for all the available builds.

### Option 3: Build from source

You can also build gamedl from source if you have Go installed on your machine.

#### Prerequisites

- Go 1.25 or later
  - `brew install go`

#### Building

```bash
# Clone the repository
git clone <repository-url>
cd gamedl

# Install dependencies
go mod tidy

# Build the CLI
go build -o gamedl .
```

## Setup

### Environment Variables

To download games, the tool requires different API credentials depending on the provider you want to use.
The env vars for each provider are mentioned below.

It's a good idea to set these if you plan to use the download command.
Feel free to ignore this section if you are only using the analysis command.

#### BetGenius

```bash
export BG_FIXTURE_KEY="your_betgenius_fixture_api_key"
export BG_FIXTURE_USER="your_betgenius_fixture_username"
export BG_FIXTURE_PASSWORD="your_betgenius_fixture_password"
export BG_STATS_KEY="your_betgenius_stats_api_key"
export BG_STATS_USER="your_betgenius_stats_username"
export BG_STATS_PASSWORD="your_betgenius_stats_password"
```

#### SportRadar (NCAAB, NCAAF)

```bash
export SPORTRADAR_NCAAB_KEY="your_sportradar_ncaab_api_key"
export SPORTRADAR_NCAAF_KEY="your_sportradar_ncaaf_api_key"
```

## Usage

### Download Command

Download game data from sports providers:

```bash
# Basic usage
./gamedl download --competition nfl --provider betgenius --seasons 2024

# Download multiple seasons
./gamedl download --competition ncaab --provider sr --seasons 2023,2024

# With custom output directory and concurrency
./gamedl download --competition nfl --provider bg --seasons 2024 --output-dir ./my_data --concurrency 4

# Using environment variables
GAMEDL_DOWNLOAD_COMPETITION=nfl GAMEDL_DOWNLOAD_PROVIDER=betgenius ./gamedl download --seasons 2024
```

#### Download Options

- `--competition, -c`: Competition to download (values allowed: 'nfl', 'ncaab' or 'ncaaf') **(required)**
- `--provider, -p`: Data provider (values allowed: 'sportradar', 'sr', 'betgenius', 'genius' or 'bg') **(required)**
- `--seasons, -s`: Seasons to download, comma-separated. e.g '2023,2024' (default: all seasons available in the provider)
- `--output-dir, -o`: Directory to store downloaded game files (default: downloaded_games")
- `--concurrency`: Number of concurrent downloads (default: 10)

#### Supported Combinations

| Competition | BetGenius | SportRadar |
|-------------|-----------|------------|
| NFL         | ✅        | ❌         |
| NCAAB       | ❌        | ✅         |
| NCAAF       | ❌        | ✅         |

### Analyze Command

Analyze previously downloaded game data:

```bash
# Basic analysis
./gamedl analyze --competition nfl --analysis action-types --seasons 2024

# NCAAB analysis with custom directories
./gamedl analyze --competition ncaab --analysis review-types --input-dir ./my_data --output ./analysis_results

# Using environment variables
GAMEDL_ANALYZE_COMPETITION=nfl GAMEDL_ANALYZE_ANALYSIS=action-types ./gamedl analyze --seasons 2024
```

#### Analyze Options

- `--competition, -c`: Competition to analyze (values allowed: 'nfl', 'ncaab' or ncaaf) **(required)**
- `--analysis, -a`: Analysis type to perform (values allowed: 'action-types' or 'review-types') **(required)**
- `--input-dir, -i`: Directory containing downloaded game files (default: "downloaded_games")
- `--output, -o`: Output directory for analysis results (default: "analysis_results")
- `--seasons, -s`: Seasons to include in analysis, comma-separated. e.g '2023,2024' (default: all seasons available)

#### Available Analysis Types

| Competition | Analysis Name | Description                                         |
|-------------|---------------|-----------------------------------------------------|
| NFL         | action-types  | Analyzes play-by-play action types and sequences    |
| NCAAB       | review-types  | Analyzes challenge reviews and related events       |
| NCAAF       | review-types  | Analyzes overturned play reviews and related events |

## Configuration

If there's a value for some flag that you pass often, it might make sense to set it via an environment variable or in a configuration file so you don't have to repeat it every time.

This CLI tool follows the 12-factor app principles for configuration management, where particular config values can be set via command-line flags, environment variables, and/or configuration files.
Flags take precedence over environment variables, which in turn take precedence over configuration files.

### Configuration Options

The following table shows all configuration options and how they can be set:

#### Global Options

| Config Key | Environment Variable | CLI Flag   | Description                                                                                   |
|------------|----------------------|------------|-----------------------------------------------------------------------------------------------|
| N/A        | N/A                  | `--config` | Config file to use (default `.gamedl.yaml` in the current directory or in the home directory) |

#### Download Command Options

| Config Key             | Environment Variable          | CLI Flag              | Description                                   |
|------------------------|-------------------------------|-----------------------|-----------------------------------------------|
| `download.competition` | `GAMEDL_DOWNLOAD_COMPETITION` | `--competition, -c`   | Competition to download (nfl, ncaab, ncaaf)   |
| `download.provider`    | `GAMEDL_DOWNLOAD_PROVIDER`    | `--provider, -p`      | Data provider (sportradar, betgenius)         |
| `download.seasons`     | `GAMEDL_DOWNLOAD_SEASONS`     | `--seasons, -s`       | Seasons to download (comma-separated)         |
| `download.output-dir`  | `GAMEDL_DOWNLOAD_OUTPUT_DIR`  | `--output-dir, -o`    | Directory to store downloaded game files      |
| `download.concurrency` | `GAMEDL_DOWNLOAD_CONCURRENCY` | `--concurrency`       | Number of concurrent downloads                |

#### Analyze Command Options

| Config Key            | Environment Variable          | CLI Flag            | Description                                    |
|-----------------------|-------------------------------|---------------------|------------------------------------------------|
| `analyze.competition` | `GAMEDL_ANALYZE_COMPETITION`  | `--competition, -c` | Competition to analyze (nfl, ncaab, ncaaf)     |
| `analyze.analysis`    | `GAMEDL_ANALYZE_ANALYSIS`     | `--analysis, -a`    | Analysis name to perform                       |
| `analyze.input-dir`   | `GAMEDL_ANALYZE_INPUT_DIR`    | `--input-dir, -i`   | Directory containing downloaded game files     |
| `analyze.output`      | `GAMEDL_ANALYZE_OUTPUT`       | `--output, -o`      | Output directory for analysis results          |
| `analyze.seasons`       | `GAMEDL_ANALYZE_SEASONS`        | `--seasons, -s`     | Seasons to include in analysis (comma-separated) |

### Configuration File

The config keys in the table above refer to the options that can be set in a YAML configuration file.
By default, the tool looks for a file named `.gamedl.yaml` in the current working directory and in the user's home directory.
You can specify a different config file using the `--config, -f` flag.

Example `~/.gamedl.yaml`:

```yaml
# Download defaults
download:
  competition: nfl
  provider: betgenius
  seasons: [2023, 2024]
  output-dir: "downloaded_games"
  concurrency: 10

# Analyze defaults  
analyze:
  competition: nfl
  analysis: action-types
  input-dir: "downloaded_games"
  output: "analysis_results"
  seasons: [2021, 2022, 2023, 2024]
```

## Output

### Downloaded Data

Game data is stored in directories organized by competition and year:

```
downloaded_games/
├── nfl/
│   ├── 2023/
│   │   ├── game1.json
│   │   └── game2.json
│   └── 2024/
│       ├── game1.json
│       └── game2.json
└── ncaab/
    ├── 2023/
    │   ├── game1.json
    │   └── game2.json
    └── 2024/
        ├── game1.json
        └── game2.json
```

### Analysis Results

Analysis results are saved as JSON files in the specified output directory:

#### NFL Analysis
- `actions_to_games.json`: Action types mapped to games
- `sub_actions_to_games.json`: Sub-action types mapped to games
- `action_type_count.json`: Count of each action type
- `sub_action_type_count.json`: Count of each sub-action type

#### NCAAB Analysis
- `review_events_to_games.json`: Review events mapped to games
- `event_type_count.json`: Count of each event type
- `review_games_ncaab/`: Sample game files for review events

#### NCAAF Analysis
- `types_to_games.json`: Review types mapped to games
- `review_type_count.json`: Count of each review type
- `review_games/`: Sample game files for review types

## Examples

### Complete Workflow

1. **Download NFL data from BetGenius:**
   ```bash
   ./gamedl download --competition nfl --provider betgenius --seasons 2024
   ```

2. **Analyze the downloaded data:**
   ```bash
   ./gamedl analyze --competition nfl --analysis action-types --seasons 2024 --output ./nfl_analysis
   ```

3. **Download NCAAB data from SportRadar:**
   ```bash
   ./gamedl download --competition ncaab --provider sr --seasons 2024
   ```

4. **Analyze NCAAB reviews:**
   ```bash
   ./gamedl analyze --competition ncaab --analysis review-types --seasons 2024
   ```

### Using Environment Variables

```bash
# Set up environment
export GAMEDL_DOWNLOAD_COMPETITION=nfl
export GAMEDL_DOWNLOAD_PROVIDER=betgenius
export GAMEDL_DOWNLOAD_CONCURRENCY=4

# Download with environment defaults
./gamedl download --seasons 2024

# Override specific values
./gamedl download --competition ncaab --provider sr --seasons 2023,2024
```

## Help

Get help for any command:

```bash
./gamedl --help                    # General help
./gamedl download --help           # Download command help
./gamedl analyze --help            # Analyze command help
```
