package table

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Ensure Table implements required interfaces
var _ fyne.Widget = (*Table)(nil)
var _ fyne.Focusable = (*Table)(nil)

// escapeableEntry is a custom Entry widget that handles ESC and Enter keys
type escapeableEntry struct {
	widget.Entry
	onEscape func()
	onSubmit func()
}

func newEscapeableEntry(text string, onEscape, onSubmit func()) *escapeableEntry {
	e := &escapeableEntry{
		onEscape: onEscape,
		onSubmit: onSubmit,
	}
	e.ExtendBaseWidget(e)
	e.SetText(text)
	return e
}

func (e *escapeableEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		if e.onEscape != nil {
			e.onEscape()
			return // Don't pass to base Entry
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if e.onSubmit != nil {
			e.onSubmit()
			return // Don't pass to base Entry
		}
	}
	e.Entry.TypedKey(key)
}

// keyboardForwardingTable wraps widget.Table to forward keyboard and focus events
// to the parent Table widget for proper keyboard shortcut handling
type keyboardForwardingTable struct {
	*widget.Table
	onTypedKey      func(*fyne.KeyEvent)
	onTypedShortcut func(fyne.Shortcut)
	onFocusGain     func()
	onFocusLost     func()
	onDoubleTap     func(*fyne.PointEvent)
}

// TypedKey forwards keyboard events to the parent Table handler
func (t *keyboardForwardingTable) TypedKey(key *fyne.KeyEvent) {
	// Forward to our custom handler which handles all key events
	// Note: widget.Table doesn't have a TypedKey method to call
	if t.onTypedKey != nil {
		t.onTypedKey(key)
	}
}

// TypedShortcut forwards shortcut events to parent Table
func (t *keyboardForwardingTable) TypedShortcut(shortcut fyne.Shortcut) {
	// Note: widget.Table doesn't have TypedShortcut, so we only forward to our handler
	if t.onTypedShortcut != nil {
		t.onTypedShortcut(shortcut)
	}
}

// FocusGained forwards focus events to parent Table
func (t *keyboardForwardingTable) FocusGained() {
	// Call base table's FocusGained
	t.Table.FocusGained()
	// Notify our wrapper
	if t.onFocusGain != nil {
		t.onFocusGain()
	}
}

// FocusLost forwards focus events to parent Table
func (t *keyboardForwardingTable) FocusLost() {
	// Call base table's FocusLost
	t.Table.FocusLost()
	// Notify our wrapper
	if t.onFocusLost != nil {
		t.onFocusLost()
	}
}

// DoubleTapped forwards double-tap events to parent Table
func (t *keyboardForwardingTable) DoubleTapped(ev *fyne.PointEvent) {
	// Forward to our wrapper for column auto-resize handling
	if t.onDoubleTap != nil {
		t.onDoubleTap(ev)
	}
}

// Table is a table widget with sorting, filtering, and inline editing
type Table struct {
	widget.BaseWidget

	config *Config
	data   []interface{}

	// Internal widget reference
	table *keyboardForwardingTable

	// Edit widget reference (separate from state as it's a UI object)
	editingEntry *widget.Entry

	// Filter UI widgets (only created if ShowSearch is true)
	filterEntry           *widget.Entry
	regexCheckbox         *widget.Check
	caseSensitiveCheckbox *widget.Check
	clearFilterBtn        *widget.Button
	filterToggleCheckbox  *widget.Check
	checkboxWithBg        *fyne.Container
	filterControlsBox     *fyne.Container
	filterSection         *fyne.Container
	filterTopContainer    *fyne.Container
	filterVisible         bool

	// Runtime state
	state *TableState

	// Event handlers (can be customized)
	KeyHandler   KeyHandler
	MouseHandler MouseHandler
	FocusHandler FocusHandler
}

// NewTable creates a new table
func NewTable(config *Config) *Table {
	st := &Table{
		config:       config,
		data:         []interface{}{},
		state:        NewTableState(),
		KeyHandler:   NewDefaultKeyHandler(),
		MouseHandler: NewDefaultMouseHandler(),
		FocusHandler: NewDefaultFocusHandler(),
	}

	// Build list of visible columns (exclude hidden ones)
	st.RebuildVisibleColumns()
	// Build list of visible rows (apply tree filtering)
	st.RebuildVisibleRows()

	st.createTable()
	st.createFilterUI()
	st.ExtendBaseWidget(st) // Must be called AFTER createFilterUI so renderer sees filter section

	// Force initial refresh to ensure all rows display correctly
	if st.table != nil {
		st.table.Refresh()
	}

	// Verify ShowHeaderColumn is still false after initialization
	if st.table != nil {
		st.logger().Info(fmt.Sprintf("[TABLE] Post-init verification: ShowHeaderColumn=%v",
			st.table.ShowHeaderColumn))
	}

	return st
}

// logger returns the configured logger or a default one
func (st *Table) logger() Logger {
	if st.config.Logger != nil {
		return st.config.Logger
	}
	return NoopLogger{}
}

// isColumnReadOnly checks if a column is marked as read-only (non-activatable)
func (st *Table) isColumnReadOnly(col int) bool {
	if col < 0 || col >= len(st.config.Columns) {
		return true // Invalid column is read-only
	}
	return st.config.Columns[col].ReadOnly
}

// createTable initializes the Fyne table widget
func (st *Table) createTable() {
	baseTable := widget.NewTable(
		st.tableLength,
		st.tableCreateCell,
		st.tableUpdateCell,
	)

	// Wrap the table to forward keyboard and focus events
	st.table = &keyboardForwardingTable{
		Table:           baseTable,
		onTypedKey:      st.TypedKey,
		onTypedShortcut: st.TypedShortcut,
		onFocusGain:     st.FocusGained,
		onFocusLost:     st.FocusLost,
		onDoubleTap:     st.handleDoubleTap,
	}

	// OnSelected is triggered by single click in Fyne
	st.table.OnSelected = func(id widget.TableCellID) {
		st.logger().Info(fmt.Sprintf("[ONSELECTED] Callback triggered: id={Row:%d,Col:%d} isReselecting=%v",
			id.Row, id.Col, st.state.isReselecting))
		st.handleCellClick(id)
		st.logger().Info(fmt.Sprintf("[ONSELECTED] Callback completed: isReselecting=%v", st.state.isReselecting))
	}

	// Make header row sticky (doesn't scroll)
	st.table.StickyRowCount = 1
	st.table.ShowHeaderRow = st.config.ShowHeaders // Control header visibility
	st.table.ShowHeaderColumn = false              // Hide Fyne's default A-D column labels
	st.logger().Info(fmt.Sprintf("[TABLE] ShowHeaderColumn set to false, ShowHeaderRow=%v",
		st.config.ShowHeaders))

	// Disable manual column resizing if configured
	// Note: Fyne doesn't expose a direct API to disable manual resize,
	// but setting ShowHeaderRow=false effectively disables it
	// Additional approach: could intercept drag events if needed

	// Set column widths (only for visible columns)
	for displayIdx, actualIdx := range st.state.visibleColumns {
		col := st.config.Columns[actualIdx]
		if col.Width > 0 {
			st.table.SetColumnWidth(displayIdx, col.Width)
		}
	}

	// Set row heights
	if st.config.HeaderHeight > 0 {
		st.table.SetRowHeight(0, st.config.HeaderHeight)
	}
	// Note: Fyne doesn't support per-row heights easily, default row height applies to all data rows
}

