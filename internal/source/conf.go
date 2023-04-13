package source

import (
	"path"

	"github.com/modern-magic/grm/internal/fs"
)

type S uint8

const (
	Npm S = iota
	Yarn
	HuaWei
	Tencent
	NpmMirror
	System
)

var DefaultSource = map[string]S{
	"https://registry.npmjs.org/":                  Npm,
	"https://registry.yarnpkg.com/":                Yarn,
	"https://repo.huaweicloud.com/repository/npm/": HuaWei,
	"https://mirrors.cloud.tencent.com/npm/":       Tencent,
	"https://registry.npmmirror.com/":              NpmMirror,
}

var DefaultKey = map[S]string{
	Npm:       "https://registry.npmjs.org/",
	Yarn:      "https://registry.yarnpkg.com/",
	HuaWei:    "https://repo.huaweicloud.com/repository/npm/",
	Tencent:   "https://mirrors.cloud.tencent.com/npm/",
	NpmMirror: "https://registry.npmmirror.com/",
}

func EnsureDefaultKey(input string) S {
	var s S
	switch input {
	case Npm.String():
		s = Npm
	case Yarn.String():
		s = Yarn
	case HuaWei.String():
		s = HuaWei
	case Tencent.String():
		s = Tencent
	case NpmMirror.String():
		s = NpmMirror
	default:
		s = System
	}
	return s
}

var SourceToString = []string{
	"npm",
	"yarn",
	"huawei",
	"tencet",
	"npmMirror",
	"system",
}

func (s S) String() string {
	return SourceToString[s]
}

type GrmConfig struct {
	baseDir  string
	ConfPath string
	aliases  []string
	parse    *GrmIni
}

func NewGrmConf() *GrmConfig {

	conf := &GrmConfig{
		baseDir:  path.Join(fs.SystemPreffix, "grm"),
		ConfPath: path.Join(fs.SystemPreffix, ".npmrc"),
	}
	conf.parse = NewGrmIniParse(conf)
	return conf
}

func (g *GrmConfig) ListAllPath() []string {
	dir := []string{}
	list := make([]string, 0, len(dir)+len(SourceToString))
	list = append(list, SourceToString...)
	list = append(list, dir...)
	g.aliases = list
	return list
}

func (g *GrmConfig) GetCurrentPath() string {
	g.parse.Get("registry")
	return g.parse.Path
}

func (g *GrmConfig) GetCurrentAlias() string {
	if g.parse.Path == "" {
		return ""
	}
	return ""
}

func (g *GrmConfig) SetCurrentPath(target string) bool {
	return g.parse.Set("registry", target)
}

func (g *GrmConfig) GetCurrentConf() string {
	return g.parse.ToString()
}
