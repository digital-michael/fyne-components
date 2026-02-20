package table

import (
	"testing"
)

func TestNewTableState(t *testing.T) {
	state := NewTableState()

	if state.sortColumn != -1 {
		t.Errorf("Expected sortColumn = -1, got %d", state.sortColumn)
	}
	if !state.sortAsc {
		t.Error("Expected sortAsc = true")
	}
	if state.selectedRow != -1 {
		t.Errorf("Expected selectedRow = -1, got %d", state.selectedRow)
	}
	if state.selectedCol != -1 {
		t.Errorf("Expected selectedCol = -1, got %d", state.selectedCol)
	}
	if state.editingRow != -1 {
		t.Errorf("Expected editingRow = -1, got %d", state.editingRow)
	}
	if state.editingCol != -1 {
		t.Errorf("Expected editingCol = -1, got %d", state.editingCol)
	}
	if state.hasFocus {
		t.Error("Expected hasFocus = false")
	}
}

func TestTableState_Selection(t *testing.T) {
	state := NewTableState()

	// Initially no selection
	if state.HasSelection() {
		t.Error("Expected no selection initially")
	}

	// Set selection
	state.selectedRow = 5
	state.selectedCol = 2
	if !state.HasSelection() {
		t.Error("Expected selection after setting selectedRow")
	}
	if state.GetSelectedRow() != 5 {
		t.Errorf("Expected selectedRow = 5, got %d", state.GetSelectedRow())
	}
	if state.GetSelectedCol() != 2 {
		t.Errorf("Expected selectedCol = 2, got %d", state.GetSelectedCol())
	}

	// Clear selection
	state.ClearSelection()
	if state.HasSelection() {
		t.Error("Expected no selection after clear")
	}
	if state.selectedRow != -1 {
		t.Errorf("Expected selectedRow = -1 after clear, got %d", state.selectedRow)
	}
	if state.selectedCol != -1 {
		t.Errorf("Expected selectedCol = -1 after clear, got %d", state.selectedCol)
	}
}

func TestTableState_Edit(t *testing.T) {
	state := NewTableState()

	// Initially not editing
	if state.IsEditing() {
		t.Error("Expected not editing initially")
	}

	// Set edit state
	state.editingRow = 3
	state.editingCol = 1
	state.editingValue = "original"
	if !state.IsEditing() {
		t.Error("Expected editing after setting editingRow and editingCol")
	}
	if state.GetEditingRow() != 3 {
		t.Errorf("Expected editingRow = 3, got %d", state.GetEditingRow())
	}
	if state.GetEditingCol() != 1 {
		t.Errorf("Expected editingCol = 1, got %d", state.GetEditingCol())
	}
	if state.GetEditingValue() != "original" {
		t.Errorf("Expected editingValue = 'original', got %q", state.GetEditingValue())
	}

	// Clear edit state
	state.ClearEdit()
	if state.IsEditing() {
		t.Error("Expected not editing after clear")
	}
	if state.editingRow != -1 {
		t.Errorf("Expected editingRow = -1 after clear, got %d", state.editingRow)
	}
	if state.editingCol != -1 {
		t.Errorf("Expected editingCol = -1 after clear, got %d", state.editingCol)
	}
	if state.editingValue != "" {
		t.Errorf("Expected editingValue = '' after clear, got %q", state.editingValue)
	}
}

func TestTableState_Filter(t *testing.T) {
	state := NewTableState()

	// Initially no filter
	if state.HasFilter() {
		t.Error("Expected no filter initially")
	}

	// Set filter
	state.filterText = "test"
	state.filterRegex = true
	state.filterCaseSensitive = true
	if !state.HasFilter() {
		t.Error("Expected filter after setting filterText")
	}
	if state.GetFilterText() != "test" {
		t.Errorf("Expected filterText = 'test', got %q", state.GetFilterText())
	}
	if !state.IsFilterRegex() {
		t.Error("Expected filterRegex = true")
	}
	if !state.IsFilterCaseSensitive() {
		t.Error("Expected filterCaseSensitive = true")
	}

	// Clear filter
	state.ClearFilter()
	if state.HasFilter() {
		t.Error("Expected no filter after clear")
	}
	if state.filterText != "" {
		t.Errorf("Expected filterText = '' after clear, got %q", state.filterText)
	}
	if state.filterRegex {
		t.Error("Expected filterRegex = false after clear")
	}
	if state.filterCaseSensitive {
		t.Error("Expected filterCaseSensitive = false after clear")
	}
}

func TestTableState_Sort(t *testing.T) {
	state := NewTableState()

	// Initially not sorted
	if state.IsSorted() {
		t.Error("Expected not sorted initially")
	}

	// Set sort
	state.sortColumn = 2
	state.sortAsc = false
	if !state.IsSorted() {
		t.Error("Expected sorted after setting sortColumn")
	}
	if state.GetSortColumn() != 2 {
		t.Errorf("Expected sortColumn = 2, got %d", state.GetSortColumn())
	}
	if state.GetSortAscending() {
		t.Error("Expected sortAsc = false")
	}

	// Clear sort
	state.ClearSort()
	if state.IsSorted() {
		t.Error("Expected not sorted after clear")
	}
	if state.sortColumn != -1 {
		t.Errorf("Expected sortColumn = -1 after clear, got %d", state.sortColumn)
	}
	if !state.sortAsc {
		t.Error("Expected sortAsc = true after clear")
	}
}

