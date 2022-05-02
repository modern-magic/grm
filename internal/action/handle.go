package action

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/registry"
)

func getCurrent() string {
	return registry.ReadNpm()
}

func ShowSources(sources *registry.RegistryDataSource) {

	outLen := len(sources.Keys) + 3

	cur := getCurrent()

	for _, key := range sources.Keys {
		prefix := ""
		uri := sources.Registry[key]
		if cur == uri {
			prefix = "* "
		}

		log := internal.StringJoin(prefix, key, getDashLine(key, outLen), uri, registry.Eol())

		if prefix == "" {
			fmt.Printf("%s", log)
		} else {
			logger.Success(log)

		}
	}

}

// show current registry uri and alias

func ShowCurrent() {
	cur := getCurrent()
	logger.Info(internal.StringJoin("[Grm]: you are using", cur))
}

func SetCurrent(sources *registry.RegistryDataSource, args []string) int {

	name := "npm"

	if len(args) >= 1 {
		name = args[0]
	}
	uri, ok := sources.Registry[name]
	if !ok {
		logger.Error(internal.StringJoin("[Grm]: Can't found alias", name, "in your .nrmrc file. Please check it exist.", registry.Eol()))
		return 1
	}
	registry.WriteNpm(registry.RegsitryInfo{
		Uri: uri,
	})
	logger.Success(internal.StringJoin("[Grm]: use", name, "success~", registry.Eol()))
	return 0
}

// del .nrm file registry alias

func DelRegistry(sources *registry.RegistryDataSource, args []string) int {

	if len(args) == 0 {
		return 0
	}
	name := args[0]

	_, ok := sources.Registry[name]

	if !ok {
		logger.Error(internal.StringJoin("[Grm]: Can't found alias", name, "in your .nrmrc file. Please check it exist.", registry.Eol()))
		return 1
	}
	flag, err := registry.DelNrm(name)
	if flag {
		logger.Success(internal.StringJoin("[Grm]: del registry", name, "success!", registry.Eol()))
		return 0
	}
	logger.Error(internal.StringJoin("[Grm]: del registry fail", err.Error(), registry.Eol()))
	return 1

}

func AddRegistry(args []string) int {

	name := ""
	home := ""
	uri := ""

	if len(args) <= 1 {
		logger.Error(internal.StringJoin("[Grm]: name and registry url is must be entry", registry.Eol()))
		return 1
	}
	name = args[0]
	uri = args[1]
	if len(args) == 2 {
		home = uri
	}
	if len(args) >= 3 {
		home = args[2]
	}

	flag, err := addRegistryImpl(name, uri, home)

	if flag {
		logger.Success(internal.StringJoin("[Grm]: add registry success!", registry.Eol()))
		return 0
	}
	logger.Error(internal.StringJoin("[Grm]: add registry fail", err.Error(), registry.Eol()))
	return 1
}

func FetchRegistry(sources *registry.RegistryDataSource, args []string) int {

	keys := make([]string, 0)

	if len(args) == 0 {
		keys = append(keys, sources.Keys...)
	} else {
		keys = append(keys, args[0])
	}
	if len(keys) == 1 {
		if _, ok := sources.Registry[keys[0]]; !ok {
			logger.Warn(internal.StringJoin("[Grm]: warning! can't found alias", keys[0], "will fetch npm source", registry.Eol()))
			keys[0] = "npm"
		}
	}
	for _, key := range keys {
		fetchRegistryImpl(sources.Registry[key], key)
	}
	return 0
}

func fetchRegistryImpl(uri, name string) {
	ctx := internal.Fetch(uri)
	log := "[Grm]: fetch " + name

	isTimeout := ctx.IsTimeout

	if isTimeout {
		log = internal.StringJoin(log, "state", ctx.Status)
	} else {
		log = internal.StringJoin(log, fmt.Sprintf("%.2f%s", ctx.Time, "s"), "state:", ctx.Status)
	}

	log = log + registry.Eol()

	if isTimeout {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s", c.Dim, log, c.Reset)
		})
		return
	}

	if ctx.StatusCode != 200 {
		logger.Error(log)
	} else {
		logger.Success(log)
	}
}

func addRegistryImpl(name, uri, home string) (bool, error) {
	return registry.WriteNrm(name, uri, home)

}

func getDashLine(key string, length int) string {
	final := math.Max(2, (float64(length) - float64(len(key))))
	bar := make([]string, int(final))
	for i := range bar {
		bar[i] = "-"
	}
	return strings.Join(bar[:], "-")
}