// createFilterUI creates the search/filter UI controls
func (st *Table) createFilterUI() {
	if !st.config.ShowSearch {
		st.logger().Info("[FILTER] ShowSearch is false, skipping filter UI creation")
		return
	}

	st.logger().Info("[FILTER] Creating filter UI widgets")
	st.filterVisible = false // Start with filter hidden (can be toggled via external controls)

	// Create filter entry with placeholder
	st.filterEntry = widget.NewEntry()
	placeholder := st.config.SearchPlaceholder
	if placeholder == "" {
		placeholder = "Search..."
	}
	st.filterEntry.SetPlaceHolder(placeholder)

	// Create regex checkbox (forward declare needed for filterEntry callback)
	st.regexCheckbox = widget.NewCheck("Regex", func(checked bool) {
		if st.filterEntry != nil {
			st.SetFilter(st.filterEntry.Text, checked)
			st.RequestFocus()
		}
	})
	st.regexCheckbox.SetChecked(false)

	// Create case sensitive checkbox
	st.caseSensitiveCheckbox = widget.NewCheck("Match case", func(checked bool) {
		st.SetFilterCaseSensitive(checked)
		st.RequestFocus()
	})
	st.caseSensitiveCheckbox.SetChecked(false)

	// Create clear filter button
	st.clearFilterBtn = widget.NewButton("Clear", func() {
		if st.filterEntry != nil {
			st.filterEntry.SetText("")
			st.ClearFilter()
			st.RequestFocus()
		}
	})

	// Wire up filter entry callbacks
	st.filterEntry.OnChanged = func(text string) {
		useRegex := false
		if st.regexCheckbox != nil {
			useRegex = st.regexCheckbox.Checked
		}
		st.SetFilter(text, useRegex)
	}
	st.filterEntry.OnSubmitted = func(text string) {
		st.RequestFocus()
	}

	// Create the filter controls container
	st.filterControlsBox = container.NewBorder(
		nil,
		nil,
		nil,
		container.NewHBox(st.regexCheckbox, st.caseSensitiveCheckbox, st.clearFilterBtn),
		st.filterEntry,
	)

	// Create toggle checkbox to show/hide filter
	st.filterToggleCheckbox = widget.NewCheck("Search/Filter", func(checked bool) {
		st.filterVisible = checked
		st.updateFilterVisibility()
	})
	st.filterToggleCheckbox.SetChecked(false) // Start unchecked (filter hidden)

	// Use checkbox directly (for backwards compatibility, checkboxWithBg points to same widget)
	st.checkboxWithBg = container.NewStack(st.filterToggleCheckbox)

	// Create the filter section container that will hold the filter controls (no checkbox)
	// The checkbox is accessed separately via GetFilterToggleCheckbox() for Settings dialog
	st.filterSection = container.NewVBox()   // Start empty since filterVisible = false
	st.filterTopContainer = st.filterSection // For backwards compatibility

	st.logger().Info(fmt.Sprintf("[FILTER] Filter UI created: filterSection=%v, checkboxWithBg=%v", st.filterSection != nil, st.checkboxWithBg != nil))
}

// ========================================
// Handler Routing Functions
// ========================================

// TypedKey handles keyboard events - delegates to KeyHandler
func (st *Table) TypedKey(key *fyne.KeyEvent) {
	if st.KeyHandler != nil {
		st.KeyHandler.HandleKey(key, st)
	}
}

// TypedShortcut handles keyboard shortcuts - delegates to KeyHandler
func (st *Table) TypedShortcut(shortcut fyne.Shortcut) {
	if st.KeyHandler != nil {
		st.KeyHandler.HandleShortcut(shortcut, st)
	}
}

// TypedRune implements Focusable interface
func (st *Table) TypedRune(r rune) {
	// Not handling typed runes currently
}

// FocusGained implements Focusable interface - delegates to FocusHandler
func (st *Table) FocusGained() {
	if st.FocusHandler != nil {
		st.FocusHandler.HandleFocusGained(st)
	}
}

// FocusLost implements Focusable interface - delegates to FocusHandler
func (st *Table) FocusLost() {
	if st.FocusHandler != nil {
		st.FocusHandler.HandleFocusLost(st)
	}
}

// Tapped handles tap events - delegates focus request to FocusHandler
func (st *Table) Tapped(*fyne.PointEvent) {
	if st.FocusHandler != nil {
		st.FocusHandler.RequestFocus(st)
	}
}

// handleCellClick handles clicking on a cell - selects row but doesn't start editing
func (st *Table) handleCellClick(id widget.TableCellID) {
	if st.MouseHandler != nil {
		st.MouseHandler.HandleCellClick(id, st)
	}
}

// handleHeaderClick handles clicking on a column header to sort
func (st *Table) handleHeaderClick(colIndex int) {
	if st.MouseHandler != nil {
		st.MouseHandler.HandleHeaderClick(colIndex, st)
	}
}

// handleDoubleTap handles double-tap/double-click events
func (st *Table) handleDoubleTap(ev *fyne.PointEvent) {
	if st.MouseHandler != nil {
		st.MouseHandler.HandleDoubleTap(ev, st)
	}
}

// ========================================
// Public Functions
// ========================================

// RebuildVisibleColumns rebuilds the list of visible column indices
func (st *Table) RebuildVisibleColumns() {
	st.state.visibleColumns = make([]int, 0, len(st.config.Columns))
	for i, col := range st.config.Columns {
		if !col.Hidden {
			st.state.visibleColumns = append(st.state.visibleColumns, i)
		}
	}
}

// SetColumnVisibility sets whether a column is visible or hidden
func (st *Table) SetColumnVisibility(columnID string, visible bool) {
	// Find the column by ID
	for i := range st.config.Columns {
		if st.config.Columns[i].ID == columnID {
			st.config.Columns[i].Hidden = !visible
			st.RebuildVisibleColumns()

			// Validate selectedCol is still in visibleColumns
			if st.state.selectedCol >= 0 {
				isVisible := false
				for _, visCol := range st.state.visibleColumns {
					if visCol == st.state.selectedCol {
						isVisible = true
						break
					}
				}
				// If current selectedCol is now hidden, find next visible cell
				if !isVisible && len(st.state.visibleColumns) > 0 {
					// Try to find next visible column to the right
					nextCol := -1
					for _, visCol := range st.state.visibleColumns {
						if visCol > st.state.selectedCol {
							nextCol = visCol
							break
						}
					}

					// If no column to the right, wrap to first visible column
					if nextCol == -1 {
						nextCol = st.state.visibleColumns[0]
					}

					st.state.selectedCol = nextCol
					st.logger().Info(fmt.Sprintf("Selected column was hidden, moved to next visible col %d", st.state.selectedCol))

					// Update the visual selection in Fyne's table
					if st.state.selectedRow >= 0 && st.table != nil {
						st.SetSelectedCell(st.state.selectedRow, st.state.selectedCol)
					}
				}
			}

			// Instead of recreating, just refresh the existing table
			// Fyne's table widget will call tableLength() which returns updated column count
			if st.table != nil {
				// Update column widths for visible columns
				for displayIdx, actualIdx := range st.state.visibleColumns {
					col := st.config.Columns[actualIdx]
					if col.Width > 0 {
						st.table.SetColumnWidth(displayIdx, col.Width)
					}
				}
				// Refresh the table to apply changes
				st.table.Refresh()
			}
			return
		}
	}
}