func TestTableState_ViewState(t *testing.T) {
	state := NewTableState()

	// Set visible rows/columns
	state.visibleRows = []int{0, 1, 2, 5, 8}
	state.visibleColumns = []int{0, 2, 3}

	if state.VisibleRowCount() != 5 {
		t.Errorf("Expected VisibleRowCount = 5, got %d", state.VisibleRowCount())
	}
	if state.VisibleColumnCount() != 3 {
		t.Errorf("Expected VisibleColumnCount = 3, got %d", state.VisibleColumnCount())
	}

	// Test getters return copies
	rows := state.GetVisibleRows()
	cols := state.GetVisibleColumns()

	// Modify copies
	rows[0] = 999
	cols[0] = 888

	// Original should be unchanged
	if state.visibleRows[0] != 0 {
		t.Error("GetVisibleRows should return a copy, not a reference")
	}
	if state.visibleColumns[0] != 0 {
		t.Error("GetVisibleColumns should return a copy, not a reference")
	}
}

func TestTableState_Focus(t *testing.T) {
	state := NewTableState()

	if state.HasFocus() {
		t.Error("Expected no focus initially")
	}

	state.hasFocus = true
	if !state.HasFocus() {
		t.Error("Expected focus after setting hasFocus = true")
	}
}

func TestTableState_Reset(t *testing.T) {
	state := NewTableState()

	// Set various state
	state.sortColumn = 3
	state.sortAsc = false
	state.selectedRow = 5
	state.selectedCol = 2
	state.editingRow = 1
	state.editingCol = 1
	state.editingValue = "test"
	state.filterText = "filter"
	state.filterRegex = true
	state.filterCaseSensitive = true
	state.hasFocus = true
	state.visibleRows = []int{1, 2, 3}
	state.visibleColumns = []int{0, 1}

	// Reset
	state.Reset()

	// Verify all state is reset
	if state.sortColumn != -1 {
		t.Errorf("Expected sortColumn = -1 after reset, got %d", state.sortColumn)
	}
	if !state.sortAsc {
		t.Error("Expected sortAsc = true after reset")
	}
	if state.selectedRow != -1 {
		t.Errorf("Expected selectedRow = -1 after reset, got %d", state.selectedRow)
	}
	if state.selectedCol != -1 {
		t.Errorf("Expected selectedCol = -1 after reset, got %d", state.selectedCol)
	}
	if state.editingRow != -1 {
		t.Errorf("Expected editingRow = -1 after reset, got %d", state.editingRow)
	}
	if state.editingCol != -1 {
		t.Errorf("Expected editingCol = -1 after reset, got %d", state.editingCol)
	}
	if state.editingValue != "" {
		t.Errorf("Expected editingValue = '' after reset, got %q", state.editingValue)
	}
	if state.filterText != "" {
		t.Errorf("Expected filterText = '' after reset, got %q", state.filterText)
	}
	if state.filterRegex {
		t.Error("Expected filterRegex = false after reset")
	}
	if state.filterCaseSensitive {
		t.Error("Expected filterCaseSensitive = false after reset")
	}
	if state.hasFocus {
		t.Error("Expected hasFocus = false after reset")
	}
	if len(state.visibleRows) != 0 {
		t.Errorf("Expected empty visibleRows after reset, got %v", state.visibleRows)
	}
	if len(state.visibleColumns) != 0 {
		t.Errorf("Expected empty visibleColumns after reset, got %v", state.visibleColumns)
	}
}

func TestTableState_Snapshot(t *testing.T) {
	state := NewTableState()

	// Set state
	state.sortColumn = 2
	state.sortAsc = false
	state.filterText = "test"
	state.filterRegex = true
	state.filterCaseSensitive = true
	state.selectedRow = 5
	state.selectedCol = 3

	// Create snapshot
	snap := state.Snapshot()

	// Modify state
	state.sortColumn = 0
	state.sortAsc = true
	state.filterText = "modified"
	state.selectedRow = 10

	// Restore from snapshot
	state.RestoreFromSnapshot(snap)

	// Verify restoration
	if state.sortColumn != 2 {
		t.Errorf("Expected sortColumn = 2 after restore, got %d", state.sortColumn)
	}
	if state.sortAsc {
		t.Error("Expected sortAsc = false after restore")
	}
	if state.filterText != "test" {
		t.Errorf("Expected filterText = 'test' after restore, got %q", state.filterText)
	}
	if !state.filterRegex {
		t.Error("Expected filterRegex = true after restore")
	}
	if !state.filterCaseSensitive {
		t.Error("Expected filterCaseSensitive = true after restore")
	}
	if state.selectedRow != 5 {
		t.Errorf("Expected selectedRow = 5 after restore, got %d", state.selectedRow)
	}
	if state.selectedCol != 3 {
		t.Errorf("Expected selectedCol = 3 after restore, got %d", state.selectedCol)
	}
}

func TestTableState_String(t *testing.T) {
	state := NewTableState()
	state.sortColumn = 2
	state.selectedRow = 5
	state.filterText = "test"

	// Just verify it doesn't panic
	str := state.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}
