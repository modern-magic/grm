package main

import (
	"fmt"
	"os"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/action"
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

	source := registry.NewSource()
	return parserSourceForRun(args, source)
}

func parserSourceForRun(args []string, source registry.Source) int {

	act := action.NewAction(action.ActionOptions{
		Source:     source.GetSource(),
		UserSource: source.GetUserSource(),
		Args:       args,
	})

	for _, arg := range args {
		switch arg {
		case "ls":
			return act.View(action.ViewOptions{All: true})
		case "current":
			return act.View(action.ViewOptions{All: false})
		case "use":
			return act.Use()
		case "add":
			return act.Join()
		case "del":
			return act.Drop()
		case "test":
			return act.Test()
		}
	}
	return 0
}
