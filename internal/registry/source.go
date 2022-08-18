package registry

import (
	"errors"
	"path"
	"reflect"

	"github.com/modern-magic/grm/internal/fs"
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
		panic(SourceStatus.Invalid)
	}
	return info.Uri

}

var (
	Grmrc = path.Join(fs.SystemPreffix, ".grmrc.yaml")
	Npmrc = path.Join(fs.SystemPreffix, ".npmrc")
)

type RegistryStatus struct {
	NotFound string
	Invalid  string
	Exists   string
}

var SourceStatus = RegistryStatus{
	NotFound: "Not found",
	Invalid:  "Invalid Source",
	Exists:   "Already exists",
}

type RegistryDataSource struct {
	fs         fs.FS
	Registry   map[string]string
	Keys       []string
	PresetKeys []string
	Niave      map[string]RegsitryInfo
}

func (r *RegistryDataSource) Drop(name string) error {
	if _, ok := r.Niave[name]; ok {
		delete(r.Niave, name)
		parsed := parsr(r.Niave)
		err := r.fs.WriteYAML(Grmrc, parsed)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(SourceStatus.NotFound)
}

func (r *RegistryDataSource) Insert(name string, uri string, home string) error {
	if _, ok := r.Niave[name]; ok {
		return errors.New(SourceStatus.Exists)
	}
	r.Niave[name] = RegsitryInfo{
		Home: home,
		Uri:  uri,
	}
	parsed := parsr(r.Niave)
	err := r.fs.WriteYAML(Grmrc, parsed)
	if err != nil {
		return err
	}
	return nil
}

func parsr(original map[string]RegsitryInfo) map[string]registryYAML {
	parserd := make(map[string]registryYAML, 0)
	for k, v := range original {
		parserd[k] = registryYAML{
			Home:     v.Home,
			Registry: v.Uri,
		}
	}
	return parserd
}

func NewRegistrySourceData(fs fs.FS) RegistryDataSource {
	return RegistryDataSource{
		fs:         fs,
		Registry:   make(map[string]string),
		Keys:       make([]string, 0),
		PresetKeys: make([]string, 0),
		Niave:      make(map[string]RegsitryInfo),
	}
}
