package table

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ColumnConfig defines the configuration for a single table column.
type ColumnConfig struct {
	ID         string             // Unique identifier for this column
	Label      string             // Display label in header
	Width      float32            // Column width in pixels (0 = auto)
	MinWidth   float32            // Minimum width (optional)
	Alignment  fyne.TextAlign     // Text alignment (default: TextAlignLeading)
	Sortable   bool               // Whether column can be sorted
	Visible    bool               // Whether column is visible
	Editable   bool               // Whether cells in this column are editable
	Formatter  func(interface{}) string // Custom cell formatter (optional)
}

// Config defines the configuration for creating a new Table widget.
type Config struct {
	// ID is a unique identifier for this table instance
	ID string

	// Columns defines the table structure
	Columns []ColumnConfig

	// Data is the slice of data items to display
	// Each item can be a map[string]interface{}, struct, or any type
	Data []interface{}

	// ShowHeaders controls whether column headers are displayed
	ShowHeaders bool

	// Callbacks for user interactions
	OnRowSelected func(row int, data interface{})
	OnCellEdited  func(row int, col string, newValue string, data interface{})
	OnKeyPressed  func(key *fyne.KeyEvent) bool // Return true to prevent default handling

	// Styling options
	HeaderBackgroundColor color.Color
	HeaderTextColor       color.Color
	SelectionColor        color.Color
	AlternateRowColor     color.Color

	// Logger for debugging (nil = NoopLogger)
	Logger Logger

	// Advanced options
	AllowMultiSort   bool // Allow sorting by multiple columns
	InitialSortCol   string
	InitialSortDesc  bool
	DisableSelection bool
}

// NewDefaultConfig creates a Config with sensible defaults.
func NewDefaultConfig(id string) *Config {
	return &Config{
		ID:                    id,
		ShowHeaders:           true,
		HeaderBackgroundColor: theme.PrimaryColor(),
		HeaderTextColor:       theme.ForegroundColor(),
		SelectionColor:        theme.SelectionColor(),
		Logger:                NoopLogger{},
	}
}

// Validate checks if the configuration is valid.
// Returns error if required fields are missing or invalid.
func (c *Config) Validate() error {
	if c.ID == "" {
		return ErrInvalidConfig("ID is required")
	}
	if len(c.Columns) == 0 {
		return ErrInvalidConfig("at least one column is required")
	}
	// More validation can be added here
	return nil
}
