package table

import (
	"fmt"
	"strings"
	"testing"
)

// TestData represents a simple test struct with various fields
type TestData struct {
	ID       int
	Name     string
	Status   string
	Priority int
	Active   bool
}

// Helper function to create a test config with standard columns
func createTestConfig() *Config {
	config := NewConfig("test-table")
	config.Columns = []ColumnConfig{
		{
			ID:    "id",
			Title: "ID",
			Width: 50,
			GetCellValue: func(data interface{}) string {
				return fmt.Sprintf("%d", data.(TestData).ID)
			},
		},
		{
			ID:    "name",
			Title: "Name",
			Width: 150,
			GetCellValue: func(data interface{}) string {
				return data.(TestData).Name
			},
		},
		{
			ID:    "status",
			Title: "Status",
			Width: 100,
			GetCellValue: func(data interface{}) string {
				return data.(TestData).Status
			},
		},
		{
			ID:    "priority",
			Title: "Priority",
			Width: 80,
			GetCellValue: func(data interface{}) string {
				return fmt.Sprintf("%d", data.(TestData).Priority)
			},
		},
	}
	// Enable filtering on all columns for test purposes
	config.FilterColumns = []string{"id", "name", "status", "priority"}
	return config
}

// Helper function to create test data
func createTestData() []interface{} {
	return []interface{}{
		TestData{ID: 1, Name: "Alice", Status: "Active", Priority: 1, Active: true},
		TestData{ID: 2, Name: "Bob", Status: "Inactive", Priority: 3, Active: false},
		TestData{ID: 3, Name: "Charlie", Status: "Active", Priority: 2, Active: true},
		TestData{ID: 4, Name: "alice", Status: "Pending", Priority: 1, Active: false},
		TestData{ID: 5, Name: "David", Status: "Active", Priority: 4, Active: true},
	}
}

// Helper function to create a minimal test table without Fyne widget initialization
func createTestTable(config *Config) *Table {
	table := &Table{
		config:       config,
		data:         []interface{}{},
		state:        NewTableState(),
		KeyHandler:   NewDefaultKeyHandler(),
		MouseHandler: NewDefaultMouseHandler(),
		FocusHandler: NewDefaultFocusHandler(),
	}

	// Build visible columns list (needed for many operations)
	table.RebuildVisibleColumns()

	return table
}

// TestLogger is a simple logger for testing
type TestLogger struct {
	logs []string
}

func (l *TestLogger) Debug(msg string, keyvals ...interface{}) {
	l.logs = append(l.logs, "DEBUG: "+msg)
}

func (l *TestLogger) Info(msg string, keyvals ...interface{}) {
	l.logs = append(l.logs, "INFO: "+msg)
}

func (l *TestLogger) Warn(msg string, keyvals ...interface{}) {
	l.logs = append(l.logs, "WARN: "+msg)
}

func (l *TestLogger) Error(msg string, keyvals ...interface{}) {
	l.logs = append(l.logs, "ERROR: "+msg)
}

// ========== Test: NewTable ==========

func TestNewTableInitialization(t *testing.T) {
	// Test that basic Table initialization works without full widget creation
	config := createTestConfig()

	// Create table struct manually to avoid Fyne widget initialization in tests
	table := &Table{
		config:       config,
		data:         []interface{}{},
		state:        NewTableState(),
		KeyHandler:   NewDefaultKeyHandler(),
		MouseHandler: NewDefaultMouseHandler(),
		FocusHandler: NewDefaultFocusHandler(),
	}

	if table.config != config {
		t.Error("Config not properly assigned")
	}

	if table.data == nil {
		t.Error("Data should be initialized (empty slice)")
	}

	if table.state == nil {
		t.Error("State should be initialized")
	}

	if table.KeyHandler == nil {
		t.Error("KeyHandler should be initialized")
	}

	if table.MouseHandler == nil {
		t.Error("MouseHandler should be initialized")
	}

	if table.FocusHandler == nil {
		t.Error("FocusHandler should be initialized")
	}
}

func TestNewTableWithCustomLogger(t *testing.T) {
	config := createTestConfig()
	customLogger := &TestLogger{logs: []string{}}
	config.Logger = customLogger

	table := &Table{
		config: config,
		state:  NewTableState(),
		data:   []interface{}{},
	}

	// Verify logger is set by triggering a log message
	table.logger().Info("test message")
	if len(customLogger.logs) == 0 {
		t.Error("Custom logger not properly assigned or not functioning")
	}
}

// ========== Test: SetData / GetData ==========

func TestSetDataEmpty(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	table.SetData([]interface{}{})

	if len(table.GetData()) != 0 {
		t.Errorf("Expected empty data, got %d rows", len(table.GetData()))
	}
}

