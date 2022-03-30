package internal

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type RegistryInner struct {
	home     string
	registry string
}

var npm = RegistryInner{
	home:     "https://www.npmjs.org",
	registry: "https://registry.npmjs.org/",
}

var yarn = RegistryInner{
	home:     "https://yarnpkg.com",
	registry: "https://registry.yarnpkg.com/",
}

var tencet = RegistryInner{
	home:     "https://mirrors.cloud.tencent.com/npm/",
	registry: "https://mirrors.cloud.tencent.com/npm/",
}
var cnpm = RegistryInner{
	home:     "https://cnpmjs.org",
	registry: "https://r.cnpmjs.org/",
}

var taobao = RegistryInner{
	home:     "https://npmmirror.com",
	registry: "https://registry.npmmirror.com/",
}

var npmMirror = RegistryInner{
	home:     "https://skimdb.npmjs.com/",
	registry: "https://skimdb.npmjs.com/registry/",
}

var registries = map[string]RegistryInner{
	"npm":       npm,
	"yarn":      yarn,
	"tencet":    tencet,
	"cnpm":      cnpm,
	"taobao":    taobao,
	"npmMirror": npmMirror,
}

type Actions interface {
	ShowList()
	SetUse()
	ShowCurrent()
}

func getCurrentRegisitry() string {
	return ReadFile(Npmrc)
}

func ShowList() {
	cur := getCurrentRegisitry()
	getAllRegistries()
	outLen := len(registries) + 3
	for k, v := range registries {
		reg := v.registry
		pre := "  "
		if cur == reg {
			pre = "* "
		}
		fmt.Print("\n", pre, k, " ", getDashLine(k, outLen), " ", reg)
	}
}

func SetUse(name string) {
	getAllRegistries()
	alias := make([]string, 0)
	for k := range registries {
		alias = append(alias, k)
	}
	exists := in(name, alias)
	if !exists {
		return
	}
	registry := registries[name]
	WriteFile(Npmrc, registry)
	fmt.Print("use ", name, " success!")
}

func ShowCurrent() {
	cur := getCurrentRegisitry()
	getAllRegistries()
	alias := make([]string, 0)
	for _, v := range registries {
		alias = append(alias, v.registry)
	}
	exist := in(cur, alias)
	if exist {
		fmt.Print("you are using url: ", cur)
		return
	}
	fmt.Print("can't fount alias")
}

func getAllRegistries() {
	all := ReadAllFile(Nrmrc)
	for _, v := range all {
		inner := RegistryInner{
			home:     v.path,
			registry: v.path,
		}
		registries[v.name] = inner
	}
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
