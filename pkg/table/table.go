package table

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Table is a feature-rich table widget for displaying and editing tabular data.
//
// The Table widget provides extensive functionality including:
//   - Multi-column sorting
//   - Cell selection and navigation
//   - In-place editing
//   - Flexible column configuration
//   - Data binding and refresh
//   - Keyboard shortcuts
type Table struct {
	widget.BaseWidget

	config *Config
	logger Logger

	// State management (implementation in internal/table)
	// These fields will be populated during Phase 2 migration
}

// NewTable creates a new Table widget with the specified configuration.
//
// Example:
//
//	table := table.NewTable(&table.Config{
//		ID: "users",
//		Columns: []table.ColumnConfig{
//			{ID: "name", Label: "Name", Width: 200},
//			{ID: "email", Label: "Email", Width: 300},
//		},
//		Data: userData,
//		OnRowSelected: handleSelection,
//	})
func NewTable(config *Config) *Table {
	if config == nil {
		panic("table.NewTable: config cannot be nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("table.NewTable: invalid config: %v", err))
	}

	// Set defaults
	if config.Logger == nil {
		config.Logger = NoopLogger{}
	}

	t := &Table{
		config: config,
		logger: config.Logger,
	}

	t.ExtendBaseWidget(t)
	return t
}

// CreateRenderer creates the renderer for this table widget.
// This will be implemented in Phase 2 when we migrate the actual table code.
func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	// Placeholder - will be implemented in Phase 2
	t.logger.Debug("CreateRenderer called", "table", t.config.ID)
	return nil
}

// GetSelectedCell returns the currently selected cell coordinates.
// Returns (-1, -1) if no cell is selected.
func (t *Table) GetSelectedCell() (row int, col int) {
	// Implementation in Phase 2
	return -1, -1
}

// SetSelectedCell programmatically selects a cell.
func (t *Table) SetSelectedCell(row, col int) {
	// Implementation in Phase 2
	t.logger.Debug("SetSelectedCell", "row", row, "col", col)
}

// GetData returns the current table data.
func (t *Table) GetData() []interface{} {
	return t.config.Data
}

// SetData updates the table data and refreshes the display.
func (t *Table) SetData(data []interface{}) {
	t.config.Data = data
	t.logger.Info("SetData", "rows", len(data))
	t.Refresh()
}

// ClearSelection removes any cell or row selection.
func (t *Table) ClearSelection() {
	// Implementation in Phase 2
	t.logger.Debug("ClearSelection", "table", t.config.ID)
}

// Refresh triggers a visual refresh of the table.
func (t *Table) Refresh() {
	// Implementation in Phase 2
	t.logger.Debug("Refresh", "table", t.config.ID)
	t.BaseWidget.Refresh()
}

// GetConfig returns the table's configuration.
// Note: Modifying the returned config after table creation may have
// unexpected results. Use setter methods where available.
func (t *Table) GetConfig() *Config {
	return t.config
}