// SetColumnReadOnly sets whether a column is read-only (non-activatable)
func (st *Table) SetColumnReadOnly(columnID string, readOnly bool) {
	for i := range st.config.Columns {
		if st.config.Columns[i].ID == columnID {
			st.config.Columns[i].ReadOnly = readOnly
			st.logger().Info(fmt.Sprintf("Column %s ReadOnly set to %t", columnID, readOnly))
			return
		}
	}
}

// SetRowSelectOnlyMode updates the row selection mode and initializes column selection if needed
func (st *Table) SetRowSelectOnlyMode(rowOnlyMode bool) {
	st.config.RowSelectOnlyMode = rowOnlyMode

	// If switching to row-column mode, ensure selectedCol is initialized
	if !rowOnlyMode && st.state.selectedCol < 0 && len(st.state.visibleColumns) > 0 {
		st.state.selectedCol = st.state.visibleColumns[0]
		st.logger().Info(fmt.Sprintf("Switched to row-column mode, initialized selectedCol to %d", st.state.selectedCol))
	}

	// Refresh to update highlighting
	if st.table != nil {
		st.table.Refresh()
	}
}

// RebuildVisibleRows rebuilds the list of visible row indices based on tree state and filter
func (st *Table) RebuildVisibleRows() {
	st.state.visibleRows = make([]int, 0, len(st.data))

	// Build filter regex if needed
	var filterRegex *regexp.Regexp
	if st.state.filterText != "" && st.state.filterRegex {
		var err error
		regexPattern := st.state.filterText
		if !st.state.filterCaseSensitive {
			regexPattern = "(?i)" + regexPattern // Case-insensitive flag
		}
		filterRegex, err = regexp.Compile(regexPattern)
		if err != nil {
			st.logger().Error(fmt.Sprintf("Invalid regex filter: %v", err))
			filterRegex = nil
		}
	}

	// Iterate through all data and apply filters
	for i := range st.data {
		// Apply text filter if configured
		if st.state.filterText != "" && len(st.config.FilterColumns) > 0 {
			matched := false
			for _, colID := range st.config.FilterColumns {
				fieldValue := st.extractFieldValue(st.data[i], colID)

				if st.state.filterRegex && filterRegex != nil {
					// Regex matching
					if filterRegex.MatchString(fieldValue) {
						matched = true
						break
					}
				} else {
					// Plain text substring matching
					if st.state.filterCaseSensitive {
						// Case-sensitive
						if strings.Contains(fieldValue, st.state.filterText) {
							matched = true
							break
						}
					} else {
						// Case-insensitive
						if strings.Contains(strings.ToLower(fieldValue), strings.ToLower(st.state.filterText)) {
							matched = true
							break
						}
					}
				}
			}
			if !matched {
				continue // Skip this row
			}
		}

		// TODO: Tree expansion/collapse filtering would go here when enabled
		// For now, include all rows that passed the text filter

		// Apply tree depth filter if MaxDepth is set
		if st.config.MaxDepth > 0 && st.config.GetNodeDepth != nil {
			depth := st.config.GetNodeDepth(st.data[i])
			if depth >= st.config.MaxDepth {
				continue // Skip rows beyond max depth
			}
		}

		st.state.visibleRows = append(st.state.visibleRows, i)
	}
}

// ToggleNodeExpansion toggles the expansion state of a node
func (st *Table) ToggleNodeExpansion(rowIndex int) {
	if rowIndex < 0 || rowIndex >= len(st.data) {
		return
	}

	if st.config.GetNodeID == nil {
		return // Tree functions not configured
	}

	item := st.data[rowIndex]
	nodeID := st.config.GetNodeID(item)

	// Toggle expansion state
	st.config.ExpandedNodes[nodeID] = !st.config.ExpandedNodes[nodeID]

	// Rebuild visible rows and refresh
	st.RebuildVisibleRows()
	st.createTable()
	st.Refresh()
}

// SetMaxDepth sets the maximum tree depth to display
func (st *Table) SetMaxDepth(maxDepth int) {
	st.config.MaxDepth = maxDepth
	st.RebuildVisibleRows()
	if st.table != nil {
		// Refresh base table to ensure it detects row count changes
		// Use DoAndWait to ensure refresh runs on UI thread
		fyne.DoAndWait(func() {
			st.table.Refresh()
			st.table.Refresh()
		})
	}
	st.logger().Info(fmt.Sprintf("SetMaxDepth: %d, visible rows=%d/%d", maxDepth, len(st.state.visibleRows), len(st.data)))
}

// RequestFocus requests keyboard focus - delegates to FocusHandler
func (st *Table) RequestFocus() {
	if st.FocusHandler != nil {
		st.FocusHandler.RequestFocus(st)
	}
}

// GetSelectedCell returns the currently selected row and column (-1 if none selected)
func (st *Table) GetSelectedCell() (row int, col int) {
	return st.state.selectedRow, st.state.selectedCol
}

// GetEditingState returns true if the table is currently in inline editing mode
func (st *Table) GetEditingState() bool {
	return st.state.editingRow >= 0 && st.state.editingCol >= 0
}

// SetSelectedCell programmatically selects a cell and triggers visual update
func (st *Table) SetSelectedCell(row int, col int) {
	if row < 0 || row >= len(st.data) {
		st.logger().Error(fmt.Sprintf("SetSelectedCell: invalid row %d (data has %d rows)", row, len(st.data)))
		return
	}

	st.state.selectedRow = row
	st.state.selectedCol = col

	// Refresh to update highlighting
	if st.table != nil {
		st.table.Refresh()
	}
}

// ========== Multi-Select API ==========

// GetSelectedRows returns a slice of selected row indices (works for both single and multi-select)
func (st *Table) GetSelectedRows() []int {
	return st.state.GetSelectedRows()
}

// SetSelectedRows sets multiple rows as selected (enables multi-select mode)
func (st *Table) SetSelectedRows(rows []int) {
	st.state.SetSelectedRows(rows)

	// Refresh to update highlighting
	if st.table != nil {
		st.table.Refresh()
	}
}

// ClearSelection clears all selected rows
func (st *Table) ClearSelection() {
	st.state.ClearSelection()

	// Refresh to update highlighting
	if st.table != nil {
		st.table.Refresh()
	}
}

// GetSelectionCount returns the number of selected rows
func (st *Table) GetSelectionCount() int {
	return st.state.GetSelectionCount()
}

// ========================================

// SetData updates the table data
func (st *Table) SetData(data []interface{}) {
	st.data = data

	// Re-apply current sort if one is active
	if st.state.sortColumn >= 0 && st.state.sortColumn < len(st.config.Columns) {
		st.logger().Info(fmt.Sprintf("[SETDATA] Re-applying sort: column=%d asc=%v",
			st.state.sortColumn, st.state.sortAsc))
		st.sortData()
	}

	st.RebuildVisibleRows() // Update visible rows based on tree state
	if st.table != nil {
		st.table.Refresh()
	}

	// Auto-select first cell if configured and data exists
	if st.config.SelectFirstCellOnStartup && len(data) > 0 && len(st.state.visibleColumns) > 0 {
		st.SetSelectedCell(0, st.state.visibleColumns[0])
		st.RequestFocus()
	}
}

