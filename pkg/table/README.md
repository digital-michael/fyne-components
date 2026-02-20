# Table Widget

A full-featured, sortable, filterable table widget for Fyne applications with support for hierarchical data, inline editing, keyboard navigation, and custom rendering.

## Features

- ✅ **Sorting**: Click column headers to sort (ascending/descending)
- ✅ **Filtering**: Built-in search/filter with regex support
- ✅ **Inline Editing**: Double-click cells or press Space/Enter to edit
- ✅ **Keyboard Navigation**: Arrow keys, Tab, Page Up/Down, Home/End
- ✅ **Hierarchical Data**: Tree-like structure with indentation and visual indicators
- ✅ **Custom Rendering**: Per-column custom cell renderers
- ✅ **Flexible Alignment**: Left, center, right text alignment per column
- ✅ **Row Selection**: Single or multi-select with callbacks
- ✅ **Column Resizing**: Manual drag resize and double-click auto-resize
- ✅ **Interactive Cells**: Checkboxes, popup menus, custom actions
- ✅ **Theming**: Multiple tree icon themes and custom styling

## Installation

```bash
go get github.com/digital-michael/fyne-components
```

## Quick Start

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "github.com/digital-michael/fyne-components/pkg/table"
)

func main() {
    myApp := app.New()
    window := myApp.NewWindow("Table Demo")

    // Create configuration
    config := table.NewConfig("my-table")
    
    // Define columns
    config.Columns = []table.ColumnConfig{
        {
            ID:       "id",
            Title:    "ID",
            Width:    60,
            Sortable: true,
            Alignment: table.AlignRight,
        },
        {
            ID:       "name",
            Title:    "Name",
            Width:    200,
            Sortable: true,
            Editable: true,
        },
        {
            ID:       "status",
            Title:    "Status",
            Width:    100,
            Sortable: true,
            Alignment: table.AlignCenter,
        },
    }
    
    // Create table widget
    tableWidget := table.NewTable(config)
    
    // Set data
    data := []interface{}{
        MyData{ID: 1, Name: "Alice", Status: "Active"},
        MyData{ID: 2, Name: "Bob", Status: "Inactive"},
    }
    tableWidget.SetData(data)
    
    window.SetContent(tableWidget)
    window.ShowAndRun()
}

type MyData struct {
    ID     int
    Name   string
    Status string
}
```

## Configuration

### Table Configuration

Create a config using `table.NewConfig(id string)`, then customize:

```go
config := table.NewConfig("my-table")

// Core settings
config.RowHeight = 35.0
config.HeaderHeight = 30.0

// Features
config.ShowSearch = true              // Enable search/filter box
config.SearchPlaceholder = "Search..."
config.ShowHeaders = true             // Show column headers
config.AllowMultiSelect = false       // Single or multi-select

// Tree hierarchy
config.ShowIndentation = true         // Enable indentation
config.IndentPerLevel = 20.0          // Pixels per level
config.ShowIndentIcons = true         // Show tree icons (├ └)
config.TreeIconTheme = table.TreeThemeBranches

// Keyboard behavior
config.RowSelectOnlyMode = true       // true = arrow keys select rows only
config.SelectFirstCellOnStartup = true // Auto-select first cell

// Visual styling
config.RootNodeBackgroundColor = color.NRGBA{R: 35, G: 35, B: 65, A: 255}
config.FontSize = 12.0
config.FontFamily = ""                // Empty = system default

// Callbacks
config.OnRowSelected = func(rowIndex int, data interface{}) {
    fmt.Printf("Selected row %d\n", rowIndex)
}

