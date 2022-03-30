package internal

import (
	"io/ioutil"
	"strings"

	"github.com/go-ini/ini"
)

type Fs interface {
	ReadFile()
	WriteFile()
}

type Inis struct {
	name string
	path string
}

func getAllSections(cfg *ini.File) []string {
	return cfg.SectionStrings()
}

func getAllKeys(cfg *ini.File, key string) []string {
	return cfg.Section(key).KeyStrings()
}

func getKey(cfg *ini.File, key string) string {
	return cfg.Section(key).Key(Registry).Value()
}

func ReadAllFile(path string) []Inis {
	cfg, _ := ini.Load(path)
	inis := make([]Inis, 0)
	sections := getAllSections(cfg)
	for _, sec := range sections {
		if strings.ToUpper(sec) == "DEFAULT" {
			continue
		}
		path := getKey(cfg, sec)
		item := Inis{
			name: sec,
			path: path,
		}
		inis = append(inis, item)
	}
	return inis
}

func ReadFile(path string) string {
	cfg, _ := ini.Load(path)
	return cfg.Section("").Key(Registry).Value()
}

func WriteFile(path string, content RegistryInner) {
	cfg, _ := ini.LoadSources(ini.LoadOptions{KeyValueDelimiterOnWrite: ":"}, path)
	ini.PrettyFormat = false
	cfg.Section("").Key(Home).SetValue(content.home)
	cfg.Section("").Key(Registry).SetValue(content.registry)
	cfg.SaveTo(path)
	buf, _ := ioutil.ReadFile(path)
	arr := strings.Split(string(buf), "\r\n")
	next := make([]string, 0)
	for _, k := range arr {
		k = strings.Replace(k, "registry:", "registry=", -1)
		k = strings.Replace(k, "home:", "home=", -1)
		next = append(next, k)
	}
	str := strings.Join(next, "\r\n")
	ioutil.WriteFile(path, []byte(str), 0644)
}