// GetData returns the current data
func (st *Table) GetData() []interface{} {
	return st.data
}

// GetUnderlyingTable returns the base Fyne table widget for direct access
func (st *Table) GetUnderlyingTable() *widget.Table {
	if st.table != nil {
		return st.table.Table
	}
	return nil
}

// SetFilter updates the filter text and rebuilds visible rows
func (st *Table) SetFilter(filterText string, useRegex bool) {
	st.state.filterText = filterText
	st.state.filterRegex = useRegex
	st.RebuildVisibleRows()
	if st.table != nil {
		// Refresh base table to ensure it detects row count changes
		// Use Do (not DoAndWait) to avoid deadlock when called from UI thread
		fyne.Do(func() {
			st.table.Refresh()
			st.table.Refresh()
		})
	}
	st.logger().Info(fmt.Sprintf("Filter set: text=%q, regex=%v, caseSensitive=%v, visible rows=%d/%d", filterText, useRegex, st.state.filterCaseSensitive, len(st.state.visibleRows), len(st.data)))
}

// SetFilterCaseSensitive sets whether filtering is case-sensitive
func (st *Table) SetFilterCaseSensitive(caseSensitive bool) {
	st.state.filterCaseSensitive = caseSensitive
	if st.state.filterText != "" {
		st.RebuildVisibleRows()
		if st.table != nil {
			// Refresh base table to ensure it detects row count changes
			// Use Do (not DoAndWait) to avoid deadlock when called from UI thread
			fyne.Do(func() {
				st.table.Refresh()
				st.table.Refresh()
			})
		}
	}
}

// GetFilter returns the current filter text and regex mode
func (st *Table) GetFilter() (string, bool) {
	return st.state.filterText, st.state.filterRegex
}

// ClearFilter removes the filter and shows all rows
func (st *Table) ClearFilter() {
	st.SetFilter("", false)
}

// GetFilterToggleCheckbox returns the filter toggle checkbox for external configuration
func (st *Table) GetFilterToggleCheckbox() *widget.Check {
	return st.filterToggleCheckbox
}

// GetFilterVisible returns the current visibility state of the filter controls
func (st *Table) GetFilterVisible() bool {
	return st.filterVisible
}

// SetFilterVisible sets the visibility of the filter controls
func (st *Table) SetFilterVisible(visible bool) {
	if !st.config.ShowSearch || st.filterToggleCheckbox == nil {
		return
	}

	st.filterVisible = visible
	st.filterToggleCheckbox.SetChecked(visible)
	st.updateFilterVisibility()
}

// updateFilterVisibility updates the filter controls visibility
func (st *Table) updateFilterVisibility() {
	if st.filterSection == nil {
		return
	}

	if st.filterVisible {
		// Show the filter controls (no checkbox in main display)
		st.filterSection.Objects = []fyne.CanvasObject{
			st.filterControlsBox,
		}
	} else {
		// Hide the filter controls
		st.filterSection.Objects = []fyne.CanvasObject{}
	}
	st.filterSection.Refresh()
}

// CreateRenderer implements fyne.Widget
func (st *Table) CreateRenderer() fyne.WidgetRenderer {
	st.logger().Info(fmt.Sprintf("[FILTER] CreateRenderer called: ShowSearch=%v, filterSection=%v", st.config.ShowSearch, st.filterSection != nil))
	// If search/filter UI is enabled, show it above the table
	if st.config.ShowSearch && st.filterSection != nil {
		st.logger().Info("[FILTER] Creating renderer WITH filter section")
		content := container.NewBorder(
			st.filterSection, // top
			nil,              // bottom
			nil,              // left
			nil,              // right
			st.table,         // center
		)
		return widget.NewSimpleRenderer(content)
	}

	// Otherwise just return the table directly
	st.logger().Info("[FILTER] Creating renderer WITHOUT filter section")
	return widget.NewSimpleRenderer(st.table)
}

// ========================================
// Private Functions
// ========================================

// tableLength returns the number of rows and columns
func (st *Table) tableLength() (int, int) {
	if st.state.visibleRows != nil {
		rows := len(st.state.visibleRows) + 1 // +1 for header row
		cols := len(st.state.visibleColumns)  // Only count visible columns
		return rows, cols
	}
	rows := len(st.data) + 1             // +1 for header row (fallback)
	cols := len(st.state.visibleColumns) // Only count visible columns
	return rows, cols
}

// tableCreateCell creates a reusable cell widget
func (st *Table) tableCreateCell() fyne.CanvasObject {
	// Create a container that can hold various content (label, entry, buttons, etc.)
	// If custom font size is configured, use canvas.Text instead of widget.Label
	if st.config.FontSize > 0 {
		text := canvas.NewText("", theme.Color(theme.ColorNameForeground))
		text.TextSize = st.config.FontSize
		return container.NewStack(text)
	}
	return container.NewStack(widget.NewLabel(""))
}

// tableUpdateCell updates a cell with data
func (st *Table) tableUpdateCell(id widget.TableCellID, cell fyne.CanvasObject) {
	container := cell.(*fyne.Container)

	// Header row (row 0)
	if id.Row == 0 {
		st.renderHeaderCell(id.Col, container)
		return
	}

	// Data rows (row 1+)
	displayRowIndex := id.Row - 1

	// Map display row index to actual data index (if tree filtering is active)
	var dataIndex int
	if st.state.visibleRows != nil && displayRowIndex < len(st.state.visibleRows) {
		dataIndex = st.state.visibleRows[displayRowIndex]
	} else {
		dataIndex = displayRowIndex // Fallback if no filtering
	}

	if dataIndex >= len(st.data) {
		// Empty cell
		container.Objects = []fyne.CanvasObject{widget.NewLabel("")}
		container.Refresh()
		return
	}

	// Note: Removed excessive logging during cell updates (was causing spam during column resize)
	// To debug cell rendering, temporarily uncomment the line below:
	// st.logger().Info(fmt.Sprintf("[UPDATE] tableUpdateCell displayRow=%d dataIndex=%d data=%v",
	// 	displayRowIndex, dataIndex, st.data[dataIndex]))

	st.renderDataCell(id.Col, dataIndex, container)
}

