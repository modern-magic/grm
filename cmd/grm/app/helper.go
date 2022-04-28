package app

import (
	"fmt"
	"os"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/logger"
)

var (
	version = "V0.4.1"
)

var helperInfo = `Usage: Grm [options] [command]

Options:
  -v, --version                           output the version number
  -h, --help                              output usage information

Commands:
  ls                                      List all the registries
  current                                 Show current registry name
  use <registry>                          Change registry to registry
  test <name>                             Test response time for specific or all registries
  add <name> <registry> [home]            Add one custom registry
  del <name>                              Delete one custom registry by alias
  help                                    Print this help
`

func Run() {
	osArgs := os.Args[1:]
	registries := internal.Regis
	registries.InitlizeRegistries()
	if len(osArgs) == 0 {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s", c.Cyan, helperInfo, c.Reset)
		})
		return
	}

	for _, arg := range osArgs {
		switch arg {
		case "-h", "--help", "help":
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("%s%s%s", c.Cyan, helperInfo, c.Reset)
			})
		case "-v", "--version", "version":
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("%s[Grm]: %s%s", c.Green, version, c.Reset)
			})
		case "ls":
			internal.ShowRegistries(registries)
		case "current":
			internal.ShowCurrentRegistry()
		case "use":
			internal.SetUsageRegistry(osArgs[1:], registries)
		case "add":
			internal.AddRegistry(osArgs[1:])
		case "del":
			internal.DelRegistry(osArgs[1:], registries.NrmRegistriesKeys)
		case "test":
			internal.CurlRegistry(osArgs[1:], registries)
		}
	}

}
