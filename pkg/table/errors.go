package table

import "fmt"

// TableError represents an error that occurred during table operations.
type TableError struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

func (e *TableError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("table %s: %v", e.Op, e.Err)
	}
	return fmt.Sprintf("table %s", e.Op)
}

func (e *TableError) Unwrap() error {
	return e.Err
}

// ErrInvalidConfig creates a configuration error.
func ErrInvalidConfig(msg string) error {
	return &TableError{Op: "config", Err: fmt.Errorf("%s", msg)}
}