// renderHeaderCell renders a header cell
func (st *Table) renderHeaderCell(displayColIndex int, container *fyne.Container) {
	// Map display column index to actual column index
	if displayColIndex >= len(st.state.visibleColumns) {
		container.Objects = []fyne.CanvasObject{widget.NewLabel("")}
		container.Refresh()
		return
	}
	colIndex := st.state.visibleColumns[displayColIndex]

	if colIndex >= len(st.config.Columns) {
		container.Objects = []fyne.CanvasObject{widget.NewLabel("")}
		container.Refresh()
		return
	}

	col := st.config.Columns[colIndex]

	// Build header text with sort indicator
	headerText := col.Title
	if st.state.sortColumn == colIndex {
		if st.state.sortAsc {
			headerText += " ▲"
		} else {
			headerText += " ▼"
		}
	}

	label := widget.NewLabel(headerText)
	label.TextStyle = fyne.TextStyle{Bold: true}

	// Apply column alignment to header
	switch col.Alignment {
	case AlignCenter:
		label.Alignment = fyne.TextAlignCenter
	case AlignRight:
		label.Alignment = fyne.TextAlignTrailing
	default:
		label.Alignment = fyne.TextAlignLeading
	}

	// Make header clickable if column is sortable
	if col.Sortable {
		button := widget.NewButton(headerText, func() {
			st.handleHeaderClick(displayColIndex)
		})
		button.Importance = widget.LowImportance
		// Apply alignment to button as well
		switch col.Alignment {
		case AlignCenter:
			button.Alignment = widget.ButtonAlignCenter
		case AlignRight:
			button.Alignment = widget.ButtonAlignTrailing
		default:
			button.Alignment = widget.ButtonAlignLeading
		}
		container.Objects = []fyne.CanvasObject{button}
	} else {
		container.Objects = []fyne.CanvasObject{label}
	}

	container.Refresh()
}

// renderDataCell renders a data cell (private helper)
func (st *Table) renderDataCell(displayColIndex int, dataIndex int, cellContainer *fyne.Container) {
	// Map display column index to actual column index
	if displayColIndex >= len(st.state.visibleColumns) {
		cellContainer.Objects = []fyne.CanvasObject{widget.NewLabel("")}
		cellContainer.Refresh()
		return
	}
	colIndex := st.state.visibleColumns[displayColIndex]

	if colIndex >= len(st.config.Columns) || dataIndex >= len(st.data) {
		cellContainer.Objects = []fyne.CanvasObject{widget.NewLabel("")}
		cellContainer.Refresh()
		return
	}

	col := st.config.Columns[colIndex]
	data := st.data[dataIndex]

	// Check if this cell is being edited
	if st.state.editingRow == dataIndex && st.state.editingCol == colIndex {
		st.logger().Info(fmt.Sprintf("[DEBUG] renderDataCell: Rendering EDIT widget for row=%d col=%d", dataIndex, colIndex))
		// Show entry widget for editing with ESC/Enter handling
		var escEntry *escapeableEntry
		if st.editingEntry == nil {
			// Create escapeable entry that handles ESC and Enter keys
			escEntry = newEscapeableEntry(st.state.editingValue, st.cancelEdit, st.saveEdit)
			st.editingEntry = &escEntry.Entry // Keep reference to underlying Entry
		} else {
			// Recreate the escapeable entry with current entry value
			escEntry = newEscapeableEntry(st.editingEntry.Text, st.cancelEdit, st.saveEdit)
			st.editingEntry = &escEntry.Entry
		}

		cellContainer.Objects = []fyne.CanvasObject{escEntry}
		cellContainer.Refresh()

		// Focus the entry widget after it's in the canvas tree
		// Use goroutine with small delay to ensure Refresh completes and canvas is ready
		go func() {
			// Small delay to ensure canvas is ready (50ms should be sufficient)
			time.Sleep(50 * time.Millisecond)
			// All UI operations must run on the UI thread
			// Use Do (not DoAndWait) to avoid deadlock
			fyne.Do(func() {
				if canvas := fyne.CurrentApp().Driver().CanvasForObject(escEntry); canvas != nil {
					canvas.Focus(escEntry)
					// Select all text so user can immediately type to replace
					// Use Ctrl+A shortcut to trigger selection
					escEntry.TypedShortcut(&desktop.CustomShortcut{
						KeyName:  fyne.KeyA,
						Modifier: fyne.KeyModifierControl,
					})
				} else {
					st.logger().Warn("Could not get canvas for edit entry (after 50ms delay)")
				}
			})
		}()
		return
	}

	// If custom renderer provided, use it
	if col.Renderer != nil {
		col.Renderer(data, cellContainer, dataIndex, col.ID)
		return
	}

	// Default renderer: extract and display the specific field
	fieldValue := st.extractFieldValue(data, col.ID)

	// Determine if this cell should be highlighted FIRST
	highlightCell := false
	// Check if this row is selected (handles both single and multi-select)
	if st.state.IsRowSelected(dataIndex) {
		if st.config.RowSelectOnlyMode {
			// In row-only mode: highlight ALL visible columns on the selected row
			highlightCell = true
		} else {
			// In row-column mode: highlight only the selected cell
			if st.state.selectedCol == colIndex {
				// Verify selectedCol is actually visible
				for _, visCol := range st.state.visibleColumns {
					if visCol == st.state.selectedCol {
						highlightCell = true
						break
					}
				}
			}
		}
	}

	// Create or update the content widget
	var content fyne.CanvasObject

	// If custom font size is configured, use canvas.Text instead of widget.Label
	if st.config.FontSize > 0 {
		var text *canvas.Text
		// Try to reuse existing text widget (unwrap from container if highlighted)
		if len(cellContainer.Objects) > 0 {
			if maxContainer, ok := cellContainer.Objects[0].(*fyne.Container); ok {
				// Previously highlighted - get content from inside
				if len(maxContainer.Objects) > 1 {
					text, _ = maxContainer.Objects[1].(*canvas.Text)
				}
			} else {
				// Not highlighted - direct content
				text, _ = cellContainer.Objects[0].(*canvas.Text)
			}
		}
		if text == nil {
			text = canvas.NewText(fieldValue, theme.Color(theme.ColorNameForeground))
		} else {
			text.Text = fieldValue
		}
		text.TextSize = st.config.FontSize

		// Apply text alignment
		switch col.Alignment {
		case AlignCenter:
			text.Alignment = fyne.TextAlignCenter
		case AlignRight:
			text.Alignment = fyne.TextAlignTrailing
		default:
			text.Alignment = fyne.TextAlignLeading
		}
		content = text
	} else {
		// Use standard widget.Label for default size
		var label *widget.Label
		// Try to reuse existing label widget (unwrap from container if highlighted)
		if len(cellContainer.Objects) > 0 {
			if maxContainer, ok := cellContainer.Objects[0].(*fyne.Container); ok {
				// Previously highlighted - get content from inside
				if len(maxContainer.Objects) > 1 {
					label, _ = maxContainer.Objects[1].(*widget.Label)
				}
			} else {
				// Not highlighted - direct content
				label, _ = cellContainer.Objects[0].(*widget.Label)
			}
		}
		if label == nil {
			label = widget.NewLabel("")
		}
		label.SetText(fieldValue)

		// Apply text alignment
		switch col.Alignment {
		case AlignCenter:
			label.Alignment = fyne.TextAlignCenter
		case AlignRight:
			label.Alignment = fyne.TextAlignTrailing
		default:
			label.Alignment = fyne.TextAlignLeading
		}
		content = label
	}

	// Apply or remove highlighting based on selection state
	if highlightCell {
		// Wrap the content with selection background using Max container (stacks objects)
		selectionBg := canvas.NewRectangle(theme.Color(theme.ColorNameSelection))

		// Add border around the row (thin line on edges)
		// Determine which borders this cell needs based on its position
		isFirstCol := displayColIndex == 0
		isLastCol := displayColIndex == len(st.state.visibleColumns)-1

		// Create border lines (1 pixel width)
		borderColor := theme.Color(theme.ColorNamePrimary)
		borderWidth := float32(1)

		// Top border on all cells in selected row
		topBorder := canvas.NewRectangle(borderColor)
		topBorder.Resize(fyne.NewSize(0, borderWidth)) // Width will be set by container

		// Bottom border on all cells in selected row
		bottomBorder := canvas.NewRectangle(borderColor)
		bottomBorder.Resize(fyne.NewSize(0, borderWidth))

		// Left border only on first column
		var leftBorder fyne.CanvasObject
		if isFirstCol {
			leftBorder = canvas.NewRectangle(borderColor)
			leftBorder.Resize(fyne.NewSize(borderWidth, 0))
		}

		// Right border only on last column
		var rightBorder fyne.CanvasObject
		if isLastCol {
			rightBorder = canvas.NewRectangle(borderColor)
			rightBorder.Resize(fyne.NewSize(borderWidth, 0))
		}

		// Create a custom container that layers: background, borders, content
		cellContainer.Objects = []fyne.CanvasObject{
			container.NewStack(
				selectionBg,
				container.NewBorder(topBorder, bottomBorder, leftBorder, rightBorder, content),
			),
		}
	} else {
		// No highlighting - just the content
		cellContainer.Objects = []fyne.CanvasObject{content}
	}

	cellContainer.Refresh()
}

