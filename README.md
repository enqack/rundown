# Rundown

A beautiful, fast terminal system monitor built with Go and Bubble Tea.

![License](https://img.shields.io/badge/license-BSD--3--Clause-blue.svg)

## Features

- **Real-time System Monitoring**: CPU, memory, disk, network, and temperature metrics
- **7 Specialized Views**:
  - **Overview**: System summary with key metrics
  - **CPU**: Per-core utilization and load averages
  - **Memory**: RAM and swap usage with detailed breakdowns
  - **Disk**: Partition usage and device information
  - **Network**: Interface statistics with ingress/egress graphs
  - **Temperature**: Sensor readings for CPU, GPU, and system components
  - **Processes**: htop-style process monitor with sorting
- **4 Beautiful Themes**:
  - `base16`: Adapts to your terminal's color scheme
  - `cyberpunk`: Vibrant purple/green/pink palette
  - `monochrome`: Clean grayscale design
  - `phosphor`: Retro green CRT monitor aesthetic
- **Smooth Scrolling**: Navigate large datasets with viewport controls
- **Configurable Updates**: Adjust refresh rate with `+/-` keys
- **Keyboard-Driven**: Efficient navigation without touching the mouse

## Installation

### Using Nix

```bash
# Build and run
nix build
./result/bin/rundown

# Or install to your profile
nix profile install .#default
```

### Using Go

```bash
# Clone the repository
git clone https://github.com/enqack/rundown.git
cd rundown

# Build
go build -o rundown ./cmd/rundown

# Or install
go install ./cmd/rundown
```

### Using Mage

```bash
# Install mage if needed
go install github.com/magefile/mage@latest

# Build
mage build

# Install to $GOPATH/bin
mage install
```

## Usage

```bash
# Run with default theme (base16)
rundown

# Run with a specific theme
rundown -t cyberpunk
rundown -t monochrome
rundown -t phosphor
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `Tab` | Switch between tabs |
| `1-7` | Jump to specific tab |
| `j/k` | Scroll down/up (when content is scrollable) |
| `c/m/p/n/u/t` | Sort processes by CPU/Memory/PID/Name/User/Time |
| `l/f/p/n/s` | Sort network by Local/Foreign/Proto/Name/State |
| `+/-` | Increase/decrease update interval |
| `q` | Quit |

## Configuration

Rundown supports configuration via YAML file:

```yaml
# ~/.config/rundown/rundown.yaml
theme: cyberpunk
```

Command-line flags take precedence over config file settings.

## Development

### Prerequisites

- Go 1.22 or later
- Nix (optional, for reproducible builds)
- Mage (optional, for build automation)

### Building

```bash
# Standard Go build
CGO_ENABLED=0 go build -o rundown ./cmd/rundown

# Using Mage
mage build

# Using Nix
nix develop --command mage build
```

### Testing

```bash
# Run tests
go test -v ./...

# Using Mage
mage test
```

### Code Quality

```bash
# Run full quality pipeline (format, lint, vet, test, build)
mage all

# Individual tasks
mage fmt    # Format code
mage lint   # Run golangci-lint
mage vet    # Run go vet
```

## Architecture

- **`cmd/rundown`**: Main entry point
- **`internal/cmd`**: Cobra command setup
- **`internal/ui`**: Bubble Tea UI components and views
- **`internal/stats`**: System metrics collection
- **`internal/theme`**: Color themes and styling
- **`internal/layout`**: Layout calculations and helpers

## License

BSD-3-Clause License - see [LICENSE](LICENSE) for details.

## Author

**enqack**

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [gopsutil](https://github.com/shirou/gopsutil) - System metrics
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
