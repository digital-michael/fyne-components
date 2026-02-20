package table

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// MouseHandler defines the interface for handling mouse/touch events in a table
type MouseHandler interface {
	// HandleSingleTap processes single tap/click events
	HandleSingleTap(ev *fyne.PointEvent, table *Table)
	// HandleDoubleTap processes double tap/click events
	HandleDoubleTap(ev *fyne.PointEvent, table *Table)
	// HandleCellClick processes cell selection events
	HandleCellClick(id widget.TableCellID, table *Table)
	// HandleHeaderClick processes column header click events
	HandleHeaderClick(colIndex int, table *Table)
}

// DefaultMouseHandler provides the standard mouse behavior for tables
type DefaultMouseHandler struct{}

// NewDefaultMouseHandler creates a new default mouse handler
func NewDefaultMouseHandler() *DefaultMouseHandler {
	return &DefaultMouseHandler{}
}

// HandleSingleTap processes single tap/click events
func (h *DefaultMouseHandler) HandleSingleTap(ev *fyne.PointEvent, table *Table) {

	// Request focus when table is tapped
	table.RequestFocus()
}

// HandleDoubleTap processes double tap/click events
func (h *DefaultMouseHandler) HandleDoubleTap(ev *fyne.PointEvent, table *Table) {

	// Check if the double-tap is near a column divider in the header row
	// for column auto-resize functionality
	if table.config.EnableDoubleClickResize {
		colIndex := table.findColumnDividerAtPosition(ev.Position)
		if colIndex >= 0 && colIndex < len(table.config.Columns) {
			table.autoResizeColumn(colIndex)
		}
	}
}

// HandleCellClick processes cell selection events
func (h *DefaultMouseHandler) HandleCellClick(id widget.TableCellID, table *Table) {
	// Header row (row 0) - handle column header clicks
	if id.Row == 0 {
		h.HandleHeaderClick(id.Col, table)
		return
	}

	// Data row click - update selection
	displayRowIndex := id.Row - 1
	if displayRowIndex < 0 || displayRowIndex >= len(table.state.visibleRows) {
		return
	}

	dataIndex := table.state.visibleRows[displayRowIndex]
	if dataIndex < 0 || dataIndex >= len(table.data) {
		return
	}

	// Map display column index to actual column index
	if id.Col >= 0 && id.Col < len(table.state.visibleColumns) {
		table.state.selectedCol = table.state.visibleColumns[id.Col]
	} else {
		// Fallback to first visible column
		if len(table.state.visibleColumns) > 0 {
			table.state.selectedCol = table.state.visibleColumns[0]
		}
	}

	// Handle selection based on multi-select mode
	if table.config.AllowMultiSelect {
		// In multi-select mode: toggle row selection
		if table.state.IsRowSelected(dataIndex) {
			table.state.RemoveSelectedRow(dataIndex)
		} else {
			table.state.AddSelectedRow(dataIndex)
		}

		// Fire OnRowSelected for each selected row (skip during programmatic re-selection)
		if !table.state.isReselecting && table.config.OnRowSelected != nil {
			selectedRows := table.state.GetSelectedRows()
			for _, rowIndex := range selectedRows {
				if rowIndex >= 0 && rowIndex < len(table.data) {
					table.config.OnRowSelected(rowIndex, table.data[rowIndex])
				}
			}
		}
	} else {
		// Single-select mode: replace selection
		table.logger().Info(fmt.Sprintf("[CLICK] HandleCellClick: id={Row:%d,Col:%d} displayRowIndex=%d dataIndex=%d isReselecting=%v",
			id.Row, id.Col, displayRowIndex, dataIndex, table.state.isReselecting))

		table.state.selectedRow = dataIndex
		table.state.selectedRows = make(map[int]bool) // Clear multi-select

		if table.state.selectedRow >= 0 {
			// Fire OnRowSelected callback if configured (skip during programmatic re-selection)
			if !table.state.isReselecting && table.config.OnRowSelected != nil {
				table.logger().Info(fmt.Sprintf("[CLICK] Calling OnRowSelected for row=%d", dataIndex))
				table.config.OnRowSelected(dataIndex, table.data[dataIndex])
			} else if table.state.isReselecting {
				table.logger().Info(fmt.Sprintf("[CLICK] Skipping OnRowSelected (isReselecting=true) for row=%d", dataIndex))
			}
		}
	}

	// After selecting the cell, check if it's an interactive cell (checkbox or dropdown)
	// and automatically activate it on mouse click (but not during keyboard navigation or re-selection)
	if !table.state.isKeyboardNavigation && !table.state.isReselecting {
		h.activateInteractiveCell(table, dataIndex, table.state.selectedCol)
	}

	// Skip focus request and refresh during programmatic re-selection (already done by caller)
	if !table.state.isReselecting {
		table.RequestFocus()
		table.table.Refresh()
	}
}