// sortData sorts the data based on current sort column and direction
func (st *Table) sortData() {
	if st.state.sortColumn < 0 || st.state.sortColumn >= len(st.config.Columns) {
		return
	}

	col := st.config.Columns[st.state.sortColumn]

	st.logger().Info(fmt.Sprintf("[SORT] sortData called: sortColumn=%d (ID='%s', Title='%s') sortAsc=%v dataLen=%d",
		st.state.sortColumn, col.ID, col.Title, st.state.sortAsc, len(st.data)))

	// Use custom comparator if provided, otherwise use default string comparator
	comparator := col.Comparator
	if comparator == nil {
		st.logger().Info(fmt.Sprintf("[SORT] Using default STRING comparator for column '%s'", col.ID))
		comparator = NewStringComparator(col.ID)
	} else {
		st.logger().Info(fmt.Sprintf("[SORT] Using CUSTOM comparator for column '%s'", col.ID))
	}

	sort.Slice(st.data, func(i, j int) bool {
		cmpResult := comparator(st.data[i], st.data[j])
		if st.state.sortAsc {
			return cmpResult < 0
		}
		return cmpResult > 0
	})

	if len(st.data) > 0 {
		firstItem := fmt.Sprintf("%v", st.data[0])
		st.logger().Info(fmt.Sprintf("[SORT] Sort complete, firstItem=%s", firstItem))
	}
}

// startEdit begins editing a cell
func (st *Table) startEdit(dataIndex int, colIndex int) {
	st.logger().Info(fmt.Sprintf("[DEBUG] startEdit called: dataIndex=%d colIndex=%d", dataIndex, colIndex))

	st.state.editingRow = dataIndex
	st.state.editingCol = colIndex

	// Store original value - extract the specific field
	data := st.data[dataIndex]
	colID := st.config.Columns[colIndex].ID

	// Try to extract field value if possible
	// For prototype: use reflection to get field by column ID
	st.state.editingValue = st.extractFieldValue(data, colID)

	st.logger().Info(fmt.Sprintf("[DEBUG] Editing value: %s", st.state.editingValue))

	// Find the display column index for this actual column index
	displayColIndex := -1
	for i, visCol := range st.state.visibleColumns {
		if visCol == colIndex {
			displayColIndex = i
			break
		}
	}

	if displayColIndex < 0 {
		st.logger().Error(fmt.Sprintf("Column %d not in visibleColumns, cannot edit", colIndex))
		st.state.editingRow = -1
		st.state.editingCol = -1
		return
	}

	st.logger().Info(fmt.Sprintf("[DEBUG] Display column index: %d (for actual col %d)", displayColIndex, colIndex))

	// Refresh the specific cell to trigger renderDataCell with editing state
	if st.table != nil {
		// TableCellID uses display indices: Row is dataIndex+1 (header is row 0), Col is display column
		cellID := widget.TableCellID{Row: dataIndex + 1, Col: displayColIndex}
		st.logger().Info(fmt.Sprintf("[DEBUG] Refreshing cell: row=%d col=%d", cellID.Row, cellID.Col))
		st.table.RefreshItem(cellID)
	}
}

// saveEdit saves the edited value
func (st *Table) saveEdit() {
	if st.state.editingRow < 0 || st.state.editingCol < 0 || st.editingEntry == nil {
		return
	}

	editedRow := st.state.editingRow
	editedCol := st.state.editingCol
	newValue := st.editingEntry.Text
	data := st.data[st.state.editingRow]
	col := st.config.Columns[st.state.editingCol]

	// Call OnViewData callback if provided (for "return" action confirmation)
	if col.OnViewData != nil {
		col.OnViewData("return", data, col.ID, editedRow, editedCol)
	}

	// Call OnCellEdited callback if provided
	if st.config.OnCellEdited != nil {
		st.config.OnCellEdited(st.state.editingRow, col.ID, newValue, data)
	}

	// Note: For prototype, we don't modify st.data directly since that would lose the struct type
	// The callback handler should update the underlying data structure
	// In production, this would have proper type-aware field updates

	// Clear editing state (don't call cancelEdit to avoid triggering OnViewData again)
	st.state.editingRow = -1
	st.state.editingCol = -1
	st.editingEntry = nil
	st.state.editingValue = ""

	if st.table != nil {
		// Refresh the entire table to restore custom renderers
		st.table.Refresh()
		// Also refresh the specific cell to ensure custom renderer is applied
		if editedRow >= 0 && editedCol >= 0 {
			st.table.RefreshItem(widget.TableCellID{Row: editedRow + 1, Col: editedCol})
		}
	}

	// Restore focus and selection to the table
	st.SetSelectedCell(editedRow, editedCol)
	st.RequestFocus()
}

// cancelEdit cancels editing and restores original value
func (st *Table) cancelEdit() {
	editedRow := st.state.editingRow
	editedCol := st.state.editingCol

	// Call OnViewData callback if provided (for "escape" action confirmation)
	if editedCol >= 0 && editedCol < len(st.config.Columns) {
		col := st.config.Columns[editedCol]
		if col.OnViewData != nil && editedRow >= 0 && editedRow < len(st.data) {
			data := st.data[editedRow]
			col.OnViewData("escape", data, col.ID, editedRow, editedCol)
		}
	}

	st.state.editingRow = -1
	st.state.editingCol = -1
	st.editingEntry = nil
	st.state.editingValue = ""

	if st.table != nil {
		// Refresh the entire table to restore custom renderers
		st.table.Refresh()
		// Also refresh the specific cell to ensure custom renderer is applied
		if editedRow >= 0 && editedCol >= 0 {
			st.table.RefreshItem(widget.TableCellID{Row: editedRow + 1, Col: editedCol})
		}
	}

	// Restore focus and selection to the table
	st.SetSelectedCell(editedRow, editedCol)
	st.RequestFocus()
}

