package internal

import (
	"fmt"
	"sort"
)

// func CurlRegistry(osArgs []string, registries *Registries) {
// 	if len(osArgs) == 0 {
// 		curlAllRegistry(registries)
// 		return
// 	}
// 	keys := registries.RegistriesKeys
// 	name := osArgs[0]
// 	exist := in(name, keys)
// 	if !exist {
// 		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
// 			return fmt.Sprintf("%s[Grm]: please check alias exist%s", c.Red, c.Reset)
// 		})
// 		return
// 	}
// 	uri := registries.Registries[name].Registry
// 	ctx := curlRegistryImpl(uri)
// 	generatorCurlMessage(ctx, name)
// }

// func curlAllRegistry(registries *Registries) {

// 	keys := registries.RegistriesKeys
// 	source := registries.Registries
// 	for _, k := range keys {
// 		uri := source[k].Registry
// 		ctx := curlRegistryImpl(uri)
// 		generatorCurlMessage(ctx, k)
// 	}

// }

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

func in(tar string, arr []string) bool {

	sort.Strings(arr)
	idx := sort.SearchStrings(arr, tar)
	if idx < len(arr) && arr[idx] == tar {
		return true
	}
	return false
}
