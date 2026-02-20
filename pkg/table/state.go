package table

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
)

// TableState holds all runtime state for a Table widget
// This includes view state, sorting, filtering, selection, and editing state
type TableState struct {
	// View state
	visibleColumns []int // Indices of visible columns (filters out hidden columns)
	visibleRows    []int // Indices of visible rows (after tree filtering)

	// Sort state
	sortColumn int  // -1 = no sort, otherwise column index
	sortAsc    bool // true = ascending, false = descending

	// Focus state
	hasFocus bool // true when table has keyboard focus

	// Filter state
	filterText          string // Current filter text
	filterRegex         bool   // true = use regex matching, false = plain text
	filterCaseSensitive bool   // true = case-sensitive matching, false = case-insensitive

	// Selection state
	selectedRow  int          // -1 = no selection, otherwise data row index (single-select mode)
	selectedCol  int          // -1 = no selection, otherwise actual column index (only used when RowSelectOnlyMode=false)
	selectedRows map[int]bool // Multi-select mode: map of selected data row indices

	// Edit state (row/col/value only - widget reference stays in Table)
	editingRow   int    // -1 = not editing
	editingCol   int    // -1 = not editing
	editingValue string // Original value before edit

	// Navigation state
	isKeyboardNavigation bool // true when navigating with arrow keys (don't auto-activate)
	isReselecting        bool // true when re-selecting cell after refresh (don't fire callbacks)
}

// NewTableState creates a new TableState with default values
func NewTableState() *TableState {
	return &TableState{
		visibleColumns: []int{},
		visibleRows:    []int{},
		sortColumn:     -1,
		sortAsc:        true,
		hasFocus:       false,
		selectedRow:    -1,
		selectedCol:    -1,
		selectedRows:   make(map[int]bool),
		editingRow:     -1,
		editingCol:     -1,
	}
}

// ========================================
// Selection State Methods
// ========================================

// HasSelection returns true if at least one row is selected
func (s *TableState) HasSelection() bool {
	return s.selectedRow >= 0 || len(s.selectedRows) > 0
}

// ClearSelection clears the current selection (both single and multi-select)
func (s *TableState) ClearSelection() {
	s.selectedRow = -1
	s.selectedCol = -1
	s.selectedRows = make(map[int]bool)
}

// GetSelectedRows returns a slice of selected row indices (works for both single and multi-select)
func (s *TableState) GetSelectedRows() []int {
	if len(s.selectedRows) > 0 {
		// Multi-select mode
		rows := make([]int, 0, len(s.selectedRows))
		for row := range s.selectedRows {
			rows = append(rows, row)
		}
		return rows
	}
	// Single-select mode
	if s.selectedRow >= 0 {
		return []int{s.selectedRow}
	}
	return []int{}
}

// SetSelectedRows sets the selected rows (multi-select mode)
func (s *TableState) SetSelectedRows(rows []int) {
	s.selectedRows = make(map[int]bool)
	for _, row := range rows {
		if row >= 0 {
			s.selectedRows[row] = true
		}
	}
	// Clear single-select state
	s.selectedRow = -1
}

// AddSelectedRow adds a row to the selection (multi-select mode)
func (s *TableState) AddSelectedRow(row int) {
	if row >= 0 {
		s.selectedRows[row] = true
		s.selectedRow = -1 // Clear single-select
	}
}

// RemoveSelectedRow removes a row from the selection (multi-select mode)
func (s *TableState) RemoveSelectedRow(row int) {
	delete(s.selectedRows, row)
}

// IsRowSelected returns true if the given row is selected
func (s *TableState) IsRowSelected(row int) bool {
	if len(s.selectedRows) > 0 {
		return s.selectedRows[row]
	}
	return s.selectedRow == row
}

// GetSelectionCount returns the number of selected rows
func (s *TableState) GetSelectionCount() int {
	if len(s.selectedRows) > 0 {
		return len(s.selectedRows)
	}
	if s.selectedRow >= 0 {
		return 1
	}
	return 0
}

// ========================================
// Edit State Methods
// ========================================

// IsEditing returns true if a cell is currently being edited
func (s *TableState) IsEditing() bool {
	return s.editingRow >= 0 && s.editingCol >= 0
}

// ClearEdit clears the edit state
func (s *TableState) ClearEdit() {
	s.editingRow = -1
	s.editingCol = -1
	s.editingValue = ""
}

// ========================================
// Filter State Methods
// ========================================

// HasFilter returns true if a filter is active
func (s *TableState) HasFilter() bool {
	return s.filterText != ""
}

// ClearFilter clears the filter state
func (s *TableState) ClearFilter() {
	s.filterText = ""
	s.filterRegex = false
	s.filterCaseSensitive = false
}

// ========================================
// Sort State Methods
// ========================================

// IsSorted returns true if table is sorted
func (s *TableState) IsSorted() bool {
	return s.sortColumn >= 0
}

// ClearSort clears the sort state
func (s *TableState) ClearSort() {
	s.sortColumn = -1
	s.sortAsc = true
}

// ========================================
// View State Methods
// ========================================

// VisibleRowCount returns the number of visible rows
func (s *TableState) VisibleRowCount() int {
	return len(s.visibleRows)
}

// VisibleColumnCount returns the number of visible columns
func (s *TableState) VisibleColumnCount() int {
	return len(s.visibleColumns)
}

