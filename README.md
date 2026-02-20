# fyne-components

Production-ready Go components for the Fyne UI toolkit

## Overview

A Go module library providing high-quality, reusable components for [Fyne](https://fyne.io/) applications. Built for production use with comprehensive testing, documentation, and examples.

### Components

- **[Table Widget](pkg/table/)** - Full-featured sortable, filterable table with inline editing, hierarchical data, keyboard navigation, and custom rendering
- **[Password Component](pkg/password/)** - Password strength validation and visual feedback with entropy-based scoring

## Installation

```bash
go get github.com/digital-michael/fyne-components
```

## Quick Start

### Table Widget

```go
import (
    "fyne.io/fyne/v2/app"
    "github.com/digital-michael/fyne-components/pkg/table"
)

func main() {
    myApp := app.New()
    window := myApp.NewWindow("Table Demo")

    // Configure table
    config := table.NewConfig("my-table")
    config.Columns = []table.ColumnConfig{
        {ID: "id", Title: "ID", Width: 60, Sortable: true},
        {ID: "name", Title: "Name", Width: 200, Editable: true},
        {ID: "status", Title: "Status", Width: 100, Alignment: table.AlignCenter},
    }

    // Create and populate table
    tableWidget := table.NewTable(config)
    tableWidget.SetData([]interface{}{
        MyData{ID: 1, Name: "Alice", Status: "Active"},
        MyData{ID: 2, Name: "Bob", Status: "Inactive"},
    })

    window.SetContent(tableWidget)
    window.ShowAndRun()
}
```

**Features**: Sorting, filtering, inline editing, keyboard navigation, hierarchical data, custom rendering

[📖 Table Documentation](pkg/table/) | [🎯 Table Demo](examples/table-demo/)

### Password Component

```go
import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/digital-michael/fyne-components/pkg/password"
)

func main() {
    myApp := app.New()
    window := myApp.NewWindow("Password Demo")

    // Create password entry with strength meter
    passwordEntry := widget.NewPasswordEntry()
    meter := password.NewStrengthMeter()

    passwordEntry.OnChanged = func(text string) {
        strength := password.CalculateStrength(text)
        meter.SetStrength(strength)
        
        // Validate: require score >= 50
        if strength.Score < 50 {
            // Show error or disable submit button
        }
    }

    content := container.NewVBox(
        widget.NewLabel("Password:"),
        passwordEntry,
        meter,
    )

    window.SetContent(content)
    window.ShowAndRun()
}
```

**Features**: Strength calculation (0-100), visual meter, entropy-based scoring, pattern detection, feedback

[📖 Password Documentation](pkg/password/) | [🎯 Password Demo](examples/password-demo/)

## Components Overview

### Table Widget [![Tests](https://img.shields.io/badge/tests-18%20passing-success)](pkg/table/)

A comprehensive table widget with enterprise features:

| Feature | Description |
|---------|-------------|
| **Sorting** | Click headers to sort ascending/descending with custom comparators |
| **Filtering** | Built-in search box with regex support and case sensitivity |
| **Editing** | Inline cell editing with Enter/Escape shortcuts |
| **Navigation** | Full keyboard support (arrows, Tab, Page Up/Down, Home/End) |
| **Hierarchy** | Tree-like structures with indentation and visual indicators |
| **Selection** | Single or multi-row selection with callbacks |
| **Customization** | Per-column renderers, alignments, and interactive elements |
| **Checkboxes** | Built-in checkbox cells with state management |
| **Popup Menus** | Dropdown menus for cell values |
| **Resizing** | Manual drag and double-click auto-resize |

**Use Cases**: Data tables, file browsers, task lists, project trees, configuration editors

### Password Component [![Tests](https://img.shields.io/badge/tests-63%20passing-success)](pkg/password/)

Password strength validation with visual feedback:

| Feature | Description |
|---------|-------------|
| **Strength Levels** | 5 levels from Very Weak (0-20) to Very Strong (81-100) |
| **Entropy Calculation** | Real entropy-based scoring for accurate strength measurement |
| **Character Variety** | Checks lowercase, uppercase, numbers, special characters |
| **Pattern Detection** | Identifies sequential chars, repetition, common passwords |
| **Visual Meter** | Color-coded progress bar with real-time updates |
| **Feedback** | Actionable suggestions for improvement |
| **Lightweight** | Zero dependencies beyond Fyne |

**Use Cases**: Registration forms, password change dialogs, security settings, user account management

## Examples

Each component includes a working demonstration:

```bash
# Table demo - interactive table with hierarchical tasks
cd examples/table-demo
go run main.go

# Password demo - real-time strength validation
cd examples/password-demo
go run main.go
```

## Testing

All components include comprehensive test coverage:

```bash
# Run all tests
go test ./pkg/...

# Run tests with coverage
go test ./pkg/... -cover

# Generate coverage report
make coverage
```

**Test Stats:**
- Table Widget: 18 tests covering sorting, state, data operations
- Password Component: 63 tests covering all strength levels and edge cases
- Total: 81 tests passing

## Development

### Project Structure

```
fyne-components/
├── pkg/               # Public library code (importable)
│   ├── table/        # Table widget (4,067 lines, 18 tests)
│   └── password/     # Password component (736 lines, 63 tests)
├── examples/          # Demo applications
│   ├── table-demo/   # Interactive table demonstration
│   └── password-demo/# Password strength demonstration
├── internal/          # Private code (not importable)
├── doc/              # Documentation
├── .github/          # CI/CD workflows
│   └── workflows/
│       └── ci.yml    # Automated testing and linting
├── go.mod            # Module definition
├── go.work           # Workspace configuration
├── Makefile          # Development automation
└── README.md         # This file
```

### Prerequisites

- Go 1.24 or later
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)
- Fyne v2.7.1 or compatible

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Workflow

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backwards-compatible)
- **PATCH**: Bug fixes (backwards-compatible)

Current: Development (targeting v0.1.0-beta)

## License

See [LICENSE](LICENSE) file for details.

## Related Projects

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit for Go
- [Resource Tech Kit](https://github.com/digital-michael/resource-tech-kit) - Original project these components were extracted from

## Support

- 📖 [Component Documentation](pkg/)
- 🐛 [Report Issues](https://github.com/digital-michael/fyne-components/issues)
- 💬 [Discussions](https://github.com/digital-michael/fyne-components/discussions)
