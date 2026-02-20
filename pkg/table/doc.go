// Package table provides a feature-rich, production-proven table widget for Fyne applications.
//
// This table widget was extracted from the Resource Tech Kit project after extensive
// production use and debugging. It provides advanced features including:
//
//   - Multi-column sorting with visual indicators
//   - Cell selection and navigation (keyboard and mouse)
//   - In-place cell editing with validation
//   - Flexible column configuration (width, alignment, visibility)
//   - Data binding with automatic refresh
//   - Keyboard shortcuts (arrow keys, Tab, Enter, Escape)
//   - Custom cell rendering
//   - Row selection callbacks
//
// # Basic Usage
//
//	columns := []table.ColumnConfig{
//		{ID: "name", Label: "Name", Width: 200},
//		{ID: "age", Label: "Age", Width: 100, Alignment: fyne.TextAlignTrailing},
//	}
//
//	data := []interface{}{
//		map[string]interface{}{"name": "Alice", "age": 30},
//		map[string]interface{}{"name": "Bob", "age": 25},
//	}
//
//	config := &table.Config{
//		ID:          "myTable",
//		Columns:     columns,
//		Data:        data,
//		ShowHeaders: true,
//		OnRowSelected: func(row int, data interface{}) {
//			fmt.Printf("Selected row %d\n", row)
//		},
//	}
//
//	tableWidget := table.NewTable(config)
//
// # Configuration
//
// The Config struct provides extensive customization options:
//
//   - Columns: Define column properties (ID, label, width, alignment, sorting)
//   - Data: Slice of data items (maps, structs, or any interface{})
//   - Callbacks: OnRowSelected, OnCellEdited, OnKeyPressed for event handling
//   - Styling: Header colors, selection colors, custom renderers
//   - Logging: Optional Logger interface for debugging
//
// # Keyboard Navigation
//
//   - Arrow Keys: Move cell selection
//   - Tab/Shift+Tab: Move to next/previous cell
//   - Enter: Start editing selected cell (if editable)
//   - Escape: Cancel editing
//   - Page Up/Down: Scroll by page
//   - Home/End: Jump to first/last row
//
// # State Management
//
// All table state is accessible via public methods:
//
//   - GetSelectedCell() - Query current selection
//   - SetSelectedCell(row, col) - Programmatically select cell
//   - GetData() - Retrieve current data
//   - SetData(data) - Update table data
//   - Refresh() - Force visual refresh
//
// # Thread Safety
//
// The table widget is NOT thread-safe. All methods must be called from the
// UI goroutine. Use fyne.CurrentApp().Driver().DoEventSync() when updating
// from background goroutines.
package table
