//go:build !windows
// +build !windows

package logger

import "os"

func writeStringWithColor(file *os.File, text string) {
	file.WriteString(text) //nolint
}
