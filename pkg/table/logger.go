package table

// Logger defines the interface for table widget logging.
// By default, the table uses NoopLogger (silent). Users can inject
// their own logger implementation for debugging.
//
// The logger interface is compatible with common logging libraries
// and follows structured logging conventions with key-value pairs.
type Logger interface {
	// Debug logs debug-level messages with optional key-value pairs
	Debug(msg string, keyvals ...interface{})

	// Info logs info-level messages with optional key-value pairs
	Info(msg string, keyvals ...interface{})

	// Warn logs warning-level messages with optional key-value pairs
	Warn(msg string, keyvals ...interface{})

	// Error logs error-level messages with optional key-value pairs
	Error(msg string, keyvals ...interface{})
}

// NoopLogger is a logger that discards all log messages.
// This is the default logger used when Config.Logger is nil.
type NoopLogger struct{}

// Debug implements Logger.Debug by discarding the message
func (NoopLogger) Debug(msg string, keyvals ...interface{}) {}

// Info implements Logger.Info by discarding the message
func (NoopLogger) Info(msg string, keyvals ...interface{}) {}

// Warn implements Logger.Warn by discarding the message
func (NoopLogger) Warn(msg string, keyvals ...interface{}) {}

// Error implements Logger.Error by discarding the message
func (NoopLogger) Error(msg string, keyvals ...interface{}) {}

// StdLogger is a simple logger adapter that writes to a provided writer
// (typically os.Stdout or os.Stderr). This provides basic logging without
// requiring external dependencies.
//
// Example usage:
//
//	import "os"
//
//	logger := table.NewStdLogger(os.Stdout)
//	config := &table.Config{
//		Logger: logger,
//		// ... other config
//	}
type StdLogger struct {
	output func(format string, args ...interface{})
}

// NewStdLogger creates a logger that writes to the provided print function.
// Typically used with fmt.Printf or log.Printf.
//
// Example:
//
//	logger := table.NewStdLogger(fmt.Printf)
func NewStdLogger(printf func(format string, args ...interface{})) Logger {
	return &StdLogger{output: printf}
}

// Debug logs debug messages with [DEBUG] prefix
func (l *StdLogger) Debug(msg string, keyvals ...interface{}) {
	l.output("[DEBUG] %s %v\n", msg, keyvals)
}

// Info logs info messages with [INFO] prefix
func (l *StdLogger) Info(msg string, keyvals ...interface{}) {
	l.output("[INFO] %s %v\n", msg, keyvals)
}

// Warn logs warning messages with [WARN] prefix
func (l *StdLogger) Warn(msg string, keyvals ...interface{}) {
	l.output("[WARN] %s %v\n", msg, keyvals)
}

// Error logs error messages with [ERROR] prefix
func (l *StdLogger) Error(msg string, keyvals ...interface{}) {
	l.output("[ERROR] %s %v\n", msg, keyvals)
}
