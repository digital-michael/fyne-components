package table

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// KeyHandler defines the interface for handling keyboard events in a table
type KeyHandler interface {
	// HandleKey processes regular key events (arrow keys, Enter, etc.)
	HandleKey(key *fyne.KeyEvent, table *Table)
	// HandleShortcut processes keyboard shortcuts (Ctrl+Enter, etc.)
	HandleShortcut(shortcut fyne.Shortcut, table *Table)
}

// DefaultKeyHandler provides the standard keyboard behavior for tables
type DefaultKeyHandler struct{}

// NewDefaultKeyHandler creates a new default key handler
func NewDefaultKeyHandler() *DefaultKeyHandler {
	return &DefaultKeyHandler{}
}

// HandleKey processes regular key events
func (h *DefaultKeyHandler) HandleKey(key *fyne.KeyEvent, table *Table) {
	// Don't handle keys if we're currently editing (text entry mode)
	if table.state.editingRow >= 0 && table.state.editingCol >= 0 {
		return
	}

	switch key.Name {
	case fyne.KeyUp:
		h.handleArrowKey("up", table)
	case fyne.KeyDown:
		h.handleArrowKey("down", table)
	case fyne.KeyLeft:
		h.handleArrowKey("left", table)
	case fyne.KeyRight:
		h.handleArrowKey("right", table)
	case fyne.KeyPageUp:
		h.handlePageNavigation("pageup", table)
	case fyne.KeyPageDown:
		h.handlePageNavigation("pagedown", table)
	case fyne.KeyHome:
		h.handlePageNavigation("home", table)
	case fyne.KeyEnd:
		h.handlePageNavigation("end", table)
	case fyne.KeySpace:
		// Space key - unified "activate cell" behavior
		if table.state.selectedRow >= 0 && table.state.selectedCol >= 0 {
			// Check what type of cell this is
			if table.state.selectedCol < len(table.config.Columns) {
				col := table.config.Columns[table.state.selectedCol]

				table.logger().Info(fmt.Sprintf("[DEBUG] SPACE key: row=%d col=%d colID=%s Editable=%v ReadOnly=%v",
					table.state.selectedRow, table.state.selectedCol, col.ID, col.Editable, col.ReadOnly))

				// Priority 1: Popup menu
				if col.PopupOptions != nil && !table.isColumnReadOnly(table.state.selectedCol) {
					table.logger().Info("[DEBUG] Showing popup menu")
					h.handleSpaceKeyPopup(table)
					return
				}

				// Priority 2: Checkbox toggle
				if col.ShowCheckbox && !table.isColumnReadOnly(table.state.selectedCol) {
					table.logger().Info("[DEBUG] Toggling checkbox")
					h.handleSpaceKeyCheckbox(table)
					return
				}

				// Priority 3: Inline text editing
				if col.Editable && !table.isColumnReadOnly(table.state.selectedCol) {
					table.logger().Info("[DEBUG] Starting inline edit")
					table.startEdit(table.state.selectedRow, table.state.selectedCol)
					return
				}

				table.logger().Info("[DEBUG] No action for SPACE key")
			}
		}
		// Note: Don't forward to table.table.TypedKey as it would cause infinite recursion
	case fyne.KeyReturn, fyne.KeyEnter:
		// ENTER - show popup menu for popup columns, toggle checkbox, or start inline editing
		if table.state.selectedRow >= 0 && table.state.selectedCol >= 0 && table.state.selectedCol < len(table.config.Columns) {
			col := table.config.Columns[table.state.selectedCol]

			// Priority 1: Popup menu
			if col.PopupOptions != nil && !table.isColumnReadOnly(table.state.selectedCol) {
				h.handleSpaceKeyPopup(table)
				return
			}

			// Priority 2: Checkbox toggle
			if col.ShowCheckbox && !table.isColumnReadOnly(table.state.selectedCol) {
				h.handleSpaceKeyCheckbox(table)
				return
			}

			// Priority 3: Inline text editing (same as SPACE key)
			if col.Editable && !table.isColumnReadOnly(table.state.selectedCol) {
				table.startEdit(table.state.selectedRow, table.state.selectedCol)
				return
			}
		}
	}
}

