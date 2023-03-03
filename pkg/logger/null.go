package logger

type Null struct {
}

var _ Logger = &Null{}

// NewNull returns new null logger instance which can be used in unit tests for example.
func NewNull() *Null {
	return &Null{}
}

// Debug prints debug message
func (n *Null) Debug(msg string, args ...any) {
}

// Error prints error message
func (n *Null) Error(_ error) {
}