// extractFieldValue tries to extract a field value from data by column ID
func (st *Table) extractFieldValue(data interface{}, colID string) string {
	if data == nil {
		return ""
	}

	// Try reflection to get field by name (capitalize first letter for exported fields)
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", data)
	}

	// Try exact match first
	field := v.FieldByName(colID)
	if !field.IsValid() {
		// Try capitalized version (e.g., "name" -> "Name")
		caser := cases.Title(language.English)
		capitalized := caser.String(colID)
		field = v.FieldByName(capitalized)
	}
	if !field.IsValid() {
		// Try uppercase (e.g., "id" -> "ID")
		upper := strings.ToUpper(colID)
		field = v.FieldByName(upper)
	}
	if !field.IsValid() {
		// Try case-insensitive search through all fields
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if strings.EqualFold(t.Field(i).Name, colID) {
				field = v.Field(i)
				break
			}
		}
	}

	if field.IsValid() {
		return fmt.Sprintf("%v", field.Interface())
	}

	// Fallback to string representation
	return fmt.Sprintf("%v", data)
}

// findColumnDividerAtPosition detects if a position is near a column divider
func (st *Table) findColumnDividerAtPosition(pos fyne.Position) int {
	// Only auto-resize if in the header row area
	headerHeight := st.config.HeaderHeight
	if headerHeight == 0 {
		headerHeight = 30.0
	}

	if pos.Y > headerHeight {
		return -1 // Not in header area
	}

	// Sync config widths with actual table widths (in case of manual column resizing)
	// We need to read actual widths from table rendering since manual drag resizes
	// don't update our config
	st.syncColumnWidthsFromTable()

	// Calculate cumulative widths to find which divider was clicked
	// Only iterate through visible columns
	threshold := float32(10.0) // 10px tolerance on each side of divider for easier targeting
	xPos := float32(0)

	for _, actualIdx := range st.state.visibleColumns {
		col := st.config.Columns[actualIdx]
		colWidth := col.Width
		if colWidth == 0 {
			colWidth = 100 // Default width
		}
		xPos += colWidth

		distance := math.Abs(float64(pos.X - xPos))

		// Check if position is near this column's right edge
		if distance <= float64(threshold) {
			return actualIdx // Return actual column index
		}
	}

	return -1
}

// syncColumnWidthsFromTable syncs config widths with actual table widths
// This is needed because manual column resizing (dragging) updates the table's
// internal state but doesn't update our config
func (st *Table) syncColumnWidthsFromTable() {
	if st.table == nil {
		return
	}

	// Use reflection to access the internal columnWidths map in widget.Table
	tableValue := reflect.ValueOf(st.table).Elem()
	columnWidthsField := tableValue.FieldByName("columnWidths")

	if !columnWidthsField.IsValid() {
		st.logger().Warn("Cannot access table columnWidths field")
		return
	}

	// Use unsafe to access unexported field
	columnWidthsField = reflect.NewAt(columnWidthsField.Type(), unsafe.Pointer(columnWidthsField.UnsafeAddr())).Elem()

	// columnWidths is map[int]float32
	columnWidthsMap, ok := columnWidthsField.Interface().(map[int]float32)
	if !ok {
		st.logger().Warn("columnWidths field is not map[int]float32")
		return
	}

	// Update our config with actual widths from the table
	// Note: The table uses display column indices, so we need to map back to actual column indices
	for displayIdx, actualIdx := range st.state.visibleColumns {
		if actualWidth, exists := columnWidthsMap[displayIdx]; exists {
			if st.config.Columns[actualIdx].Width != actualWidth {
				st.config.Columns[actualIdx].Width = actualWidth
			}
		}
	}
}

// autoResizeColumn calculates and applies the optimal width for a column
func (st *Table) autoResizeColumn(colIndex int) {
	if colIndex < 0 || colIndex >= len(st.config.Columns) {
		st.logger().Error(fmt.Sprintf("Invalid column index: %d (total columns: %d)", colIndex, len(st.config.Columns)))
		return
	}

	col := st.config.Columns[colIndex]
	maxWidth := float32(0)

	// Measure header width
	headerText := col.Title
	if st.state.sortColumn == colIndex {
		headerText += " ▲" // Account for sort indicator
	}
	headerWidth := st.measureTextWidth(headerText, true)
	if headerWidth > maxWidth {
		maxWidth = headerWidth
	}

	// Measure all data cells in this column
	for i := range st.data {
		cellText := st.extractFieldValue(st.data[i], col.ID)
		cellWidth := st.measureTextWidth(cellText, false)

		// For columns with custom renderers (like hierarchical Task Name),
		// we need to account for additional visual elements
		if col.Renderer != nil && col.ID == "name" {
			// Check if this is hierarchical data with a Depth field
			v := reflect.ValueOf(st.data[i])
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() == reflect.Struct {
				depthField := v.FieldByName("Depth")
				if depthField.IsValid() && depthField.Kind() == reflect.Int {
					depth := int(depthField.Int())
					if depth > 0 {
						// Add indentation width based on config
						var indentWidth float32
						if st.config.ShowIndentation {
							indentPerLevel := st.config.IndentPerLevel
							if indentPerLevel == 0 {
								indentPerLevel = 20.0 // Default if not set
							}
							indentWidth = float32(5 + (depth * int(indentPerLevel)))
						}

						// Add tree icon width if icons are enabled
						var iconWidth float32
						if st.config.ShowIndentIcons {
							// Measure the actual icon text
							iconText := ""
							if st.config.ShowBranch {
								iconText = st.config.TreeIconTheme.Branch
							}
							if len(st.config.TreeIconTheme.Icons) > 0 {
								iconIndex := (depth - 1) % len(st.config.TreeIconTheme.Icons)
								iconText += st.config.TreeIconTheme.Icons[iconIndex]
							}
							iconWidth = st.measureTextWidth(iconText, false)
						}

						cellWidth += indentWidth + iconWidth
					}
				}
			}
		}

		if cellWidth > maxWidth {
			maxWidth = cellWidth
		}
	}

	// Add some padding
	maxWidth += 20 // 20px padding

	// Apply minimum width if configured
	if col.MinWidth > 0 && maxWidth < col.MinWidth {
		maxWidth = col.MinWidth
	}

	// Apply the new width
	st.config.Columns[colIndex].Width = maxWidth
	if st.table != nil {
		// Reapply ALL visible column widths to force Fyne to recalculate drag handler positions
		// This is necessary because Fyne's internal drag handlers aren't updated
		// when we only change one column width programmatically
		for displayIdx, actualIdx := range st.state.visibleColumns {
			c := st.config.Columns[actualIdx]
			if c.Width > 0 {
				st.table.SetColumnWidth(displayIdx, c.Width)
			}
		}

		// Force complete layout recalculation
		currentSize := st.table.Size()
		st.table.Resize(fyne.NewSize(currentSize.Width+0.1, currentSize.Height))
		st.table.Resize(currentSize)

		// Refresh both layers
		st.table.Refresh()
		st.table.Refresh()
	}

	// Save column widths if callback provided
	if st.config.SaveColumnWidths != nil {
		widths := make(map[string]float32)
		for idx, c := range st.config.Columns {
			widths[c.ID] = st.config.Columns[idx].Width
		}
		st.config.SaveColumnWidths(widths)
	}
}

