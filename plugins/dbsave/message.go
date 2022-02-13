package golog_db

import (
	"github.com/VenomPCPL/golog"
)

type Message struct {
	Level   golog.Level `bson:"level"`
	Message []byte      `bson:"message"`
	Time    int64       `bson:"time"`
	Tags    []string    `json:"tags,omitempty"`
}
