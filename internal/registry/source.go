package registry

import (
	"path"

	"github.com/modern-magic/grm/internal/fs"
)

type RegsitryInfo struct {
	Home     string
	Uri      string
	Internal bool
}

type SourceName struct {
	Npm       string
	Yarn      string
	HuaWei    string
	Tencent   string
	Cnpm      string
	TaoBao    string
	NpmMirror string
}

var PresetSourceName = SourceName{
	Npm:       "npm",
	Yarn:      "yarn",
	HuaWei:    "huawei",
	Tencent:   "tencent",
	Cnpm:      "cnpm",
	TaoBao:    "taobao",
	NpmMirror: "npmMirror",
}

var PresetRegistry = map[string]RegsitryInfo{
	PresetSourceName.Npm: {
		Home:     "https://www.npmjs.org",
		Uri:      "https://registry.npmjs.org/",
		Internal: true,
	},
	PresetSourceName.Yarn: {
		Home:     "https://yarnpkg.com",
		Uri:      "https://registry.yarnpkg.com/",
		Internal: true,
	},
	PresetSourceName.HuaWei: {
		Home:     "https://repo.huaweicloud.com/repository/npm/",
		Uri:      "https://repo.huaweicloud.com/repository/npm/",
		Internal: true,
	},
	PresetSourceName.Tencent: {
		Home:     "https://mirrors.cloud.tencent.com/npm/",
		Uri:      "https://mirrors.cloud.tencent.com/npm/",
		Internal: true,
	},
	PresetSourceName.Cnpm: {
		Home:     "https://cnpmjs.org",
		Uri:      "https://r.cnpmjs.org/",
		Internal: true,
	},
	PresetSourceName.NpmMirror: {
		Home:     "https://npmmirror.com",
		Uri:      "https://registry.npmmirror.com/",
		Internal: true,
	},
}

var (
	Grmrc = path.Join(fs.SystemPreffix, ".grmrc.yaml")
	Npmrc = path.Join(fs.SystemPreffix, ".npmrc")
)

type YAMLStruct struct {
	Home     string `yaml:"home"`
	Registry string `yaml:"registry"`
}

func Parsr(original map[string]RegsitryInfo) map[string]YAMLStruct {
	parserd := make(map[string]YAMLStruct, 0)
	for k, v := range original {
		parserd[k] = YAMLStruct{
			Home:     v.Home,
			Registry: v.Uri,
		}
	}
	return parserd
}

type Source interface {
	GetSource() map[string]RegsitryInfo
	GetUserSource() map[string]RegsitryInfo
}

type sourceImpl struct {
	fs         fs.FS
	source     map[string]RegsitryInfo
	userSource map[string]RegsitryInfo
}

func NewSource() Source {
	source := &sourceImpl{
		fs:         fs.NewFS(),
		source:     PresetRegistry,
		userSource: make(map[string]RegsitryInfo),
	}
	source.loadGRMConfig()
	return source
}

func (s *sourceImpl) loadGRMConfig() {

	content, err := s.fs.ReadYAML(Grmrc, map[string]YAMLStruct{})

	if err != nil {
		return
	}

	switch c := content.(type) {
	case map[string]YAMLStruct:
		for k, v := range c {
			s.source[k] = RegsitryInfo{
				Home: v.Home,
				Uri:  v.Registry,
			}
			s.userSource[k] = RegsitryInfo{
				Home: v.Home,
				Uri:  v.Registry,
			}
		}
	}
}

func (s *sourceImpl) GetSource() map[string]RegsitryInfo {
	return s.source
}

func (s *sourceImpl) GetUserSource() map[string]RegsitryInfo {
	return s.userSource
}