// measureTextWidth estimates the width needed for text
func (st *Table) measureTextWidth(text string, bold bool) float32 {
	// Create a temporary text object to measure size
	tempText := canvas.NewText(text, theme.Color(theme.ColorNameForeground))
	if bold {
		tempText.TextStyle = fyne.TextStyle{Bold: true}
	}
	tempText.TextSize = theme.TextSize()
	return tempText.MinSize().Width
}

// showPopupMenu displays a transient popup menu for cell options activated by SPACE key
func (st *Table) showPopupMenu(rowIndex, colIndex int, options []string, dataItem interface{}) {
	col := st.config.Columns[colIndex]

	// Sync config widths with actual table widths (in case of manual column resizing)
	st.syncColumnWidthsFromTable()

	// Get current value to mark it in the menu
	var currentValue string
	if col.GetCellValue != nil {
		currentValue = col.GetCellValue(dataItem)
	}

	// Sort options alphabetically
	sortedOptions := make([]string, len(options))
	copy(sortedOptions, options)
	sort.Strings(sortedOptions)

	// Create popup menu items
	var menuItems []*fyne.MenuItem
	for _, option := range sortedOptions {
		optionVal := option // Capture for closure

		// Add checkmark to current value
		displayText := optionVal
		if optionVal == currentValue {
			displayText = "✓ " + optionVal
		}

		menuItems = append(menuItems, fyne.NewMenuItem(displayText, func() {
			st.logger().Info(fmt.Sprintf("[POPUP-CALLBACK] Menu item clicked: %s, params: rowIndex=%d colIndex=%d", optionVal, rowIndex, colIndex))
			st.logger().Info(fmt.Sprintf("[POPUP-CALLBACK] Current state BEFORE callback: selectedRow=%d selectedCol=%d", st.state.selectedRow, st.state.selectedCol))

			// Call the callback to get new value
			newValue := col.OnPopupSelected(dataItem, optionVal, rowIndex)

			st.logger().Info(fmt.Sprintf("[POPUP-CALLBACK] After OnPopupSelected, state: selectedRow=%d selectedCol=%d", st.state.selectedRow, st.state.selectedCol))

			// Trigger OnCellEdited callback if defined
			if st.config.OnCellEdited != nil {
				st.config.OnCellEdited(rowIndex, col.ID, newValue, dataItem)
			}

			st.logger().Info(fmt.Sprintf("[POPUP-CALLBACK] After OnCellEdited, state: selectedRow=%d selectedCol=%d", st.state.selectedRow, st.state.selectedCol))

			// CRITICAL: Restore navigation state to the cell where popup was triggered
			// The popup dismissal or callbacks may have triggered OnSelected which changed our state
			// But we want arrow keys to continue from where the popup was originally shown
			st.logger().Info(fmt.Sprintf("[POPUP-CALLBACK] Restoring navigation state from (%d,%d) to (%d,%d)",
				st.state.selectedRow, st.state.selectedCol, rowIndex, colIndex))
			st.state.selectedRow = rowIndex
			st.state.selectedCol = colIndex

			// Refresh the entire table to show updated value
			// Note: We refresh the whole table instead of just the cell to avoid
			// unintended side effects from RefreshItem (which can cause selection changes)
			// We DON'T call Select() because that would trigger OnSelected callbacks
			// and potentially change our navigation state. Our selectedRow/selectedCol
			// remain unchanged, so arrow keys will work correctly from the current position.
			if st.table != nil {
				st.logger().Info(fmt.Sprintf("[REFRESH] Refreshing table after dropdown change, navigation state: row=%d col=%d", st.state.selectedRow, st.state.selectedCol))
				st.table.Refresh()
				st.logger().Info(fmt.Sprintf("[REFRESH] Table refreshed, navigation state preserved: row=%d col=%d", st.state.selectedRow, st.state.selectedCol))
			}

			st.logger().Info(fmt.Sprintf("Popup selection: %s for row %d, col %s", optionVal, rowIndex, col.ID))
		}))
	}

	// Create and show popup menu
	popupMenu := fyne.NewMenu("", menuItems...)

	// Get canvas for positioning
	canvas := fyne.CurrentApp().Driver().CanvasForObject(st.table)
	if canvas == nil {
		st.logger().Warn("Cannot show popup menu: no canvas found")
		return
	}

	// Calculate popup position at current cell
	if st.table == nil {
		st.logger().Warn("Cannot show popup menu: table not initialized")
		return
	}

	// Find display column index
	displayColIndex := -1
	for i, visCol := range st.state.visibleColumns {
		if visCol == colIndex {
			displayColIndex = i
			break
		}
	}
	if displayColIndex < 0 {
		st.logger().Warn(fmt.Sprintf("Cannot show popup menu: column %d not visible", colIndex))
		return
	}

	// Calculate cell position
	// X position: sum of all column widths before this column
	xPos := float32(0)
	for i := 0; i < displayColIndex; i++ {
		actualColIdx := st.state.visibleColumns[i]
		if actualColIdx < len(st.config.Columns) {
			xPos += st.config.Columns[actualColIdx].Width
		}
	}

	// Calculate Y position: header height + (row * row height)
	yPos := st.config.HeaderHeight + (float32(rowIndex) * st.config.RowHeight)

	// Get table's absolute position on canvas for correct popup placement
	tablePos := fyne.CurrentApp().Driver().AbsolutePositionForObject(st.table)

	// Position popup at bottom of the cell (yPos is top of cell, add rowHeight for bottom)
	popupYPos := tablePos.Y + yPos + st.config.RowHeight
	popupXPos := tablePos.X + xPos
	pos := fyne.NewPos(popupXPos, popupYPos)

	st.logger().Info(fmt.Sprintf("Showing popup at pos (%.1f, %.1f) for cell row=%d col=%d (tablePos=%v, cellTop=%.1f)", popupXPos, popupYPos, rowIndex, colIndex, tablePos, yPos))

	widget.ShowPopUpMenuAtPosition(popupMenu, canvas, pos)
}

// RenderTextWithPopupIcon renders a text cell with an optional dropdown indicator icon
// This is a helper function for columns that have popup menus
func RenderTextWithPopupIcon(text string, showIcon bool, alignment TextAlignment) *fyne.Container {
	textLabel := widget.NewLabel(text)

	if !showIcon {
		// No icon, just return centered text based on alignment
		switch alignment {
		case AlignCenter:
			return container.NewCenter(textLabel)
		case AlignRight:
			return container.NewHBox(widget.NewLabel(""), textLabel) // Right align
		default: // AlignLeft
			return container.NewHBox(textLabel, widget.NewLabel("")) // Left align with spacer
		}
	}

	// Create dropdown icon (small triangle)
	icon := widget.NewIcon(theme.MenuDropDownIcon())

	// Combine text and icon based on alignment
	switch alignment {
	case AlignCenter:
		return container.NewCenter(
			container.NewHBox(textLabel, icon),
		)
	case AlignRight:
		return container.NewHBox(widget.NewLabel(""), textLabel, icon)
	default: // AlignLeft
		return container.NewHBox(textLabel, icon, widget.NewLabel(""))
	}
}
