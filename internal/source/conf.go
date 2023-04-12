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
)

var DefaultSource = map[string]S{
	"https://registry.npmjs.org/":                  Npm,
	"https://registry.yarnpkg.com/":                Yarn,
	"https://repo.huaweicloud.com/repository/npm/": HuaWei,
	"https://mirrors.cloud.tencent.com/npm/":       Tencent,
	"https://registry.npmmirror.com/":              NpmMirror,
}

var SourceToString = []string{
	"npm",
	"yarn",
	"huawei",
	"tencet",
	"npmMirror",
}

func (s S) String() string {
	return SourceToString[s]
}

type GrmConfig struct {
	baseDir  string
	confPath string
	aliases  []string
	parse    *GrmIni
}

func NewGrmConf() *GrmConfig {

	conf := &GrmConfig{
		baseDir:  path.Join(fs.SystemPreffix, "grm"),
		confPath: path.Join(fs.SystemPreffix, ".npmrc"),
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
