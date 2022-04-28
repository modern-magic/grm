package internal

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/modern-magic/grm/internal/logger"
)

func getCurrentRegisitry() string {
	return readNpmRegistry()
}

func ShowRegistries(registries *Registries) {

	keys := registries.RegistriesKeys
	source := registries.Registries
	curRegistry := getCurrentRegisitry()
	outLen := len(keys) + 3
	for _, k := range keys {
		prefix := "  "
		registry := source[k].Registry
		if curRegistry == registry {
			prefix = "* "
		}
		fmt.Print("\n", prefix, k, " ", getDashLine(k, outLen), " ", registry)
	}

}

func SetUsageRegistry(osArgs []string, registries *Registries) {

	keys := registries.RegistriesKeys
	name := ""
	// if user don't set other alias we will use npm registry as default.
	if len(osArgs) == 0 {
		name = "npm"
	} else {
		name = osArgs[0]
	}
	exist := in(name, keys)
	if !exist {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s[Grm:]%s%s", c.Red, "can't found registry"+name+"please check it exist.", c.Reset)
		})
		return
	}
	meta := registries.Registries[name]
	writeNpmRegistry(meta)
	fmt.Print("use ", name, " success!")
}

func ShowCurrentRegistry() {
	curRegistry := getCurrentRegisitry()
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s[Grm]: You are using:\n%s%s", c.Blue, curRegistry, c.Reset)
	})
}

func AddRegistry(osArgs []string) {
	name := ""
	home := ""
	registry := ""
	if len(osArgs) <= 1 {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s[Grm]: name and registry uri is must%s", c.Red, c.Reset)
		})
		return
	}
	if len(osArgs) >= 3 {
		home = osArgs[2]
	}

	name = osArgs[0]
	registry = osArgs[1]

	meta := RegistryMeta{
		Home:     home,
		Registry: registry,
	}

	writeNrmRegistries(meta, name)
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s[Grm]: add registry url success!", c.Green)
	})
}

func DelRegistry(osArgs []string, nrmKeys []string) {
	if len(osArgs) == 0 {
		return
	}
	name := osArgs[0]

	exist := in(name, nrmKeys)
	if !exist {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s[Grm]:can't found alias %s%s", c.Red, name, c.Reset)
		})
		return
	}
	writeNrmRegistries(RegistryMeta{}, name, Delete)
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s[Grm]: del registry %s success!%s", c.Green, name, c.Reset)
	})
}

func CurlRegistry(osArgs []string, registries *Registries) {
	if len(osArgs) == 0 {
		curlAllRegistry(registries)
		return
	}
	keys := registries.RegistriesKeys
	name := osArgs[0]
	exist := in(name, keys)
	if !exist {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s[Grm]: please check alias exist%s", c.Red, c.Reset)
		})
		return
	}
	uri := registries.Registries[name].Registry
	ctx := curlRegistryImpl(uri)
	generatorCurlMessage(ctx, name)
}

func curlAllRegistry(registries *Registries) {

	keys := registries.RegistriesKeys
	source := registries.Registries
	for _, k := range keys {
		uri := source[k].Registry
		ctx := curlRegistryImpl(uri)
		generatorCurlMessage(ctx, k)
	}

}

func generatorCurlMessage(ctx FetchContext, name string) {
	tout := ctx.isTimeout
	prefix := ""
	if tout {
		prefix = "fetch" + name + " " + "timeout"
	} else {
		prefix = "fetch " + name + " " + fmt.Sprintf("%v", ctx.time) + "s"
	}

	fmt.Println(prefix, "state:", ctx.status)

}

func curlRegistryImpl(uri string) FetchContext {
	return fetch(uri)
}

func getDashLine(key string, length int) string {
	final := math.Max(2, (float64(length) - float64(len(key))))
	bar := make([]string, int(final))
	for i := range bar {
		bar[i] = "-"
	}
	return strings.Join(bar[:], "-")
}

func in(tar string, arr []string) bool {

	sort.Strings(arr)
	idx := sort.SearchStrings(arr, tar)
	if idx < len(arr) && arr[idx] == tar {
		return true
	}
	return false
}
