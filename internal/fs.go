package internal

import (
	"github.com/go-ini/ini"
)

type Fs interface {
	ReadFile()
	WriteFile()
}

/*
this file is work for wrte ini data.
*/

func ReadFile(path string) string {
	cfg, err := ini.Load(path)
	if err != nil {

	}
	val := cfg.Section("").Key(Registry).Value()
	return val

}

func WriteFile() {}