// HandleShortcut processes keyboard shortcuts
func (h *DefaultKeyHandler) HandleShortcut(shortcut fyne.Shortcut, table *Table) {
	// Don't handle shortcuts if we're currently editing
	if table.state.editingRow >= 0 && table.state.editingCol >= 0 {
		return
	}

	// Don't handle shortcuts if no row is selected
	if table.state.selectedRow < 0 {
		return
	}

	// Check for CTRL+ENTER (or CMD+ENTER on Mac)
	if typed, ok := shortcut.(*desktop.CustomShortcut); ok {
		if (typed.KeyName == fyne.KeyReturn || typed.KeyName == fyne.KeyEnter) &&
			(typed.Modifier&fyne.KeyModifierControl != 0 || typed.Modifier&fyne.KeyModifierSuper != 0) {
			h.handleKeyboardDoubleClickAction(table)
		}
	}
}

// handleArrowKey handles arrow key navigation for row/column selection
func (h *DefaultKeyHandler) handleArrowKey(direction string, table *Table) {
	// Note: We handle all arrow key navigation ourselves below.
	// Don't call table.table.TypedKey as it would cause infinite recursion
	// via the keyboardForwardingTable wrapper.

	// Initialize selection to first row if nothing selected
	if table.state.selectedRow < 0 {
		if len(table.data) > 0 {
			// Select first visible row
			if len(table.state.visibleRows) > 0 {
				table.state.selectedRow = table.state.visibleRows[0]
			} else {
				table.state.selectedRow = 0
			}
			// Always set to first visible column (even in row-only mode for scrolling)
			if len(table.state.visibleColumns) > 0 {
				table.state.selectedCol = table.state.visibleColumns[0]
			}

			// Trigger auto-scroll to selected cell
			if table.table != nil {
				// Set flag to prevent auto-activation during initial keyboard selection
				table.state.isKeyboardNavigation = true

				// First visible row is always displayRow 1 (row 0 is header)
				displayRow := 1
				table.table.Select(widget.TableCellID{Row: displayRow, Col: 0})
				table.table.Refresh()

				// Clear the keyboard navigation flag after selection is complete
				table.state.isKeyboardNavigation = false
			}
		}
		return
	}

	oldRow := table.state.selectedRow
	oldCol := table.state.selectedCol

	switch direction {
	case "up":
		// Move to previous row in visible rows
		if len(table.state.visibleRows) > 0 {
			// Find current position in visibleRows
			currentIdx := -1
			for idx, dataIdx := range table.state.visibleRows {
				if dataIdx == table.state.selectedRow {
					currentIdx = idx
					break
				}
			}
			// Move to previous visible row
			if currentIdx > 0 {
				table.state.selectedRow = table.state.visibleRows[currentIdx-1]
			} else if currentIdx < 0 && len(table.state.visibleRows) > 0 {
				// Not found, default to first visible row
				table.state.selectedRow = table.state.visibleRows[0]
			}
		}
	case "down":
		// Move to next row in visible rows
		if len(table.state.visibleRows) > 0 {
			// Find current position in visibleRows
			currentIdx := -1
			for idx, dataIdx := range table.state.visibleRows {
				if dataIdx == table.state.selectedRow {
					currentIdx = idx
					break
				}
			}
			// Move to next visible row
			if currentIdx >= 0 && currentIdx < len(table.state.visibleRows)-1 {
				table.state.selectedRow = table.state.visibleRows[currentIdx+1]
			} else if currentIdx < 0 && len(table.state.visibleRows) > 0 {
				// Not found, default to first visible row
				table.state.selectedRow = table.state.visibleRows[0]
			}
		}
	case "left":
		// Only handle column navigation if not in row-only mode
		if !table.config.RowSelectOnlyMode {
			// Find current column display index
			currentDisplayIdx := -1
			for displayIdx, actualIdx := range table.state.visibleColumns {
				if actualIdx == table.state.selectedCol {
					currentDisplayIdx = displayIdx
					break
				}
			}
			// Move to previous column
			if currentDisplayIdx > 0 {
				prevActualIdx := table.state.visibleColumns[currentDisplayIdx-1]
				table.state.selectedCol = prevActualIdx
			} else if currentDisplayIdx < 0 && len(table.state.visibleColumns) > 0 {
				// Initialize to first visible column if not set
				table.state.selectedCol = table.state.visibleColumns[0]
			}
		} else {
		}
	case "right":
		// Only handle column navigation if not in row-only mode
		if !table.config.RowSelectOnlyMode {
			// Find current column display index
			currentDisplayIdx := -1
			for displayIdx, actualIdx := range table.state.visibleColumns {
				if actualIdx == table.state.selectedCol {
					currentDisplayIdx = displayIdx
					break
				}
			}
			// Move to next column
			if currentDisplayIdx >= 0 && currentDisplayIdx < len(table.state.visibleColumns)-1 {
				nextActualIdx := table.state.visibleColumns[currentDisplayIdx+1]
				table.state.selectedCol = nextActualIdx
			} else if currentDisplayIdx < 0 && len(table.state.visibleColumns) > 0 {
				// Initialize to first visible column if not set
				table.state.selectedCol = table.state.visibleColumns[0]
			}
		} else {
		}
	}

	// Log if selection changed
	if oldRow != table.state.selectedRow || oldCol != table.state.selectedCol {

		// Set flag to prevent auto-activation during keyboard navigation
		table.state.isKeyboardNavigation = true

		// Programmatically select the cell in the underlying table to trigger auto-scroll
		if table.table != nil {
			// Map data row index to display row index
			// selectedRow is a data index, we need to find its position in visibleRows
			displayRow := -1
			for displayIdx, dataIdx := range table.state.visibleRows {
				if dataIdx == table.state.selectedRow {
					displayRow = displayIdx + 1 // +1 for header
					break
				}
			}

			// Fallback if not found in visibleRows (shouldn't happen normally)
			if displayRow < 0 {
				displayRow = table.state.selectedRow + 1
			}

			// Map actual column to display column
			displayCol := -1
			for idx, actualCol := range table.state.visibleColumns {
				if actualCol == table.state.selectedCol {
					displayCol = idx
					break
				}
			}

			// If we have a valid display column, select it (or select first column in row-only mode)
			if displayCol >= 0 {
				table.table.Select(widget.TableCellID{Row: displayRow, Col: displayCol})
			} else if table.config.RowSelectOnlyMode && len(table.state.visibleColumns) > 0 {
				// In row-only mode, select first visible column for scrolling purposes
				table.table.Select(widget.TableCellID{Row: displayRow, Col: 0})
			}

			// Refresh to update highlighting - this will re-render all visible cells
			table.table.Refresh()
		}

		// Clear the keyboard navigation flag after selection is complete
		table.state.isKeyboardNavigation = false
	}
}

