package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/digital-michael/fyne-components/pkg/table"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Table Widget Demo")
	window.Resize(fyne.NewSize(800, 600))

	// Generate mock task data
	mockData := generateMockTasks()

	// Convert to interface{} slice
	data := make([]interface{}, len(mockData))
	for i, task := range mockData {
		data[i] = task
	}

	// Configure table
	config := table.NewConfig("demo-table")
	config.ShowHeaders = true
	config.SelectFirstCellOnStartup = true
	config.FilterColumns = []string{"id", "name", "status"}
	config.ShowSearch = true
	config.SearchPlaceholder = "Search tasks by ID, Name, or Status..."
	config.RootNodeBackgroundColor = color.NRGBA{R: 35, G: 35, B: 65, A: 255}

	// Configure columns
	config.Columns = []table.ColumnConfig{
		{
			ID:         "id",
			Title:      "ID",
			Width:      60,
			Sortable:   true,
			Editable:   false,
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
		{
			ID:         "status",
			Title:      "Status",
			Width:      120,
			Sortable:   true,
			Editable:   false,
			Alignment:  table.AlignCenter,
			Comparator: table.NewStringComparator("Status"),
		},
		{
			ID:         "priority",
			Title:      "Priority",
			Width:      80,
			Sortable:   true,
			Editable:   false,
			Alignment:  table.AlignCenter,
			Comparator: table.NewNumericComparator("Priority"),
		},
		{
			ID:        "external",
			Title:     "External",
			Width:     100,
			Sortable:  true,
			Editable:  false,
			Alignment: table.AlignCenter,
		},
	}

	// Callbacks
	config.OnRowSelected = func(row int, rowData interface{}) {
		task, ok := rowData.(TaskData)
		if ok {
			fmt.Printf("Selected: %s (ID: %d)\n", task.Name, task.ID)
		}
	}

	config.OnCellEdited = func(row int, col string, newValue string, rowData interface{}) {
		task, ok := rowData.(TaskData)
		if ok {
			fmt.Printf("Edited: Row %d, Column %s, Task: %s, New Value: %s\n",
				row, col, task.Name, newValue)
			// Update the task name if editing name column
			if col == "name" {
				task.Name = newValue
				data[row] = task
			}
		}
	}

	// Create table widget
	tableWidget := table.NewTable(config)
	tableWidget.SetData(data)

	// Create info label
	infoLabel := widget.NewLabel("Table Widget Demo - Features:")
	instructions := widget.NewLabel(
		"• Click column headers to sort\n" +
			"• Use search box to filter tasks\n" +
			"• Double-click 'Task Name' cell to edit\n" +
			"• Arrow keys to navigate\n" +
			"• Tab/Shift+Tab to move between cells\n" +
			"• Enter to edit, Escape to cancel",
	)
	instructions.Wrapping = fyne.TextWrapWord

	// Create layout
	content := container.NewBorder(
		container.NewVBox(
			infoLabel,
			instructions,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		tableWidget,
	)

	window.SetContent(content)
	window.ShowAndRun()
}
