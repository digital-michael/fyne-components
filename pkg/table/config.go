package table

import (
	"fmt"
	"image/color"
	"reflect"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

// TextAlignment specifies how text is aligned in a cell
type TextAlignment int

const (
	AlignLeft TextAlignment = iota
	AlignCenter
	AlignRight
)

// TreeIconTheme defines the visual style for tree hierarchy indicators
type TreeIconTheme struct {
	Name   string
	Branch string   // Branch character (e.g., "├")
	Icons  []string // Symbol for each depth level (cycles if needed)
}

// Predefined tree icon themes
var (
	TreeThemeBranches = TreeIconTheme{
		Name:   "Branches",
		Branch: "├",
		Icons:  []string{"─ ", "─ ", "─ ", "─ ", "─ "},
	}
	TreeThemeArrows = TreeIconTheme{
		Name:   "Arrows",
		Branch: "├",
		Icons:  []string{"▷ ", "▷ ", "▷ ", "▷ ", "▷ "},
	}
	TreeThemeCircles = TreeIconTheme{
		Name:   "Circles",
		Branch: "├",
		Icons:  []string{"● ", "○ ", "● ", "○ ", "● "},
	}
	TreeThemeDiamonds = TreeIconTheme{
		Name:   "Diamonds",
		Branch: "├",
		Icons:  []string{"◆ ", "◇ ", "◆ ", "◇ ", "◆ "},
	}
	TreeThemeSquares = TreeIconTheme{
		Name:   "Squares",
		Branch: "├",
		Icons:  []string{"■ ", "□ ", "■ ", "□ ", "■ "},
	}
	TreeThemeBullets = TreeIconTheme{
		Name:   "Bullets",
		Branch: "├",
		Icons:  []string{"• ", "• ", "• ", "• ", "• "},
	}
	TreeThemeAngles = TreeIconTheme{
		Name:   "Angles",
		Branch: "├",
		Icons:  []string{"> ", "> ", "> ", "> ", "> "},
	}
)

// GetAllTreeThemes returns all available tree themes
func GetAllTreeThemes() []TreeIconTheme {
	return []TreeIconTheme{
		TreeThemeBranches,
		TreeThemeArrows,
		TreeThemeCircles,
		TreeThemeDiamonds,
		TreeThemeSquares,
		TreeThemeBullets,
		TreeThemeAngles,
	}
}

// ColumnConfig defines a single column
type ColumnConfig struct {
	ID       string // Unique column identifier
	Title    string // Header text
	Width    float32
	MinWidth float32

	// Behavior flags
	Sortable bool // true = clickable header for sorting
	Editable bool // true = inline editing enabled
	ReadOnly bool // true = prevent keyboard activation
	Hidden   bool // true = column is hidden by default

	// Visual styling
	Alignment TextAlignment // Text alignment (default: AlignLeft)

	// Custom rendering and logic
	Renderer   CellRenderer   // Custom cell content renderer
	Comparator SortComparator // Custom sort logic (nil = default string compare)

	// Popup menu for interactive cells
	PopupOptions    func(data interface{}) []string                            // Returns menu options for SPACE activation
	OnPopupSelected func(data interface{}, option string, rowIndex int) string // Returns new value when option selected
	ShowPopupIcon   bool                                                       // true = show dropdown icon in cell
	GetCellValue    func(data interface{}) string                              // Returns current cell text value for display

	// Checkbox for interactive boolean cells
	ShowCheckbox      bool                                               // true = render as checkbox
	GetCheckboxValue  func(data interface{}) bool                        // Returns current checkbox state
	OnCheckboxChanged func(data interface{}, checked bool, rowIndex int) // Called when checkbox is toggled
	CheckboxLabel     string                                             // Optional label text

	// Action callbacks
	OnViewData func(action string, data interface{}, colID string, rowIndex int, colIndex int) // Called for ESC, Return, Double-Click, etc.
}

// CellRenderer renders cell content
// Parameters:
//   - data: the data item for this row
//   - container: the Fyne container to populate
//   - rowIndex: the row index in the data array
//   - colID: the column identifier
type CellRenderer func(data interface{}, container fyne.CanvasObject, rowIndex int, colID string)

// SortComparator compares two data items for sorting
// Returns:
//   - negative if a < b
//   - zero if a == b
//   - positive if a > b
type SortComparator func(a, b interface{}) int

// Config defines the structure and behavior of a sortable table
type Config struct {
	// Core settings
	ID           string         // Unique identifier for persistence
	Columns      []ColumnConfig // Column definitions
	RowHeight    float32        // Default: 35px
	HeaderHeight float32        // Default: 30px

	// Features
	AllowMultiSelect  bool          // true = multi-select, false = single-select
	ShowSearch        bool          // true = show search box above table
	SearchPlaceholder string        // Search box placeholder text
	FilterTitle       string        // Card title for filter section (default: "Search/Filter")
	ShowToolbar       bool          // true = show toolbar with bulk actions
	TreeIconTheme     TreeIconTheme // Visual style for hierarchical indicators
	ShowBranch        bool          // true = show branch character (├), false = hide it
	RowSelectOnlyMode bool          // true = arrow keys select rows only, false = select row+column

	// Column Resizing
	EnableDoubleClickResize bool // true = double-click column divider to auto-resize

	// Header Control
	ShowHeaders bool // true = show column headers (also enables manual drag-resize), false = hide headers

	// Startup Selection
	SelectFirstCellOnStartup bool // true = automatically select cell (0,0) and set focus after data loaded

	// Indentation Control
	ShowIndentIcons bool    // true = show visual indent icons (├ └), false = hide them
	IndentPerLevel  float32 // Pixels to indent per hierarchy level (0 = no indent, default: 20)
	ShowIndentation bool    // true = apply indentation spacing, false = no indentation

	// Filter Control
	FilterColumns []string // Column IDs to search/filter (empty = no filtering)

	// Tree Hierarchy Control
	MaxDepth         int                                // Maximum depth to display (0 or nil = show all levels)
	ExpandedNodes    map[interface{}]bool               // Track which nodes are expanded (nil = all expanded)
	GetNodeID        func(data interface{}) interface{} // Get unique ID for a node (for expand/collapse tracking)
	GetNodeDepth     func(data interface{}) int         // Get depth of a node
	GetNodeParentID  func(data interface{}) interface{} // Get parent ID of a node (nil = root)
	IsNodeExpandable func(data interface{}) bool        // Check if node has children

	// Visual Styling
	RootNodeBackgroundColor color.Color // Background color for root nodes (depth 0), nil = no background
	FontFamily              string      // Font family name (empty = default)
	FontSize                float32     // Font size in points (0 = default)

	// Logging
	Logger Logger // Logger interface for structured logging (nil = use NoopLogger)

	// Callbacks
	OnRowSelected func(rowIndex int, data interface{})
	OnRowAction   func(action string, rowIndex int, data interface{})
	OnCellEdited  func(rowIndex int, colID string, newValue string, data interface{})

	// Persistence (optional)
	SaveColumnWidths func(widths map[string]float32)
	LoadColumnWidths func() map[string]float32
}

// NewConfig creates a default table configuration
func NewConfig(id string) *Config {
	return &Config{
		ID:                      id,
		Columns:                 []ColumnConfig{},
		RowHeight:               35.0,
		HeaderHeight:            30.0,
		AllowMultiSelect:        false,
		ShowSearch:              false,
		SearchPlaceholder:       "Search...",
		FilterTitle:             "Search/Filter",
		ShowToolbar:             false,
		ShowBranch:              true,
		RowSelectOnlyMode:       true,
		EnableDoubleClickResize: true,
		ShowHeaders:             true,
		ShowIndentIcons:         true,
		IndentPerLevel:          20.0,
		ShowIndentation:         true,
		MaxDepth:                0,                          // Show all levels by default
		ExpandedNodes:           make(map[interface{}]bool), // All nodes collapsed initially
		RootNodeBackgroundColor: nil,                        // No background by default
		FontFamily:              "",                         // System default
		FontSize:                0,                          // System default (usually 12-14pt)
		Logger:                  NoopLogger{},               // Default noop logger
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
	return nil
}

// NewStringComparator creates a comparator that extracts a string field and compares lexicographically
func NewStringComparator(fieldName string) SortComparator {
	return func(a, b interface{}) int {
		valA := extractFieldString(a, fieldName)
		valB := extractFieldString(b, fieldName)
		if valA < valB {
			return -1
		} else if valA > valB {
			return 1
		}
		return 0
	}
}

// NewNumericComparator creates a comparator that extracts a numeric field and compares numerically
func NewNumericComparator(fieldName string) SortComparator {
	return func(a, b interface{}) int {
		valA := extractFieldNumeric(a, fieldName)
		valB := extractFieldNumeric(b, fieldName)
		if valA < valB {
			return -1
		} else if valA > valB {
			return 1
		}
		return 0
	}
}

// extractFieldString extracts a string value from a struct field by name (case-insensitive)
func extractFieldString(data interface{}, fieldName string) string {
	if data == nil {
		return ""
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", data)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if strings.EqualFold(field.Name, fieldName) {
			fieldValue := v.Field(i)
			if fieldValue.CanInterface() {
				return fmt.Sprintf("%v", fieldValue.Interface())
			}
		}
	}

	return fmt.Sprintf("%v", data)
}

// extractFieldNumeric extracts a numeric value from a struct field by name (case-insensitive)
func extractFieldNumeric(data interface{}, fieldName string) float64 {
	if data == nil {
		return 0
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return toFloat64(data)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if strings.EqualFold(field.Name, fieldName) {
			fieldValue := v.Field(i)
			if fieldValue.CanInterface() {
				return toFloat64(fieldValue.Interface())
			}
		}
	}

	return 0
}

// toFloat64 converts various numeric types to float64
func toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num
		}
		return 0
	default:
		return 0
	}
}
