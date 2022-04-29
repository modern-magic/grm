package registry

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
	"github.com/modern-magic/grm/internal/logger"
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

func WriteNpm(info RegsitryInfo) {
	buf, err := ioutil.ReadFile(Npmrc)
	if err != nil {
		return
	}
	slice := strings.Split(string(buf), eol())
	next := make([]string, 0)
	for _, k := range slice {
		if strings.Index(k, "registry=") == 0 || strings.Index(k, "home=") == 0 {
			tmp := strings.Split(k, "=")
			tmp[1] = info.Uri
			k = strings.Join(tmp, "=")
		}
		next = append(next, k)
	}
	str := strings.Join(next, eol())
	err = ioutil.WriteFile(Npmrc, []byte(str), 0644)
	if err != nil {
		logger.PrintError("[Grm]: Error With" + err.Error())
		return
	}
}
