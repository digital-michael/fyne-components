# fyne-components

Go lang components for the Fyne UI toolkit

## Overview

This is a Go module library providing reusable components for [Fyne](https://fyne.io/) applications. Import this package to use custom widgets, layouts, and themes in your Fyne projects.

## Installation

```bash
go get github.com/digital-michael/fyne-components
```

## Project Structure

```
fyne-components/
├── pkg/          # Public library code (importable by external projects)
├── internal/     # Private code (not importable externally)
├── cmd/          # Command-line applications and tools
├── examples/     # Example code demonstrating library usage
├── doc/          # Additional documentation
├── go.mod        # Go module definition
├── Makefile      # Development automation
└── README.md     # This file
```

## Development

### Prerequisites

- Go 1.23 or later
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

### Makefile Targets

Run `make help` to see all available targets:

```bash
make help          # Show all available targets
make all           # Run deps, fmt, vet, lint, and test
make deps          # Download and tidy dependencies
make fmt           # Format code
make vet           # Run go vet
make lint          # Run golangci-lint
make test          # Run tests
make coverage      # Generate HTML coverage report
make clean         # Clean build artifacts
make ci            # Run full CI pipeline
```

### Quick Start for Development

```bash
# Install dependencies
make deps

# Run all checks and tests
make all

# Run tests with coverage
make coverage
```

## Usage

```go
import "github.com/digital-michael/fyne-components/pkg/your-package"
```

## License

See [LICENSE](LICENSE) file for details.