func TestSetDataSingleRow(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	data := []interface{}{
		TestData{ID: 1, Name: "Alice", Status: "Active", Priority: 1},
	}
	table.SetData(data)

	if len(table.GetData()) != 1 {
		t.Errorf("Expected 1 row, got %d", len(table.GetData()))
	}

	if table.GetData()[0].(TestData).Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", table.GetData()[0].(TestData).Name)
	}
}

func TestSetDataMultipleRows(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	data := createTestData()
	table.SetData(data)

	if len(table.GetData()) != 5 {
		t.Errorf("Expected 5 rows, got %d", len(table.GetData()))
	}
}

func TestSetDataReplacesPreviousData(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	// First set
	data1 := []interface{}{
		TestData{ID: 1, Name: "Alice", Status: "Active", Priority: 1},
		TestData{ID: 2, Name: "Bob", Status: "Inactive", Priority: 2},
	}
	table.SetData(data1)

	if len(table.GetData()) != 2 {
		t.Errorf("Expected 2 rows after first SetData, got %d", len(table.GetData()))
	}

	// Second set (replacing)
	data2 := []interface{}{
		TestData{ID: 3, Name: "Charlie", Status: "Active", Priority: 3},
	}
	table.SetData(data2)

	if len(table.GetData()) != 1 {
		t.Errorf("Expected 1 row after second SetData, got %d", len(table.GetData()))
	}

	if table.GetData()[0].(TestData).Name != "Charlie" {
		t.Error("Data not properly replaced")
	}
}

func TestSetDataPreservesSortOrder(t *testing.T) {
	config := createTestConfig()
	config.Columns[1].Sortable = true
	config.Columns[1].Comparator = func(a, b interface{}) int {
		return strings.Compare(
			a.(TestData).Name,
			b.(TestData).Name,
		)
	}

	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Manually sort by name (column 1)
	table.state.sortColumn = 1
	table.state.sortAsc = true
	table.sortData()

	// First should be "Alice"
	if table.GetData()[0].(TestData).Name != "Alice" {
		t.Errorf("Expected first sorted name 'Alice', got '%s'", table.GetData()[0].(TestData).Name)
	}

	// Now add new data - should re-apply sort
	newData := append(data, TestData{ID: 6, Name: "Aaron", Status: "Active", Priority: 1})
	table.SetData(newData)

	// First should now be "Aaron" (alphabetically first)
	if table.GetData()[0].(TestData).Name != "Aaron" {
		t.Errorf("Expected first sorted name 'Aaron' after SetData, got '%s'", table.GetData()[0].(TestData).Name)
	}
}

// ========== Test: Filtering ==========

func TestSetFilterPlainText(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Filter for "Active" - should match "Active" and "Inactive" (case-insensitive)
	table.SetFilter("Active", false)

	// Should show 4 rows: Alice (Active), Bob (Inactive), Charlie (Active), David (Active)
	// Row 3 (alice, Pending) should be excluded
	visibleCount := len(table.state.visibleRows)
	if visibleCount != 4 {
		t.Errorf("Expected 4 visible rows with 'Active' filter (matches Active and Inactive), got %d", visibleCount)
	}

	// Verify the excluded row is correct
	excluded := true
	for _, rowIdx := range table.state.visibleRows {
		if rowIdx == 3 { // alice with status "Pending"
			excluded = false
			break
		}
	}
	if !excluded {
		t.Error("Row 3 (alice, Pending) should be excluded")
	}

	// Verify filter state
	filterText, filterRegex := table.GetFilter()
	if filterText != "Active" {
		t.Errorf("Expected filter text 'Active', got '%s'", filterText)
	}
	if filterRegex {
		t.Error("Expected filterRegex to be false")
	}
}

func TestSetFilterCaseInsensitive(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Filter for "alice" (lowercase) - should match both "Alice" and "alice"
	table.state.filterCaseSensitive = false
	table.SetFilter("alice", false)

	visibleCount := len(table.state.visibleRows)
	if visibleCount != 2 {
		t.Errorf("Expected 2 visible rows with case-insensitive 'al' filter, got %d", visibleCount)
	}
}

func TestSetFilterCaseSensitive(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Enable case-sensitive filtering
	table.SetFilterCaseSensitive(true)
	table.SetFilter("Alice", false)

	// Should only match "Alice" (capital A), not "alice"
	visibleCount := len(table.state.visibleRows)
	if visibleCount != 1 {
		t.Errorf("Expected 1 visible row with case-sensitive 'Alice' filter, got %d", visibleCount)
	}

	// Verify the matched row is the correct one
	if len(table.state.visibleRows) > 0 {
		matchedRow := table.data[table.state.visibleRows[0]].(TestData)
		if matchedRow.Name != "Alice" {
			t.Errorf("Expected matched name 'Alice', got '%s'", matchedRow.Name)
		}
	}
}

