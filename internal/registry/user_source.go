package registry

import (
	"os"

	"github.com/modern-magic/grm/internal/fs"
)

type registryYAML struct {
	Home     string `yaml:"home"`
	Registry string `yaml:"registry"`
}

type RegistryYAML = registryYAML

func GetSystemPreffix() string {
	if fs.IsWindows() {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

type ResolverResult struct {
	Registries map[string]RegsitryInfo
	Names      []string
}

type Resolver interface {
	GetNames() []string
	GetRegistries() map[string]RegsitryInfo
	Resolve()
}

type resolver struct {
	fs       fs.FS
	registry map[string]RegsitryInfo
	names    []string
}

func NewUserResolver(fs fs.FS) Resolver {
	return &resolver{
		fs:       fs,
		registry: make(map[string]RegsitryInfo, 0),
		names:    make([]string, 0),
	}
}

func (r *resolver) Resolve() {
	content, err, _ := r.fs.ReadYAML(Grmrc, map[string]registryYAML{})
	if err != nil {
		// user may delte .grmrc.yaml.So we should set it as a empty map
		r.registry = make(map[string]RegsitryInfo, 0)
	}
	switch c := content.(type) {
	case map[string]registryYAML:
		for k, v := range c {
			r.registry[k] = RegsitryInfo{
				Home: v.Home,
				Uri:  v.Registry,
			}
			r.names = append(r.names, k)
		}
	default:
		panic("Invalid type")
	}
}

func (r *resolver) GetRegistries() map[string]RegsitryInfo {
	return r.registry
}

func (r *resolver) GetNames() []string {
	return r.names
}
