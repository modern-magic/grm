package internal

import (
	"fmt"
	"math"
	"strings"
)

type registryInner struct {
	home     string
	registry string
}

var npm = registryInner{
	home:     "https://www.npmjs.org",
	registry: "https://registry.npmjs.org/",
}

var yarn = registryInner{
	home:     "https://yarnpkg.com",
	registry: "https://registry.yarnpkg.com/",
}

var tencet = registryInner{
	home:     "https://mirrors.cloud.tencent.com/npm/",
	registry: "https://mirrors.cloud.tencent.com/npm/",
}
var cnpm = registryInner{
	home:     "https://cnpmjs.org",
	registry: "https://r.cnpmjs.org/",
}

var taobao = registryInner{
	home:     "https://npmmirror.com",
	registry: "https://registry.npmmirror.com/",
}

var npmMirror = registryInner{
	home:     "https://skimdb.npmjs.com/",
	registry: "https://skimdb.npmjs.com/registry/",
}

var registries = map[string]registryInner{
	"npm":       npm,
	"yarn":      yarn,
	"tencet":    tencet,
	"cnpm":      cnpm,
	"taobao":    taobao,
	"npmMirror": npmMirror,
}


type Actions interface {
	ShowList()
}

func getCurrentRegisitry() string {
	return ReadFile(Npmrc)
}

func ShowList() {
	showList()
}

func showList() {
	curRegistry := getCurrentRegisitry()
	outLen := len(registries) + 3
	for k, v := range registries {
		reg := v.registry
		pre := "  "
		if curRegistry == reg {
			pre = "* "
		}
		fmt.Print("\n", pre, k, " ", getDashLine(k, outLen), " ", reg)
	}
}

func getPresetRegistries() {
}

func getDashLine(key string, length int) string {
	final := math.Max(2, (float64(length) - float64(len(key))))
	bar := make([]string, int(final))
	for i := range bar {
		bar[i] = "-"
	}
	return strings.Join(bar[:], "-")
}
