package internal

import (
	"os"
	"path"
	"runtime"
)

func isWin() bool {
	return runtime.GOOS == "windows"
}

func eol() string {
	if isWin() {
		return "\r\n"
	}
	return "\n"
}

func getSystemEnv(key string) string {
	return os.Getenv(key)
}

func getSystemPreffix() string {
	win := isWin()
	if win {
		return getSystemEnv("USERPROFILE")
	}
	return getSystemEnv("HOME")
}

var (
	Home     = "home"
	Author   = "_author"
	Registry = "registry"
	Delete   = "delete"
	Nrmrc    = path.Join(getSystemPreffix(), ".nrmrc")
	Npmrc    = path.Join(getSystemPreffix(), ".npmrc")
)
