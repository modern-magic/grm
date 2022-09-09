package action

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/registry"
)

type actionImpl struct {
	fs         fs.FS
	source     map[string]registry.RegsitryInfo
	userSource map[string]registry.RegsitryInfo
	args       []string
}

type ActionOptions struct {
	Source     map[string]registry.RegsitryInfo
	UserSource map[string]registry.RegsitryInfo
	Args       []string
}

func NewAction(options ActionOptions) Action {
	return &actionImpl{
		fs:         fs.NewFS(),
		source:     options.Source,
		userSource: options.UserSource,
		args:       options.Args,
	}
}

func loadRegistryFromSource(source map[string]registry.RegsitryInfo, args []string, invork func(r *RegistryDataSource) int) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	key := internal.PickArgs(args, 1)
	v, ok := source[key]

	if !ok {
		logger.Error(internal.StringJoin("[Grm]: Can't find alias", key, "in your .grmrc file. Please check it exist."))
		return 1
	}
	return invork(&RegistryDataSource{
		Name:     key,
		Uri:      v.Uri,
		Internal: v.Internal,
	})
}

func getDashLine(key string, distance int) string {
	final := math.Max(2, (float64(distance) - float64(len(key))))
	bar := make([]string, int(final))
	for i := range bar {
		bar[i] = "-"
	}
	return strings.Join(bar[:], "-")
}

func (action *actionImpl) getCurrent() string {
	s, _ := registry.ReadNpm()
	s = strings.Replace(s, "", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return s
}

func (action *actionImpl) View(option ViewOptions) int {

	cur := action.getCurrent()

	if !option.All {
		logger.Info(internal.StringJoin("[Grm]: you're using", cur))
		return 0
	}

	serialize := make([]RegistryDataSource, 0, len(action.source))

	for k, v := range action.source {
		m := RegistryDataSource{
			Name: k,
			Uri:  v.Uri,
		}
		if cur != v.Uri {
			serialize = append(serialize, m)
			continue
		}
		serialize = append([]RegistryDataSource{m}, serialize...)

	}

	for i, v := range serialize {

		orginalLog := internal.StringJoin(v.Name, getDashLine(v.Name, len(serialize)), v.Uri)
		if i == 0 {
			logger.Success(internal.StringJoin("*", orginalLog))
			continue
		}
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Dim, orginalLog, c.Reset)
		})
	}

	return 0
}

func (action *actionImpl) Drop() int {

	return loadRegistryFromSource(action.source, action.args, func(r *RegistryDataSource) int {
		if r.Internal {
			logger.Error(internal.StringJoin("[Grm]: can't delete preset registry", r.Name))
			return 1
		}
		delete(action.userSource, r.Name)
		raw := registry.Parsr(action.source)
		err := action.fs.WriteYAML(registry.Grmrc, raw)
		if err != nil {
			logger.Error(internal.StringJoin("[Grm]: del registry fail", err.Error()))
			return 1
		}
		logger.Success(internal.StringJoin("[Grm]: del registry", r.Name, "success!"))
		return 0
	})

}

func (action *actionImpl) Join() int {

	defer func() {
		if err := recover(); err != nil {
			logger.Warn(internal.StringJoin("[Grm]: Please pass the registry url assets."))
			return
		}
	}()
	if len(action.args) < 1 {
		logger.Warn(internal.StringJoin("[Grm]: Please pass an alias."))
		return 1
	}
	k := internal.PickArgs(action.args, 1)
	uri := internal.PickArgs(action.args, 2)
	home := uri
	if len(action.args) > 3 {
		home = internal.PickArgs(action.args, 4)
	}
	real := []string{uri, home}

	for i := 0; i < len(real); {
		if !internal.IsUri(real[i]) {
			logger.Error("[Grm]: Please verify the uri address you entered.")
			return 1
		}

		i++
	}

	if _, ok := action.source[k]; ok {
		logger.Error(internal.StringJoin("[Grm]:", k, "has already exist."))
		return 1
	}

	action.userSource[k] = registry.RegsitryInfo{
		Home: home,
		Uri:  uri,
	}

	raw := registry.Parsr(action.userSource)
	err := action.fs.WriteYAML(registry.Grmrc, raw)
	if err != nil {
		logger.Error(internal.StringJoin("[Grm]: Add registry fail", err.Error()))
		return 1
	}
	logger.Success(internal.StringJoin("[Grm]: Add registry success!"))
	return 0
}

func (action *actionImpl) Test() int {
	return 1
}

func (action *actionImpl) Use() int {
	return loadRegistryFromSource(action.source, action.args, func(r *RegistryDataSource) int {
		err := registry.WriteNpm(r.Uri)
		if err != nil {
			logger.Error(internal.StringJoin("[Grm]: use registry fail", err.Error()))
			return 1
		}
		logger.Success(internal.StringJoin("[Grm]: use", r.Name, "success!"))
		return 0
	})

}
