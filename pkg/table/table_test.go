package table

import (
	"testing"

	"fyne.io/fyne/v2/widget"
)

// TestSortData tests the sorting logic with custom comparators
func TestSortData(t *testing.T) {
	tests := []struct {
		name       string
		data       []interface{}
		sortColumn int
		sortAsc    bool
		comparator SortComparator
		expected   []interface{}
	}{
		{
			name:       "Sort strings ascending",
			data:       []interface{}{"zebra", "apple", "mango"},
			sortColumn: 0,
			sortAsc:    true,
			comparator: nil,
			expected:   []interface{}{"apple", "mango", "zebra"},
		},
		{
			name:       "Sort strings descending",
			data:       []interface{}{"zebra", "apple", "mango"},
			sortColumn: 0,
			sortAsc:    false,
			comparator: nil,
			expected:   []interface{}{"zebra", "mango", "apple"},
		},
		{
			name:       "Sort numbers with custom comparator ascending",
			data:       []interface{}{30, 10, 20},
			sortColumn: 0,
			sortAsc:    true,
			comparator: func(a, b interface{}) int {
				aVal := a.(int)
				bVal := b.(int)
				return aVal - bVal
			},
			expected: []interface{}{10, 20, 30},
		},
		{
			name:       "Sort numbers with custom comparator descending",
			data:       []interface{}{30, 10, 20},
			sortColumn: 0,
			sortAsc:    false,
			comparator: func(a, b interface{}) int {
				aVal := a.(int)
				bVal := b.(int)
				return aVal - bVal
			},
			expected: []interface{}{30, 20, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig("test")
			config.Columns = []ColumnConfig{
				{
					ID:         "col1",
					Title:      "Column 1",
					Sortable:   true,
					Comparator: tt.comparator,
				},
			}

			table := &Table{
				config: config,
				data:   tt.data,
				state: &TableState{
					sortColumn: tt.sortColumn,
					sortAsc:    tt.sortAsc,
				},
			}

			table.sortData()

			for i, expected := range tt.expected {
				if table.data[i] != expected {
					t.Errorf("Expected data[%d] = %v, got %v", i, expected, table.data[i])
				}
			}
		})
	}
}

// TestHandleHeaderClick tests clicking column headers to toggle sort
func TestHandleHeaderClick(t *testing.T) {
	config := NewConfig("test")
	config.Columns = []ColumnConfig{
		{ID: "col1", Title: "Column 1", Sortable: true},
		{ID: "col2", Title: "Column 2", Sortable: false},
	}

	table := &Table{
		config:       config,
		data:         []interface{}{"zebra", "apple", "mango"},
		MouseHandler: NewDefaultMouseHandler(),
		state: &TableState{
			sortColumn:     -1,
			sortAsc:        true,
			visibleColumns: []int{0, 1}, // Both columns visible
		},
	}
	// Need to create a minimal table widget for refresh
	table.table = &keyboardForwardingTable{Table: &widget.Table{}}

	table.handleHeaderClick(0)
	if table.state.sortColumn != 0 {
		t.Errorf("Expected sortColumn = 0, got %d", table.state.sortColumn)
	}
	if !table.state.sortAsc {
		t.Error("Expected sortAsc = true after first click")
	}
	if table.data[0] != "apple" {
		t.Errorf("Expected data[0] = apple after sort, got %v", table.data[0])
	}

	table.handleHeaderClick(0)
	if table.state.sortAsc {
		t.Error("Expected sortAsc toggled to false")
	}
	if table.data[0] != "zebra" {
		t.Errorf("Expected data[0] = zebra after descending sort, got %v", table.data[0])
	}

	table.handleHeaderClick(1)
	if table.state.sortColumn != 0 {
		t.Errorf("Expected sortColumn still = 0, got %d", table.state.sortColumn)
	}
}

// TestSortDataWithStructs tests sorting complex data structures
func TestSortDataWithStructs(t *testing.T) {
	type Task struct {
		ID       int
		Name     string
		Priority int
	}

	data := []interface{}{
		Task{ID: 3, Name: "Task C", Priority: 2},
		Task{ID: 1, Name: "Task A", Priority: 1},
		Task{ID: 2, Name: "Task B", Priority: 3},
	}

	config := NewConfig("test")
	config.Columns = []ColumnConfig{
		{
			ID:       "priority",
			Title:    "Priority",
			Sortable: true,
			Comparator: func(a, b interface{}) int {
				taskA := a.(Task)
				taskB := b.(Task)
				return taskA.Priority - taskB.Priority
			},
		},
	}

	table := &Table{
		config: config,
		data:   data,
		state: &TableState{
			sortColumn: 0,
			sortAsc:    true,
		},
	}

	table.sortData()

	expectedPriorities := []int{1, 2, 3}
	for i, expected := range expectedPriorities {
		task := table.data[i].(Task)
		if task.Priority != expected {
			t.Errorf("Expected data[%d].Priority = %d, got %d", i, expected, task.Priority)
		}
	}
}

