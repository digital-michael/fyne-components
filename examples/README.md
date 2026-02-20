# examples/

Interactive demonstrations of fyne-components library features.

## Available Demos

### [Table Demo](table-demo/)

Interactive demonstration of the table widget with hierarchical task data.

**Features Shown:**
- Sortable columns (ID, Name, Status, Priority, External)
- String and numeric comparators
- Text alignment options
- Inline editing
- Hierarchical data structure
- Keyboard navigation

**Run:**
```bash
cd table-demo
go run main.go
```

[📖 Demo README](table-demo/README.md) | [📖 Table Docs](../pkg/table/)

### [Password Demo](password-demo/)

Interactive password strength validation with real-time visual feedback.

**Features Shown:**
- Real-time strength calculation
- Visual strength meter with color coding
- All strength levels (Very Weak to Very Strong)
- Detailed strength breakdown
- Feedback suggestions
- Entropy display

**Run:**
```bash
cd password-demo
go run main.go
```

[📖 Demo README](password-demo/README.md) | [📖 Password Docs](../pkg/password/)

## Building Demos

Each demo can be built as a standalone binary:

```bash
# Build table demo
cd table-demo && go build

# Build password demo
cd password-demo && go build
```

## Demo Statistics

- **Table Demo**: 184 lines (main.go + task_data.go)
- **Password Demo**: ~150 lines (main.go)
- **Both demos**: Zero external dependencies beyond Fyne and fyne-components
