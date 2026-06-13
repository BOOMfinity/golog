package gcore

import (
	"context"
	"unsafe"
)

// Context provides read-only access to all important information about log entry.
//
// Strings and slices are pre-allocated, make a copy before storing or passing to other methods!
type Context interface {
	// Modules returns fixed-size slices for modules and scopes.
	// Both slices are always the same size, where each element represents a sub-logger in the path.
	Modules() (modules, scopes []string)
	// Level returns the severity level of this log entry.
	Level() Level
	// Attributes returns global (saved in [Logger]) and local (added directly to the [Message]) attributes.
	Attributes() (global, local []Attribute)
	// Message returns the message of current log entry.
	//
	// Returns an empty string if message is empty or not set.
	Message() string
	// HasMessage returns true if message is not empty.
	//
	// It is equivalent to Context.Message() != "".
	HasMessage() bool
	// StackTrace returns the stack trace of the current goroutine.
	//
	// If Message.Stack() wasn't called, the returned string is empty.
	StackTrace() string
	// HasStackTrace returns true if this log entry has a stack trace.
	//
	// It's just a shortcut for Context.StackTrace() != "".
	HasStackTrace() bool
	// Error returns an error added by the user to this log entry.
	//
	// It returns nil if no error was set.
	Error() error
	// HasError returns true if this log entry contains an error.
	//
	// It's just a shortcut for Context.Error() != nil.
	HasError() bool
	// Context returns the context of this log entry.
	//
	// If no context was set, it falls back to the [Logger]'s context or context.Background().
	Context() context.Context
	// ExitCode returns the exit code used to terminate the process.
	//
	// If the returned value is 0, that means the exit code was not set and process will not exit.
	ExitCode() int
	// HasExitCode returns true if exit code was set.
	HasExitCode() bool
}

type contextImpl Message

func (c *contextImpl) Modules() (modules, scopes []string) {
	return c.parent.modules, c.parent.scopes
}

func (c *contextImpl) Level() Level {
	return c.level
}

func (c *contextImpl) Attributes() (global, local []Attribute) {
	return c.parent.attrs, c.attrs
}

func (c *contextImpl) Message() string {
	return unsafe.String(unsafe.SliceData(c.message), len(c.message))
}

func (c *contextImpl) HasMessage() bool {
	return len(c.message) > 0
}

func (c *contextImpl) StackTrace() string {
	return unsafe.String(unsafe.SliceData(c.stack), len(c.stack))
}

func (c *contextImpl) HasStackTrace() bool {
	return len(c.stack) > 0
}

func (c *contextImpl) Error() error {
	return c.err
}

func (c *contextImpl) HasError() bool {
	return c.err != nil
}

func (c *contextImpl) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	if c.parent.ctx != nil {
		return c.parent.ctx
	}
	return context.Background()
}

func (c *contextImpl) ExitCode() int {
	return c.exitCode
}

func (c *contextImpl) HasExitCode() bool {
	return c.exitCode != 0
}