// TestSortEmptyData tests sorting with no data
func TestSortEmptyData(t *testing.T) {
	config := NewConfig("test")
	config.Columns = []ColumnConfig{
		{ID: "col1", Title: "Column 1", Sortable: true},
	}

	table := &Table{
		config: config,
		data:   []interface{}{},
		state: &TableState{
			sortColumn: 0,
			sortAsc:    true,
		},
	}

	table.sortData()

	if len(table.data) != 0 {
		t.Errorf("Expected data length = 0, got %d", len(table.data))
	}
}

// TestStartEdit tests starting an edit operation
func TestStartEdit(t *testing.T) {
	config := NewConfig("test")
	config.Columns = []ColumnConfig{
		{ID: "col1", Title: "Column 1", Editable: true},
	}

	table := &Table{
		config: config,
		data:   []interface{}{"apple", "banana"},
		state: &TableState{
			editingRow:     -1,
			editingCol:     -1,
			visibleColumns: []int{0}, // Column 0 visible
		},
	}

	table.startEdit(0, 0)

	if table.state.editingRow != 0 {
		t.Errorf("Expected editingRow = 0, got %d", table.state.editingRow)
	}
	if table.state.editingCol != 0 {
		t.Errorf("Expected editingCol = 0, got %d", table.state.editingCol)
	}
	if table.state.editingValue != "apple" {
		t.Errorf("Expected editingValue = apple, got %s", table.state.editingValue)
	}
}

// TestCancelEdit tests canceling an edit operation
func TestCancelEdit(t *testing.T) {
	config := NewConfig("test")
	config.Columns = []ColumnConfig{
		{ID: "col1", Title: "Column 1", Editable: true},
	}

	table := &Table{
		config: config,
		data:   []interface{}{"apple"},
		state: &TableState{
			editingRow:   0,
			editingCol:   0,
			editingValue: "apple",
		},
	}

	table.cancelEdit()

	if table.state.editingRow != -1 {
		t.Errorf("Expected editingRow = -1, got %d", table.state.editingRow)
	}
	if table.state.editingCol != -1 {
		t.Errorf("Expected editingCol = -1, got %d", table.state.editingCol)
	}
}

// TestHandleCellClick tests cell click handling for selection
func TestHandleCellClick(t *testing.T) {
	config := NewConfig("test")
	config.Columns = []ColumnConfig{
		{ID: "col1", Title: "Column 1", Editable: true},
		{ID: "col2", Title: "Column 2", Editable: false},
	}

	table := &Table{
		config:       config,
		data:         []interface{}{"apple"},
		MouseHandler: NewDefaultMouseHandler(),
		state: &TableState{
			selectedRow:    -1,
			selectedCol:    -1,
			visibleColumns: []int{0, 1}, // Both columns visible
			visibleRows:    []int{0},    // One data row
		},
	}
	// Need table widget for refresh
	table.table = &keyboardForwardingTable{Table: &widget.Table{}}

	// Click on first column cell - should select row 0, col 0
	table.handleCellClick(widget.TableCellID{Row: 1, Col: 0})
	if table.state.selectedRow != 0 {
		t.Errorf("Expected selectedRow = 0, got %d", table.state.selectedRow)
	}
	if table.state.selectedCol != 0 {
		t.Errorf("Expected selectedCol = 0, got %d", table.state.selectedCol)
	}

	// Click on second column cell - should select row 0, col 1
	table.handleCellClick(widget.TableCellID{Row: 1, Col: 1})
	if table.state.selectedRow != 0 {
		t.Errorf("Expected selectedRow = 0, got %d", table.state.selectedRow)
	}
	if table.state.selectedCol != 1 {
		t.Errorf("Expected selectedCol = 1, got %d", table.state.selectedCol)
	}

	// Click on header row should not change selection
	prevRow := table.state.selectedRow
	table.handleCellClick(widget.TableCellID{Row: 0, Col: 0})
	if table.state.selectedRow != prevRow {
		t.Errorf("Expected selectedRow unchanged at %d, got %d", prevRow, table.state.selectedRow)
	}
}