config.OnCellEdited = func(rowIndex int, colID string, newValue string, data interface{}) {
    fmt.Printf("Edited: row=%d, col=%s, value=%s\n", rowIndex, colID, newValue)
}
```

### Column Configuration

Each column is defined with a `ColumnConfig`:

```go
type ColumnConfig struct {
    ID       string          // Unique column identifier
    Title    string          // Header text
    Width    float32         // Column width in pixels
    MinWidth float32         // Minimum width for resizing

    // Behavior flags
    Sortable bool           // Enable sorting for this column
    Editable bool           // Enable inline editing
    ReadOnly bool           // Prevent keyboard activation
    Hidden   bool           // Hide column by default

    // Visual styling
    Alignment TextAlignment // Left, Center, or Right

    // Custom logic
    Renderer   CellRenderer   // Custom cell renderer
    Comparator SortComparator // Custom sort function
}
```

#### Text Alignment

```go
table.AlignLeft    // Default
table.AlignCenter
table.AlignRight
```

#### Custom Comparators

For proper sorting of different data types:

```go
// String field sorting
Comparator: table.NewStringComparator("FieldName")

// Numeric field sorting
Comparator: table.NewNumericComparator("FieldName")

// Custom sorting logic
Comparator: func(a, b interface{}) int {
    itemA := a.(MyType)
    itemB := b.(MyType)
    if itemA.Value < itemB.Value {
        return -1
    } else if itemA.Value > itemB.Value {
        return 1
    }
    return 0
}
```

## Interactive Features

### Inline Editing

Enable editing for specific columns:

```go
{
    ID:       "name",
    Title:    "Name",
    Editable: true,  // Allow editing
}
```

Users can:
- **Double-click** a cell to start editing
- **Press Space or Enter** on selected cell to edit
- **Press Enter** to save changes
- **Press Escape** to cancel

Handle edits in the callback:

```go
config.OnCellEdited = func(rowIndex int, colID string, newValue string, data interface{}) {
    item := data.(MyType)
    switch colID {
    case "name":
        item.Name = newValue
        // Update your data source here
    }
}
```

### Checkboxes

Create interactive checkbox cells:

```go
{
    ID:    "completed",
    Title: "Done",
    ShowCheckbox: true,
    CheckboxLabel: "Yes|No",
    GetCheckboxValue: func(data interface{}) bool {
        item := data.(MyType)
        return item.Completed
    },
    OnCheckboxChanged: func(data interface{}, checked bool, rowIndex int) {
        // Handle checkbox toggle
    },
}
```

### Popup Menus

Add dropdown menus to cells:

```go
{
    ID:    "status",
    Title: "Status",
    ShowPopupIcon: true,
    PopupOptions: func(data interface{}) []string {
        return []string{"Active", "Inactive", "Pending"}
    },
    OnPopupSelected: func(data interface{}, option string, rowIndex int) string {
        // Handle selection, return new display value
        return option
    },
}
```

### Custom Cell Rendering

For complete control over cell appearance:

```go
{
    ID:    "custom",
    Title: "Custom",
    Renderer: func(data interface{}, container fyne.CanvasObject, rowIndex int, colID string) {
        item := data.(MyType)
        
        // Create custom content
        label := container.(*fyne.Container).Objects[0].(*widget.Label)
        label.SetText(fmt.Sprintf("🎨 %s", item.Value))
    },
}
```

## Hierarchical Data

Display tree-like structures with automatic indentation:

```go
type Task struct {
    ID     int
    Name   string
    Depth  int  // 0 = root, 1 = child, 2 = grandchild
}

// Configure tree functions
config.GetNodeDepth = func(data interface{}) int {
    return data.(Task).Depth
}

config.GetNodeID = func(data interface{}) interface{} {
    return data.(Task).ID
}

config.IsNodeExpandable = func(data interface{}) bool {
    // Return true if node has children
    return checkHasChildren(data.(Task))
}

