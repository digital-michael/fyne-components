# Table Widget Demo

Interactive demonstration of the fyne-components table widget.

## Features Demonstrated

- ✅ Hierarchical task data with 19 sample tasks
- ✅ Sortable columns (click headers to sort)
- ✅ Multiple column types (ID, text, status, priority, boolean)
- ✅ Text alignment per column (left, center, right)
- ✅ Inline editing for Task Name column (double-click to edit)
- ✅ Row selection with callbacks
- ✅ Keyboard navigation (arrow keys, Tab/Shift-Tab)
- ✅ String and numeric comparators for sorting

## Building

```bash
go build
```

## Running

```bash
./table-demo
```

## Usage

### Navigation
- **Arrow Keys**: Navigate between cells
- **Tab**: Move to next cell
- **Shift+Tab**: Move to previous cell
- **Page Up/Down**: Scroll through rows

### Sorting
- **Click Column Header**: Toggle sort (ascending/descending)
- **ID Column**: Numeric sort
- **Task Name**: Alphabetic sort
- **Status**: Alphabetic sort  
- **Priority**: Numeric sort
- **External**: Boolean sort (false then true)

### Editing
- **Double-click Task Name cell**: Enter edit mode
- **Enter**: Save changes
- **Escape**: Cancel editing

## Data Structure

The demo uses a hierarchical `TaskData` struct with:
- **ID**: Unique task identifier
- **Name**: Task name/description
- **Status**: Current status (In Progress, Complete, Not Started, Pending)
- **Priority**: Numeric priority (1-3)
- **Depth**: Hierarchy level (for future tree view features)
- **External**: Boolean flag for external/internal tasks

## Code Example

```go
// Create configuration
config := table.NewConfig("demo-table")
config.Columns = []table.ColumnConfig{
    {
        ID:         "id",
        Title:      "ID",
        Width:      60,
        Sortable:   true,
        Alignment:  table.AlignRight,
        Comparator: table.NewNumericComparator("ID"),
    },
    {
        ID:         "name",
        Title:      "Task Name",
        Width:      300,
        Sortable:   true,
        Editable:   true,
        Comparator: table.NewStringComparator("Name"),
    },
}

// Create table widget
tableWidget := table.NewTable(config)

// Set data
data := []interface{}{task1, task2, task3}
tableWidget.SetData(data)

// Add callbacks
config.OnRowSelected = func(row int, rowData interface{}) {
    // Handle selection
}
```

## Dependencies

- `fyne.io/fyne/v2@v2.7.1`
- `github.com/digital-michael/fyne-components`(local)

## Related

- [Table Widget Source](../../pkg/table/)
- [Table Widget Tests](../../pkg/table/table_test.go)
- [Password Demo](../password-demo/)
