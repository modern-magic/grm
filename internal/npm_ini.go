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

/*
 * because go-ini can't compatible wtih npm-ini format, so we should do this by weself.
 */

func writeNpmRegistry(r RegistryMeta) {
	buf, err := ioutil.ReadFile(Npmrc)
	if err != nil {
		return
	}
	bufArr := strings.Split(string(buf), eol())
	next := make([]string, 0)
	for _, k := range bufArr {
		if strings.Index(k, "registry=") == 0 || strings.Index(k, "home=") == 0 {
			tmp := strings.Split(k, "=")
			tmp[1] = r.Registry
			k = strings.Join(tmp, "=")
		}
		next = append(next, k)
	}
	str := strings.Join(next, eol())
	ioutil.WriteFile(Npmrc, []byte(str), 0644)
}
