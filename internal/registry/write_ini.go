package registry

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
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

func Eol() string {
	return eol()
}

func GetSystemPreffix() string {
	if isWin() {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func readini(file string) (*ini.File, error) {
	cfg, err := ini.LooseLoad(file)
	return cfg, err
}

func ReadNrm() (map[string]RegsitryInfo, []string) {
	cfg, _ := readini(Nrmrc)
	secs := cfg.SectionStrings()
	keys := make([]string, 0)
	nrmRegistry := make(map[string]RegsitryInfo)
	for _, sec := range secs {
		if strings.ToUpper(sec) == Default {
			continue
		}
		info := RegsitryInfo{
			Home: cfg.Section(sec).Key(Home).Value(),
			Uri:  cfg.Section(sec).Key(Registry).Value(),
		}
		nrmRegistry[sec] = info
		keys = append(keys, sec)
	}
	return nrmRegistry, keys
}

func WriteNrm(name, uri, home string) (bool, error) {
	cfg := writeNrmImpl(name, uri, home)
	err := cfg.SaveTo(Nrmrc)
	return err == nil, err
}

func writeNrmImpl(name, uri, home string) *ini.File {
	cfg, _ := readini(Nrmrc)
	ini.PrettyFormat = false
	blockPtr, _ := cfg.NewSection(name)
	blockPtr.Key(Home).SetValue(home)
	blockPtr.Key(Registry).SetValue(uri)
	return cfg
}

func DelNrm(name string) (bool, error) {
	cfg := writeNrmImpl(name, "", "")
	cfg.DeleteSection(name)
	err := cfg.SaveTo(Nrmrc)
	return err == nil, err
}

func ReadNpm() string {
	cfg, _ := readini(Npmrc)
	return cfg.Section("").Key(Registry).Value()
}

func WriteNpm(info RegsitryInfo) (bool, error) {
	buf, err := ioutil.ReadFile(Npmrc)
	if err != nil {
		return false, err
	}

	return writeNpmImpl(buf, info.Uri)
}

func writeNpmImpl(buf []byte, uri string) (bool, error) {
	sections := strings.Split(string(buf), eol())
	next := make([]string, 0)
	for _, key := range sections {
		if strings.Index(key, "registry=") == 0 || strings.Index(key, "home=") == 0 {
			temp := strings.Split(key, "=")
			temp[1] = uri
			line := strings.Join(temp, "=")
			next = append(next, line)
		} else {
			next = append(next, key)
		}
	}
	after := strings.Join(next, eol())
	err := ioutil.WriteFile(Npmrc, []byte(after), 0644)
	return err == nil, err
}
