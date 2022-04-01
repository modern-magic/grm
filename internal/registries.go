package internal

type RegistryMeta struct {
	Home     string
	Registry string
}

var npm = RegistryMeta{
	Home:     "https://www.npmjs.org",
	Registry: "https://registry.npmjs.org/",
}

var yarn = RegistryMeta{
	Home:     "https://yarnpkg.com",
	Registry: "https://registry.yarnpkg.com/",
}

var tencet = RegistryMeta{
	Home:     "https://mirrors.cloud.tencent.com/npm/",
	Registry: "https://mirrors.cloud.tencent.com/npm/",
}
var cnpm = RegistryMeta{
	Home:     "https://cnpmjs.org",
	Registry: "https://r.cnpmjs.org/",
}

var taobao = RegistryMeta{
	Home:     "https://npmmirror.com",
	Registry: "https://registry.npmmirror.com/",
}

var npmMirror = RegistryMeta{
	Home:     "https://skimdb.npmjs.com/",
	Registry: "https://skimdb.npmjs.com/Registry/",
}

var presetKeys = []string{
	"npm", "yarn", "tencet", "cnpm", "taobao", "npmMirror",
}
var presetRegistries = []RegistryMeta{
	npm, yarn, tencet, cnpm, taobao, npmMirror,
}

type Registries struct {
	Registries        map[string]RegistryMeta
	RegistriesKeys    []string
	NrmRegistriesKeys []string
}

// func (r *Registries) SetRegistries(source RegistryMeta, key string) {
// 	r.Registries[key] = source
// 	r.RegistriesKeys = append(r.RegistriesKeys, key)
// }

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