// Set data
tableWidget.SetData([]interface{}{
    Task{ID: 1, Name: "Project Alpha", Depth: 0},
    Task{ID: 2, Name: "Planning", Depth: 1},
    Task{ID: 3, Name: "Requirements", Depth: 2},
})
```

### Tree Icon Themes

Multiple visual styles available:

```go
table.TreeThemeBranches  // ├ └ (default)
table.TreeThemeArrows    // ► ▶
table.TreeThemeCircles   // ○ ● 
table.TreeThemeDiamonds  // ◇ ◆
table.TreeThemeSquares   // ■ □
table.TreeThemeBullets   // • • •
table.TreeThemeAngles    // > > >

config.TreeIconTheme = table.TreeThemeCircles
```

## Keyboard Navigation

### Standard Navigation

- **Arrow Keys**: Move selection (Up, Down, Left, Right)
- **Tab**: Move to next cell
- **Shift+Tab**: Move to previous cell
- **Home**: First row
- **End**: Last row
- **Page Up/Down**: Scroll by page

### Editing

- **Space**: Edit selected cell (if editable)
- **Enter**: Edit cell or save changes
- **Escape**: Cancel editing
- **Double-click**: Start editing

### Selection

- **Click**: Select single row/cell
- **Ctrl+Click**: Multi-select (if enabled)
- **Shift+Click**: Range select (if enabled)

## Search and Filtering

Enable search box:

```go
config.ShowSearch = true
config.FilterColumns = []string{"name", "status"}  // Which columns to search
```

Users can:
- Type in search box to filter rows
- Check "Regex" for pattern matching
- Check "Match case" for case-sensitive search
- Click "Clear" to remove filter

Programmatic filtering:

```go
tableWidget.SetFilter("search text", false)  // text, isRegex
tableWidget.ClearFilter()
```

## API Reference

### Creating Tables

```go
func NewConfig(id string) *Config
func NewTable(config *Config) *Table
```

### Data Management

```go
func (t *Table) SetData(data []interface{})
func (t *Table) GetData() []interface{}
func (t *Table) Refresh()
```

### Filtering

```go
func (t *Table) SetFilter(text string, isRegex bool)
func (t *Table) SetFilterCaseSensitive(sensitive bool)
func (t *Table) ClearFilter()
```

### Selection

```go
func (t *Table) GetSelectedRow() int
func (t *Table) SetSelectedRow(row int)
func (t *Table) ClearSelection()
```

### Focus

```go
func (t *Table) RequestFocus()
func (t *Table) FocusGained()
func (t *Table) FocusLost()
```

## Examples

See the [table-demo](../../examples/table-demo/) for a complete working example with:
- Multiple column types (ID, text, status, priority)
- Sortable columns with custom comparators
- Inline editing
- Hierarchical task data
- Full keyboard navigation

## Testing

The table widget includes comprehensive tests:

```bash
cd pkg/table
go test -v
```

Test coverage includes:
- Sorting (string, numeric, custom comparators)
- State management (selection, editing, filtering)
- Data operations
- Keyboard handling
- Cell rendering

## Performance Considerations

- **Large Datasets**: The widget uses Fyne's native table which efficiently handles large datasets via virtual scrolling
- **Custom Renderers**: Keep render functions lightweight to maintain smooth scrolling
- **Filtering**: Regex filtering on very large datasets may impact performance; use plain text search when possible

## Migration from RTK

If migrating from `resource-tech-kit/internal/ui/widgets`:

| RTK | fyne-components |
|-----|-----------------|
| `widgets.Table` | `table.Table` |
| `widgets.TableConfig` | `table.Config` |
| `widgets.NewTableConfig()` | `table.NewConfig()` |
| `widgets.NewTable()` | `table.NewTable()` |
| `Label` field | `Title` field |
| `fyne.TextAlign*` | `table.Align*` |

## License

See the [LICENSE](../../LICENSE) file for details.

## Related

- [Table Demo](../../examples/table-demo/) - Interactive demonstration
- [Password Component](../password/) - Password strength validation
- [Design Principles](../../DESIGN-PRINCIPLES.md) - Architecture guidelines
