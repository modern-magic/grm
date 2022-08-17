package main

import (
	"fmt"
	"os"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/action"
	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/registry"
)

var helperInfo = `Usage: Grm [options] [command]

Options:
  -v, --version                           Output the version number
  -h, --help                              Output usage information

Commands:
  ls                                      List all the registries
  current                                 Show current registry name
  use  <name>                             Change registry to registry
  test <name>                             Test response time for specific or all registries
  add  <name> <registry> [home]           Add one custom registry
  del  <name>                             Delete one custom registry by alias
  help                                    Print this help
`

func Run() int {
	return runImpl(os.Args[1:])
}

func newRegistrySourceData() registry.RegistryDataSource {
	return registry.RegistryDataSource{
		Registry:   make(map[string]string),
		Keys:       make([]string, 0),
		PresetKeys: make([]string, 0),
	}
}

func runImpl(args []string) int {

	if len(args) == 0 {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s", c.Cyan, helperInfo, c.Reset)
		})
		return 0
	}

	for _, arg := range args {
		switch arg {
		case "-h", "--help", "help":
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("%s%s%s", c.Cyan, helperInfo, c.Reset)
			})
			return 0
		case "-v", "--version", "version":
			logger.Info(internal.StringJoin("[Grm]: version", grmVersion))
			return 0
		}
	}

	// initlize nrm & npm preset source.
	sources := newRegistrySourceData()
	return parserSourceForRun(args, &sources)
}

func parserSourceForRun(args []string, source *registry.RegistryDataSource) int {

	source.Keys = append(source.Keys, registry.GetPresetRegistryNames()...)
	source.PresetKeys = append(source.PresetKeys, source.Keys...)
	user := registry.NewUserResolver(&fs.FsImpl{})
	user.Resolve()
	for _, key := range registry.GetPresetRegistryNames() {
		source.Registry[key] = registry.GetPresetRegistryInfo(key)
	}

	source.Keys = append(source.Keys, user.GetNames()...)

	for _, key := range user.GetNames() {
		source.Registry[key] = user.GetRegistries()[key].Uri
	}

	for _, arg := range args {
		switch arg {
		case "ls":
			return action.ShowSources(source)
		case "current":
			return action.ShowCurrent()
		case "use":
			return action.SetCurrent(source, args[1:])
		case "add":
			return action.AddRegistry(source, args[1:])
		case "del":
			return action.DelRegistry(source, args[1:])
		case "test":
			return action.FetchRegistry(source, args[1:])
		}
	}
	return 0
}