// handlePageNavigation handles page-based navigation (PgUp, PgDown, Home, End)
func (h *DefaultKeyHandler) handlePageNavigation(direction string, table *Table) {
	// Initialize selection to first row if nothing selected
	if table.state.selectedRow < 0 {
		if len(table.data) > 0 {
			table.state.selectedRow = 0
			if len(table.state.visibleColumns) > 0 {
				table.state.selectedCol = table.state.visibleColumns[0]
			}
		}
		return
	}

	oldRow := table.state.selectedRow
	pageSize := 10 // Number of rows to jump for PgUp/PgDown

	switch direction {
	case "pageup":
		// Move up by pageSize rows
		table.state.selectedRow -= pageSize
		if table.state.selectedRow < 0 {
			table.state.selectedRow = 0
		}
	case "pagedown":
		// Move down by pageSize rows
		table.state.selectedRow += pageSize
		if table.state.selectedRow >= len(table.data) {
			table.state.selectedRow = len(table.data) - 1
		}
	case "home":
		// Jump to first row
		table.state.selectedRow = 0
	case "end":
		// Jump to last row
		if len(table.data) > 0 {
			table.state.selectedRow = len(table.data) - 1
		}
	}

	// Log if selection changed
	if oldRow != table.state.selectedRow {

		// Set flag to prevent auto-activation during keyboard navigation
		table.state.isKeyboardNavigation = true

		// Programmatically select the cell in the underlying table to trigger auto-scroll
		if table.table != nil {
			// Map data row to display row (+1 for header)
			displayRow := table.state.selectedRow + 1

			// Map actual column to display column
			displayCol := -1
			for idx, actualCol := range table.state.visibleColumns {
				if actualCol == table.state.selectedCol {
					displayCol = idx
					break
				}
			}

			// If we have a valid display column, select it
			if displayCol >= 0 {
				table.table.Select(widget.TableCellID{Row: displayRow, Col: displayCol})
			} else if table.config.RowSelectOnlyMode && len(table.state.visibleColumns) > 0 {
				// In row-only mode, select first visible column for scrolling purposes
				table.table.Select(widget.TableCellID{Row: displayRow, Col: 0})
			}

			// Refresh to update highlighting
			table.table.Refresh()
		}

		// Clear the keyboard navigation flag after selection is complete
		table.state.isKeyboardNavigation = false
	}
}

// handleKeyboardEdit starts inline editing on the first editable column of selected row
func (h *DefaultKeyHandler) handleKeyboardEdit(table *Table) {
	if table.state.selectedRow < 0 || table.state.selectedRow >= len(table.data) {
		return
	}

	// Find first editable column (skip columns with popup menus)
	for _, colIndex := range table.state.visibleColumns {
		col := table.config.Columns[colIndex]
		if col.Editable && col.PopupOptions == nil {
			table.startEdit(table.state.selectedRow, colIndex)
			return
		}
	}

}

