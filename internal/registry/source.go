package registry

import (
	"path"
	"reflect"
)

type RegsitryInfo struct {
	Home string
	Uri  string
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
		Home: "https://www.npmjs.org",
		Uri:  "https://registry.npmjs.org/",
	},
	PresetSourceName.Yarn: {
		Home: "https://yarnpkg.com",
		Uri:  "https://registry.yarnpkg.com/",
	},
	PresetSourceName.HuaWei: {
		Home: "https://repo.huaweicloud.com/repository/npm/",
		Uri:  "https://repo.huaweicloud.com/repository/npm/",
	},
	PresetSourceName.Tencent: {
		Home: "https://mirrors.cloud.tencent.com/npm/",
		Uri:  "https://mirrors.cloud.tencent.com/npm/",
	},
	PresetSourceName.Cnpm: {
		Home: "https://cnpmjs.org",
		Uri:  "https://r.cnpmjs.org/",
	},
	PresetSourceName.TaoBao: {
		Home: "https://npmmirror.com",
		Uri:  "https://registry.npmmirror.com/",
	},
	PresetSourceName.NpmMirror: {
		Home: "https://skimdb.npmjs.com/",
		Uri:  "https://skimdb.npmjs.com/Registry/",
	},
}

func GetPresetRegistryNames() []string {
	var dict interface{} = PresetSourceName
	names := reflect.ValueOf(dict)
	sequenNames := make([]string, names.NumField())
	for i := 0; i < names.NumField(); i++ {
		sequenNames[i] = names.Field(i).String()
	}
	return sequenNames
}

func GetPresetRegistryInfo(kind string) string {
	info, ok := PresetRegistry[kind]
	if !ok {
		panic("Invalid Source")
	}
	return info.Uri

}

var (
	Home     = "home"
	Author   = "_author"
	Registry = "registry"
	Delete   = "delete"
	Default  = "DEFAULT"
	Nrmrc    = path.Join(GetSystemPreffix(), ".nrmrc")
	Grmrc    = path.Join(GetSystemPreffix(), ".grmrc.yaml")
	Npmrc    = path.Join(GetSystemPreffix(), ".npmrc")
)

type RegistryDataSource struct {
	Registry     map[string]string
	Keys         []string
	UserRegistry map[string]string
}
