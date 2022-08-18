//go:build !windows
// +build !windows

package fs

import "os"

var SystemPreffix = os.Getenv("HOME")
