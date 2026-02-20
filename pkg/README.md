# pkg/

Public library code for fyne-components - importable by external projects.

## Available Components

### [Table Widget](table/)

Full-featured sortable, filterable table with inline editing and hierarchical data support.

```go
import "github.com/digital-michael/fyne-components/pkg/table"
```

**Features**: Sorting, filtering, inline editing, keyboard navigation, tree structures, custom rendering

[📖 Documentation](table/) | [🎯 Demo](../examples/table-demo/)

### [Password Component](password/)

Password strength validation and visual feedback with entropy-based scoring.

```go
import "github.com/digital-michael/fyne-components/pkg/password"
```

**Features**: Strength calculation (0-100), visual meter, pattern detection, feedback suggestions

[📖 Documentation](password/) | [🎯 Demo](../examples/password-demo/)

## Usage

Each component has its own README with comprehensive documentation, API reference, and examples.

## Testing

All components include comprehensive test suites:

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -cover
```

**Test Coverage:**
- Table: 18 tests
- Password: 63 tests
- Total: 81 tests passing