func TestSetFilterRegex(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Test regex with anchor and case-sensitivity: ^P (starts with capital P)
	// This clearly demonstrates regex functionality (^ anchor) vs plain substring matching
	// Only row 3 has Status="Pending" which starts with capital P
	table.state.filterCaseSensitive = true // Enable case-sensitive for this test
	table.SetFilter("^P", true)

	visibleCount := len(table.state.visibleRows)
	if visibleCount != 1 {
		t.Errorf("Expected 1 visible row with regex '^P' filter (Status=Pending), got %d", visibleCount)
		// Debug: show what matched
		for _, rowIdx := range table.state.visibleRows {
			rowData := table.data[rowIdx].(TestData)
			t.Logf("Matched row: %+v", rowData)
		}
	}

	// Verify it's the correct row (alice with Status=Pending)
	if len(table.state.visibleRows) > 0 {
		matchedRow := data[table.state.visibleRows[0]].(TestData)
		if matchedRow.Status != "Pending" {
			t.Errorf("Expected matched row to have Status='Pending', got Status='%s'", matchedRow.Status)
		}
	}
}

func TestClearFilter(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Apply filter
	table.SetFilter("Active", false)
	if len(table.state.visibleRows) == len(data) {
		t.Error("Filter should reduce visible rows")
	}

	// Clear filter
	table.ClearFilter()

	// All rows should be visible
	if len(table.state.visibleRows) != len(data) {
		t.Errorf("Expected %d visible rows after ClearFilter, got %d", len(data), len(table.state.visibleRows))
	}

	// Filter text should be empty
	filterText, _ := table.GetFilter()
	if filterText != "" {
		t.Errorf("Expected empty filter text after ClearFilter, got '%s'", filterText)
	}
}

func TestSetFilterEmptyString(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Apply empty filter (should show all)
	table.SetFilter("", false)

	if len(table.state.visibleRows) != len(data) {
		t.Errorf("Expected all %d rows visible with empty filter, got %d", len(data), len(table.state.visibleRows))
	}
}

// ========== Test: Selection ==========

func TestSetSelectedCellValid(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	table.SetSelectedCell(2, 1)

	row, col := table.GetSelectedCell()
	if row != 2 || col != 1 {
		t.Errorf("Expected selected cell (2, 1), got (%d, %d)", row, col)
	}
}

func TestSetSelectedCellInvalidRow(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Try to select row 100 (out of bounds)
	table.SetSelectedCell(100, 1)

	// Selection should not change (should remain at default -1, -1)
	row, _ := table.GetSelectedCell()
	if row == 100 {
		t.Error("Invalid row index should not be accepted")
	}
}

func TestSetSelectedCellNegativeIndex(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Try to select negative row
	table.SetSelectedCell(-5, 1)

	// Selection should not change
	row, _ := table.GetSelectedCell()
	if row == -5 {
		t.Error("Negative row index should not be accepted")
	}
}

func TestSetSelectedRows(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Select multiple rows
	table.SetSelectedRows([]int{0, 2, 4})

	selectedRows := table.GetSelectedRows()
	if len(selectedRows) != 3 {
		t.Errorf("Expected 3 selected rows, got %d", len(selectedRows))
	}

	if selectedRows[0] != 0 || selectedRows[1] != 2 || selectedRows[2] != 4 {
		t.Errorf("Expected selected rows [0, 2, 4], got %v", selectedRows)
	}

	if table.GetSelectionCount() != 3 {
		t.Errorf("Expected selection count 3, got %d", table.GetSelectionCount())
	}
}

func TestClearSelection(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Select some rows
	table.SetSelectedRows([]int{0, 2, 4})

	// Clear selection
	table.ClearSelection()

	if table.GetSelectionCount() != 0 {
		t.Errorf("Expected 0 selected rows after ClearSelection, got %d", table.GetSelectionCount())
	}

	selectedRows := table.GetSelectedRows()
	if len(selectedRows) != 0 {
		t.Errorf("Expected empty selection after ClearSelection, got %v", selectedRows)
	}
}

// ========== Test: Column Visibility ==========