// handleKeyboardDoubleClickAction triggers the action callback for the selected row/column
func (h *DefaultKeyHandler) handleKeyboardDoubleClickAction(table *Table) {
	if table.state.selectedRow < 0 || table.state.selectedRow >= len(table.data) {
		return
	}

	// Call action callback for the selected column only
	if table.state.selectedCol >= 0 && table.state.selectedCol < len(table.config.Columns) {
		col := table.config.Columns[table.state.selectedCol]
		if col.OnViewData != nil {
			col.OnViewData("ctrl-enter", table.data[table.state.selectedRow], col.ID, table.state.selectedRow, table.state.selectedCol)
		}
	}
}

// handleSpaceKeyPopup shows popup menu for interactive cells
func (h *DefaultKeyHandler) handleSpaceKeyPopup(table *Table) {
	rowIndex := table.state.selectedRow
	colIndex := table.state.selectedCol

	// Validate indices
	if rowIndex < 0 || rowIndex >= len(table.data) || colIndex < 0 || colIndex >= len(table.config.Columns) {
		return
	}

	// Check if column is read-only (non-activatable)
	if table.isColumnReadOnly(colIndex) {
		return
	}

	col := table.config.Columns[colIndex]

	// Check if column has popup options defined
	if col.PopupOptions == nil || col.OnPopupSelected == nil {
		return
	}

	// Get current data item
	dataItem := table.data[rowIndex]

	// Get available options for this cell
	options := col.PopupOptions(dataItem)
	if len(options) == 0 {
		return
	}

	// Show popup menu at current cell position
	table.showPopupMenu(rowIndex, colIndex, options, dataItem)
}

// handleSpaceKeyCheckbox toggles checkbox state for interactive boolean cells
func (h *DefaultKeyHandler) handleSpaceKeyCheckbox(table *Table) {
	rowIndex := table.state.selectedRow
	colIndex := table.state.selectedCol

	// Validate indices
	if rowIndex < 0 || rowIndex >= len(table.data) || colIndex < 0 || colIndex >= len(table.config.Columns) {
		return
	}

	// Check if column is read-only (non-activatable)
	if table.isColumnReadOnly(colIndex) {
		return
	}

	col := table.config.Columns[colIndex]

	// Check if column has checkbox defined
	if !col.ShowCheckbox || col.GetCheckboxValue == nil || col.OnCheckboxChanged == nil {
		return
	}

	// Get current data item
	dataItem := table.data[rowIndex]

	// Get current checkbox state and toggle it
	currentValue := col.GetCheckboxValue(dataItem)
	newValue := !currentValue

	// Call the change handler
	col.OnCheckboxChanged(dataItem, newValue, rowIndex)

	// Trigger OnCellEdited callback if defined
	if table.config.OnCellEdited != nil {
		valueStr := "false"
		if newValue {
			valueStr = "true"
		}
		table.config.OnCellEdited(rowIndex, col.ID, valueStr, dataItem)
	}

	// CRITICAL: Restore navigation state to the cell where checkbox was toggled
	// The OnCellEdited callback may have triggered state changes
	// But we want arrow keys to continue from where the checkbox toggle occurred
	table.logger().Info(fmt.Sprintf("[CHECKBOX-KEY] Restoring navigation state from (%d,%d) to (%d,%d)",
		table.state.selectedRow, table.state.selectedCol, rowIndex, colIndex))
	table.state.selectedRow = rowIndex
	table.state.selectedCol = colIndex

	// Refresh the entire table to show updated state
	// Note: We refresh the whole table instead of just the cell to avoid
	// unintended side effects from RefreshItem (which can cause selection changes)
	// We DON'T call Select() because that would trigger OnSelected callbacks
	// and potentially change our navigation state. Our selectedRow/selectedCol
	// remain unchanged, so arrow keys will work correctly from the current position.
	if table.table != nil {
		table.logger().Info(fmt.Sprintf("[REFRESH-KEY] Refreshing table after checkbox toggle, navigation state: row=%d col=%d", table.state.selectedRow, table.state.selectedCol))
		table.table.Refresh()
		table.logger().Info(fmt.Sprintf("[REFRESH-KEY] Table refreshed, navigation state preserved: row=%d col=%d", table.state.selectedRow, table.state.selectedCol))
	}

	table.logger().Info(fmt.Sprintf("Checkbox toggled: %t -> %t for row %d, col %s", currentValue, newValue, rowIndex, col.ID))
}
