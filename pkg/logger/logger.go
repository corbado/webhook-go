package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(err error)
}

type Impl struct {
}

var _ Logger = &Impl{}

// New returns new logger instance.
func New() *Impl {
	return &Impl{}
}

// Debug prints debug message
func (i *Impl) Debug(msg string, args ...any) {
	log.Printf("[DEBUG] %s\n", fmt.Sprintf(msg, args...))
}

// Error prints error message
func (i *Impl) Error(err error) {
	log.Printf("[ERROR] %+v\n", err)
}