// ========================================
// Reset Methods
// ========================================

// Reset resets all state to initial values
func (s *TableState) Reset() {
	s.visibleColumns = []int{}
	s.visibleRows = []int{}
	s.sortColumn = -1
	s.sortAsc = true
	s.hasFocus = false
	s.filterText = ""
	s.filterRegex = false
	s.filterCaseSensitive = false
	s.selectedRow = -1
	s.selectedCol = -1
	s.selectedRows = make(map[int]bool)
	s.editingRow = -1
	s.editingCol = -1
	s.editingValue = ""
}

// ========================================
// Debug Methods
// ========================================

// String returns a string representation of the state for debugging
func (s *TableState) String() string {
	return fmt.Sprintf("TableState{sort:%d/%v, sel:%d/%d, edit:%d/%d, filter:%q, focus:%v, rows:%d/%d cols:%d}",
		s.sortColumn, s.sortAsc,
		s.selectedRow, s.selectedCol,
		s.editingRow, s.editingCol,
		s.filterText, s.hasFocus,
		len(s.visibleRows), len(s.visibleRows),
		len(s.visibleColumns))
}

// ========================================
// Getter Methods (for external access)
// These return copies to prevent external modification
// ========================================

// GetVisibleColumns returns a copy of visible column indices
func (s *TableState) GetVisibleColumns() []int {
	result := make([]int, len(s.visibleColumns))
	copy(result, s.visibleColumns)
	return result
}

// GetVisibleRows returns a copy of visible row indices
func (s *TableState) GetVisibleRows() []int {
	result := make([]int, len(s.visibleRows))
	copy(result, s.visibleRows)
	return result
}

// GetSortColumn returns the sort column index (-1 if not sorted)
func (s *TableState) GetSortColumn() int {
	return s.sortColumn
}

// GetSortAscending returns true if sorting ascending
func (s *TableState) GetSortAscending() bool {
	return s.sortAsc
}

// GetFilterText returns the current filter text
func (s *TableState) GetFilterText() string {
	return s.filterText
}

// IsFilterRegex returns true if regex filtering is enabled
func (s *TableState) IsFilterRegex() bool {
	return s.filterRegex
}

// IsFilterCaseSensitive returns true if case-sensitive filtering is enabled
func (s *TableState) IsFilterCaseSensitive() bool {
	return s.filterCaseSensitive
}

// GetSelectedRow returns the selected row index (-1 if none)
func (s *TableState) GetSelectedRow() int {
	return s.selectedRow
}

// GetSelectedCol returns the selected column index (-1 if none)
func (s *TableState) GetSelectedCol() int {
	return s.selectedCol
}

// GetEditingRow returns the editing row index (-1 if not editing)
func (s *TableState) GetEditingRow() int {
	return s.editingRow
}

// GetEditingCol returns the editing column index (-1 if not editing)
func (s *TableState) GetEditingCol() int {
	return s.editingCol
}

// GetEditingValue returns the original value being edited
func (s *TableState) GetEditingValue() string {
	return s.editingValue
}

// HasFocus returns true if table has keyboard focus
func (s *TableState) HasFocus() bool {
	return s.hasFocus
}

// ========================================
// TableStateSnapshot for save/restore
// ========================================

// TableStateSnapshot represents a snapshot of table state that can be saved and restored
type TableStateSnapshot struct {
	SortColumn          int
	SortAsc             bool
	FilterText          string
	FilterRegex         bool
	FilterCaseSensitive bool
	SelectedRow         int
	SelectedCol         int
	SelectedRows        map[int]bool // Multi-select state
}

// Snapshot creates a snapshot of the current state
func (s *TableState) Snapshot() *TableStateSnapshot {
	// Deep copy selectedRows map
	selectedRowsCopy := make(map[int]bool)
	for row, val := range s.selectedRows {
		selectedRowsCopy[row] = val
	}

	return &TableStateSnapshot{
		SortColumn:          s.sortColumn,
		SortAsc:             s.sortAsc,
		FilterText:          s.filterText,
		FilterRegex:         s.filterRegex,
		FilterCaseSensitive: s.filterCaseSensitive,
		SelectedRow:         s.selectedRow,
		SelectedCol:         s.selectedCol,
		SelectedRows:        selectedRowsCopy,
	}
}

// RestoreFromSnapshot restores state from a snapshot
func (s *TableState) RestoreFromSnapshot(snap *TableStateSnapshot) {
	s.sortColumn = snap.SortColumn
	s.sortAsc = snap.SortAsc
	s.filterText = snap.FilterText
	s.filterRegex = snap.FilterRegex
	s.filterCaseSensitive = snap.FilterCaseSensitive
	s.selectedRow = snap.SelectedRow
	s.selectedCol = snap.SelectedCol

	// Deep copy selectedRows map
	s.selectedRows = make(map[int]bool)
	if snap.SelectedRows != nil {
		for row, val := range snap.SelectedRows {
			s.selectedRows[row] = val
		}
	}
}

// ========================================
// Internal helper to check if entry is still valid
// (Used by Table to verify editingEntry is not stale)
// ========================================

// IsEditEntryValid checks if the given entry widget matches current edit state
func (s *TableState) IsEditEntryValid(entry *widget.Entry) bool {
	return s.IsEditing() && entry != nil
}
