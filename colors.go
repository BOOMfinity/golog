package golog

import "github.com/gookit/color"

var (
	errorColorCode   = color.StartSet + color.New(color.Red, color.Bold).Code() + "m"
	warningColorCode = color.StartSet + color.New(color.Yellow, color.Bold).Code() + "m"
	infoColorCode    = color.StartSet + color.New(color.Blue).Code() + "m"
	debugColorCode   = color.StartSet + color.New(color.Magenta).Code() + "m"
	traceColorCode   = color.StartSet + color.New(color.White).Code() + "m"
)
