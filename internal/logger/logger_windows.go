//go:build windows
// +build windows

package logger

import (
	"os"
	"strings"
	"syscall"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var setConsoleTextAttribute = kernel32.NewProc("SetConsoleTextAttribute")

const (
	FOREGROUND_BLUE uint8 = 1 << iota
	FOREGROUND_GREEN
	FOREGROUND_RED
	FOREGROUND_INTENSITY
	BACKGROUND_BLUE
	BACKGROUND_GREEN
	BACKGROUND_RED
	BACKGROUND_INTENSITY
)

var windowsEscapeSequenceMap = map[string]uint8{
	TerminalColors.Reset: FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE,
	TerminalColors.Dim:   FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE,
	TerminalColors.Bold:  FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE | FOREGROUND_INTENSITY,

	TerminalColors.Red:   FOREGROUND_RED,
	TerminalColors.Green: FOREGROUND_GREEN,
	TerminalColors.Blue:  FOREGROUND_BLUE,

	TerminalColors.Cyan:    FOREGROUND_GREEN | FOREGROUND_BLUE,
	TerminalColors.Magenta: FOREGROUND_RED | FOREGROUND_BLUE,
	TerminalColors.Yellow:  FOREGROUND_RED | FOREGROUND_GREEN,
}

func writeStringWithColor(file *os.File, text string) {
	fd := file.Fd()
	i := 0
	for i < len(text) {
		if text[i] != 033 {
			i++
			continue
		}

		window := text[i:]
		if len(window) > 8 {
			window = window[:8]
		}
		m := strings.IndexByte(window, 'm')
		if m == -1 {
			i++
			continue
		}
		m += i + 1

		attributes, ok := windowsEscapeSequenceMap[text[i:m]]
		if !ok {
			i++
			continue
		}

		file.WriteString(text[:i]) //nolint

		text = text[m:]
		i = 0
		setConsoleTextAttribute.Call(fd, uintptr(attributes)) //nolint
	}

	file.WriteString(text) //nolint
}
