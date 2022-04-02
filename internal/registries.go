package internal

type RegistryMeta struct {
	Home     string
	Registry string
}

var (
	npm = RegistryMeta{
		Home:     "https://www.npmjs.org",
		Registry: "https://registry.npmjs.org/",
	}

	yarn = RegistryMeta{
		Home:     "https://yarnpkg.com",
		Registry: "https://registry.yarnpkg.com/",
	}

	huawei = RegistryMeta{
		Home:     "https://repo.huaweicloud.com/repository/npm/",
		Registry: "https://repo.huaweicloud.com/repository/npm/",
	}

	tencet = RegistryMeta{
		Home:     "https://mirrors.cloud.tencent.com/npm/",
		Registry: "https://mirrors.cloud.tencent.com/npm/",
	}

	cnpm = RegistryMeta{
		Home:     "https://cnpmjs.org",
		Registry: "https://r.cnpmjs.org/",
	}

	taobao = RegistryMeta{
		Home:     "https://npmmirror.com",
		Registry: "https://registry.npmmirror.com/",
	}

	npmMirror = RegistryMeta{
		Home:     "https://skimdb.npmjs.com/",
		Registry: "https://skimdb.npmjs.com/Registry/",
	}
)

var presetKeys = []string{
	"npm", "yarn", "huawei", "tencet", "cnpm", "taobao", "npmMirror",
}
var presetRegistries = []RegistryMeta{
	npm, yarn, huawei, tencet, cnpm, taobao, npmMirror,
}

type Registries struct {
	Registries        map[string]RegistryMeta
	RegistriesKeys    []string
	NrmRegistriesKeys []string
}

func (r *Registries) InitlizeRegistries() {
	for idx, v := range presetKeys {
		r.Registries[v] = presetRegistries[idx]
		r.RegistriesKeys = append(r.RegistriesKeys, v)
	}
	nrmRegistries, nrmRegistriesKey := getNrmRegistries()
	for idx, v := range nrmRegistriesKey {
		r.Registries[v] = nrmRegistries[idx]
		r.RegistriesKeys = append(r.RegistriesKeys, v)
		r.NrmRegistriesKeys = append(r.NrmRegistriesKeys, v)
	}
}

var Regis = &Registries{
	Registries:        make(map[string]RegistryMeta, 0),
	RegistriesKeys:    make([]string, 0),
	NrmRegistriesKeys: make([]string, 0),
}
