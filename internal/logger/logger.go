package logger

import (
	"os"
)

type Colors struct {
	Reset string
	Bold  string
	Dim   string

	Red   string
	Green string
	Blue  string

	Cyan    string
	Magenta string
	Yellow  string

	DimCyan string
}

var TerminalColors = Colors{
	Reset: "\033[0m",
	Bold:  "\033[1m",
	Dim:   "\033[37m",

	Red:   "\033[31m",
	Green: "\033[32m",
	Blue:  "\033[34m",

	Cyan:    "\033[36m",
	Magenta: "\033[35m",
	Yellow:  "\033[33m",

	DimCyan: "\033[90m",
}

func PrintTextWithColor(file *os.File, callback func(Colors) string) {
	colors := TerminalColors
	writeStringWithColor(file, callback(colors))
}
