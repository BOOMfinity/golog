package golog

import (
	"context"
	"errors"

	"github.com/BOOMfinity/golog/v3/formats/colorfmt"
	"github.com/BOOMfinity/golog/v3/formats/jsonfmt"
	"github.com/BOOMfinity/golog/v3/gcore"
)

func ExampleInit_json() {
	log := Init("app", jsonfmt.Init())

	log.Info().Msg("hello world!")
}

func ExampleInit_color() {
	log := Init("app", colorfmt.Init())

	log.Info().Msg("hello world!")
}

// Passing context to the log entry can be used to access context-bound values in the Handler.
func Example_context() {
	// Create a sub-logger.
	log := Default().Module("db")

	// Set a base context for all future log entries.
	log = log.WithContext(context.TODO())

	// Or set context for a single entry. It overrides the logger's base context.
	log.Info(context.TODO()).Msg("hello world!")
}

func Example() {
	// You can use the default logger from anywhere in your code.
	Default().Info().Msg("hello world!")

	// Or use helpers.
	Info().Msg("hello world!")

	// Or just create custom logger.
	log := Init("app")
	log.Info().Msg("hello world!")
}

// Global attributes are added to every log entry from this logger.
// Local attributes are added to a single entry.
func Example_attributes() {
	log := Default().WithAttribute("service", "api")
	log.Info().Attr("reqID", "abc").Msg("request started")
}

// Modules create named sub-loggers, optionally with a scope.
func Example_modules() {
	db := Default().Module("db", "queries")
	db.Debug().Msg("SELECT * FROM users")
}

// Attach an error to a log entry — the formatter decides how to display it.
func Example_error() {
	err := errors.New("connection refused")
	Default().Error().Cause(err).Msg("database unreachable")
}

// Recover catches panics, logs them, and resumes execution.
func Example_recover() {
	defer Default().Recover()
	panic("unexpected nil pointer")
}

// Hooks run before the formatter — useful for external monitoring.
func Example_hook() {
	hook := func(entry gcore.Context) {
		if entry.HasError() {
			// send error to external monitoring system
		}
	}
	log := Default().WithHook(hook)
	log.Error().MsgErr(errors.New("critical"))
}

// Use a custom formatter with configuration.
func Example_custom() {
	log := Init("myapp", colorfmt.Init(colorfmt.Config{
		Base: gcore.Config{TimeFormat: "15:04:05"},
	}))
	log.Info().Msg("custom time format")
}
