package logger

import (
	"fmt"
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
}

func PrintTextWithColor(file *os.File, callback func(Colors) string) {
	colors := TerminalColors
	writeStringWithColor(file, callback(colors))
}

func Info(text string) {
	PrintTextWithColor(os.Stdout, func(c Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Blue, text, c.Reset)
	})
}

func Error(text string) {
	PrintTextWithColor(os.Stderr, func(c Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Red, text, c.Reset)
	})
}

func Success(text string) {
	PrintTextWithColor(os.Stdout, func(c Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Green, text, c.Reset)
	})
}

func Warn(text string) {
	PrintTextWithColor(os.Stdout, func(c Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Yellow, text, c.Reset)
	})
}