// activateInteractiveCell checks if the clicked cell is a checkbox or dropdown
// and automatically activates it (toggle checkbox or show popup menu)
func (h *DefaultMouseHandler) activateInteractiveCell(table *Table, rowIndex, colIndex int) {
	// Validate indices
	if rowIndex < 0 || rowIndex >= len(table.data) {
		return
	}
	if colIndex < 0 || colIndex >= len(table.config.Columns) {
		return
	}

	// Check if column is read-only
	if table.isColumnReadOnly(colIndex) {
		return
	}

	col := table.config.Columns[colIndex]
	dataItem := table.data[rowIndex]

	// Priority 1: Dropdown/Popup menu
	if col.PopupOptions != nil && col.OnPopupSelected != nil {
		options := col.PopupOptions(dataItem)
		if len(options) > 0 {
			table.showPopupMenu(rowIndex, colIndex, options, dataItem)
			return
		}
	}

	// Priority 2: Checkbox toggle
	if col.ShowCheckbox && col.GetCheckboxValue != nil && col.OnCheckboxChanged != nil {
		// Get current checkbox value
		currentValue := col.GetCheckboxValue(dataItem)

		// Toggle the value
		newValue := !currentValue

		// Call the callback to update the data
		col.OnCheckboxChanged(dataItem, newValue, rowIndex)

		// Trigger OnCellEdited callback if defined
		if table.config.OnCellEdited != nil {
			// Convert bool to string for OnCellEdited callback
			newValueStr := "false"
			if newValue {
				newValueStr = "true"
			}
			table.config.OnCellEdited(rowIndex, col.ID, newValueStr, dataItem)
		}

		// CRITICAL: Restore navigation state to the cell where checkbox was clicked
		// The OnCellEdited callback may have triggered state changes
		// But we want arrow keys to continue from where the checkbox click occurred
		table.logger().Info(fmt.Sprintf("[CHECKBOX-MOUSE] Restoring navigation state from (%d,%d) to (%d,%d)",
			table.state.selectedRow, table.state.selectedCol, rowIndex, colIndex))
		table.state.selectedRow = rowIndex
		table.state.selectedCol = colIndex

		// Refresh the entire table to show updated value
		// Note: We refresh the whole table instead of just the cell to avoid
		// unintended side effects from RefreshItem (which can cause selection changes)
		// We DON'T call Select() because that would trigger OnSelected callbacks
		// and potentially change our navigation state. Our selectedRow/selectedCol
		// remain unchanged, so arrow keys will work correctly from the current position.
		if table.table != nil {
			table.logger().Info(fmt.Sprintf("[REFRESH-MOUSE] Refreshing table after checkbox toggle, navigation state: row=%d col=%d", table.state.selectedRow, table.state.selectedCol))
			table.table.Refresh()
			table.logger().Info(fmt.Sprintf("[REFRESH-MOUSE] Table refreshed, navigation state preserved: row=%d col=%d", table.state.selectedRow, table.state.selectedCol))
		}
	}
}

// HandleHeaderClick processes column header click events
func (h *DefaultMouseHandler) HandleHeaderClick(colIndex int, table *Table) {
	// Get actual column index from visible columns
	if colIndex < 0 || colIndex >= len(table.state.visibleColumns) {
		return
	}
	actualColIndex := table.state.visibleColumns[colIndex]
	col := table.config.Columns[actualColIndex]

	// Only allow sorting if column is sortable
	if !col.Sortable {
		return
	}

	// Toggle sort direction if clicking same column, otherwise sort ascending
	if table.state.sortColumn == actualColIndex {
		table.state.sortAsc = !table.state.sortAsc
	} else {
		table.state.sortColumn = actualColIndex
		table.state.sortAsc = true
	}

	table.sortData()
	table.logger().Info("[SORT] Calling table.Refresh after sort")
	table.table.Refresh()
	table.logger().Info("[SORT] table.Refresh completed")
}
