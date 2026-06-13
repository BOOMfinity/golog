package jsonfmt

import (
	"github.com/BOOMfinity/golog/v3/gcore"
)

// Config holds the JSON formatter settings.
type Config struct {
	Base                     gcore.Config
	Marshaler                func(any) ([]byte, error)
	DisableAttributeEscaping bool
	DisableMessageEscaping   bool
}
