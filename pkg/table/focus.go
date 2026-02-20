package table

import (
	"fyne.io/fyne/v2"
)

// FocusHandler defines the interface for handling focus events and management in a table
type FocusHandler interface {
	// HandleFocusGained is called when the table gains keyboard focus
	HandleFocusGained(table *Table)
	// HandleFocusLost is called when the table loses keyboard focus
	HandleFocusLost(table *Table)
	// RequestFocus requests keyboard focus for the table
	RequestFocus(table *Table)
}

// DefaultFocusHandler provides the standard focus behavior for tables
type DefaultFocusHandler struct{}

// NewDefaultFocusHandler creates a new default focus handler
func NewDefaultFocusHandler() *DefaultFocusHandler {
	return &DefaultFocusHandler{}
}

// HandleFocusGained is called when the table gains keyboard focus
func (h *DefaultFocusHandler) HandleFocusGained(table *Table) {
	table.state.hasFocus = true
}

// HandleFocusLost is called when the table loses keyboard focus
func (h *DefaultFocusHandler) HandleFocusLost(table *Table) {
	table.state.hasFocus = false
}

// RequestFocus requests keyboard focus for the table
func (h *DefaultFocusHandler) RequestFocus(table *Table) {
	// Focus the doubleTappableTable wrapper (which is in the canvas tree via CreateRenderer)
	// NOT table.table.Table - that's embedded in table.table and not directly in canvas
	if table.table != nil {
		if canvas := fyne.CurrentApp().Driver().CanvasForObject(table.table); canvas != nil {
			canvas.Focus(table.table)
			// Verify focus was set (for debugging focus issues)
			if canvas.Focused() == nil {
				table.logger().Warn("Canvas reports NO focused object after Focus() call")
			}
		}
		// If canvas not ready, silently return - caller should use fyne.Do() to defer if needed
	}
	// If table.table is nil, silently return - not an error during initialization
}
