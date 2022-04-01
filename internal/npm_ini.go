package internal

import (
	"io/ioutil"
	"strings"

	"github.com/go-ini/ini"
)

func getNrmRegistries() ([]RegistryMeta, []string) {
	cfg, _ := ini.LooseLoad(Nrmrc)

	sections := cfg.SectionStrings()

	nrmRegistries := make([]RegistryMeta, 0)
	nrmRegistriesKey := make([]string, 0)

	for _, sec := range sections {
		if strings.ToUpper(sec) == "DEFAULT" {
			continue
		}
		registry := cfg.Section(sec).Key(Registry).Value()
		home := cfg.Section(sec).Key(Home).Value()
		info := RegistryMeta{
			Home:     home,
			Registry: registry,
		}
		nrmRegistries = append(nrmRegistries, info)
		nrmRegistriesKey = append(nrmRegistriesKey, sec)
	}
	return nrmRegistries, nrmRegistriesKey
}

func writeNrmRegistries(r RegistryMeta, section string, args ...string) {
	cfg, _ := ini.LooseLoad(Nrmrc)
	ini.PrettyFormat = false
	blockPtr, _ := cfg.NewSection(section)

	if r.Home != "" {
		blockPtr.Key(Home).SetValue(r.Home)
	}
	blockPtr.Key(Registry).SetValue(r.Registry)
	// TODO
	if len(args) == 1 && args[0] == Delete {
		cfg.DeleteSection(section)
	}
	cfg.SaveTo(Nrmrc)
}

func readNpmRegistry() string {
	cfg, _ := ini.LooseLoad(Npmrc)
	return cfg.Section("").Key(Registry).Value()
}

// ....

func writeNpmRegistry(r RegistryMeta) {
	cfg, _ := ini.LoadSources(ini.LoadOptions{KeyValueDelimiterOnWrite: ":"}, Npmrc)
	ini.PrettyFormat = false
	cfg.Section("").Key(Home).SetValue(r.Home)
	cfg.Section("").Key(Registry).SetValue(r.Registry)
	cfg.SaveTo(Npmrc)
	buf, _ := ioutil.ReadFile(Npmrc)
	arr := strings.Split(string(buf), "\r\n")
	next := make([]string, 0)
	for _, k := range arr {
		k = strings.Replace(k, "registry:", "registry=", -1)
		k = strings.Replace(k, "home:", "home=", -1)
		next = append(next, k)
	}
	str := strings.Join(next, "\r\n")
	ioutil.WriteFile(Npmrc, []byte(str), 0644)
}