func TestSetColumnVisibility(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	// Initially all 4 columns should be visible
	if len(table.state.visibleColumns) != 4 {
		t.Errorf("Expected 4 visible columns initially, got %d", len(table.state.visibleColumns))
	}

	// Hide the "status" column (column index 2)
	table.SetColumnVisibility("status", false)

	// Should now have 3 visible columns
	if len(table.state.visibleColumns) != 3 {
		t.Errorf("Expected 3 visible columns after hiding one, got %d", len(table.state.visibleColumns))
	}

	// Verify "status" is not in visible columns
	for _, colIdx := range table.state.visibleColumns {
		if table.config.Columns[colIdx].ID == "status" {
			t.Error("Status column should be hidden")
		}
	}

	// Show the column again
	table.SetColumnVisibility("status", true)

	// Should have 4 visible columns again
	if len(table.state.visibleColumns) != 4 {
		t.Errorf("Expected 4 visible columns after showing column, got %d", len(table.state.visibleColumns))
	}
}

func TestSetColumnVisibilityInvalidID(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	initialCount := len(table.state.visibleColumns)

	// Try to hide a column that doesn't exist
	table.SetColumnVisibility("nonexistent", false)

	// Should not change visible columns count
	if len(table.state.visibleColumns) != initialCount {
		t.Error("Invalid column ID should not affect visible columns")
	}
}

// ========== Test: Column Read-Only ==========

func TestSetColumnReadOnly(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	// Initially "name" column (index 1) should be editable
	if table.config.Columns[1].ReadOnly {
		t.Error("Name column should initially be editable")
	}

	// Set "name" column to read-only
	table.SetColumnReadOnly("name", true)

	if !table.config.Columns[1].ReadOnly {
		t.Error("Name column should be read-only after SetColumnReadOnly")
	}

	// Set it back to editable
	table.SetColumnReadOnly("name", false)

	if table.config.Columns[1].ReadOnly {
		t.Error("Name column should be editable again")
	}
}

func TestSetColumnReadOnlyInvalidID(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	// Try to set read-only on non-existent column (should not panic)
	table.SetColumnReadOnly("nonexistent", true)

	// All columns should remain in their original state
	for _, col := range table.config.Columns {
		if col.ID == "nonexistent" {
			t.Error("Nonexistent column should not exist")
		}
	}
}

// ========== Test: Hierarchical/Tree Features ==========

func TestSetMaxDepth(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	// Initially max depth should be unlimited (0 = all levels)
	if table.config.MaxDepth != 0 {
		t.Errorf("Expected initial MaxDepth 0, got %d", table.config.MaxDepth)
	}

	// Set max depth to 3
	table.SetMaxDepth(3)

	if table.config.MaxDepth != 3 {
		t.Errorf("Expected MaxDepth 3, got %d", table.config.MaxDepth)
	}

	// Set max depth to 5
	table.SetMaxDepth(5)

	if table.config.MaxDepth != 5 {
		t.Errorf("Expected MaxDepth 5, got %d", table.config.MaxDepth)
	}
}

// ========== Test: Rendering ==========

func TestCreateRendererDoesNotPanic(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	// CreateRenderer should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CreateRenderer panicked: %v", r)
		}
	}()

	renderer := table.CreateRenderer()

	if renderer == nil {
		t.Error("CreateRenderer returned nil")
	}
}

func TestCreateRendererHasCorrectStructure(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)

	renderer := table.CreateRenderer()

	// Renderer should have objects (table + filter UI if visible)
	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Error("Renderer should have at least one object (the table)")
	}
}

// ========== Test: RebuildVisibleRows ==========

func TestRebuildVisibleRowsNoFilter(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	table.RebuildVisibleRows()

	// All rows should be visible when no filter is applied
	if len(table.state.visibleRows) != len(data) {
		t.Errorf("Expected %d visible rows, got %d", len(data), len(table.state.visibleRows))
	}
}

func TestRebuildVisibleRowsWithFilter(t *testing.T) {
	config := createTestConfig()
	table := createTestTable(config)
	data := createTestData()
	table.SetData(data)

	// Apply filter for "Active" - will match both "Active" and "Inactive"
	table.state.filterText = "Active"
	table.RebuildVisibleRows()

	// Should have fewer visible rows than total
	if len(table.state.visibleRows) >= len(data) {
		t.Error("Filter should reduce visible rows")
	}

	// Should match 4 rows (Alice-Active, Bob-Inactive, Charlie-Active, David-Active)
	if len(table.state.visibleRows) != 4 {
		t.Errorf("Expected 4 filtered rows, got %d", len(table.state.visibleRows))
	}

	// Verify filtered rows contain "active" (case-insensitive)
	for _, rowIdx := range table.state.visibleRows {
		rowData := data[rowIdx].(TestData)
		status := strings.ToLower(rowData.Status)
		if !strings.Contains(status, "active") {
			t.Errorf("Filtered row should contain 'active' in status, got %+v", rowData)
		}
	}
}
