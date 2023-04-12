package main

import (
	"os"

	"github.com/modern-magic/grm/internal"
	"github.com/modern-magic/grm/internal/action"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/registry"
)

var helperText = `
` + `Running version ` + grmVersion + `.` + `

Usage:

    grm current                            : Display active registry name.
    grm use <name>                         : Use this registry.
    grm list                               : List the registry info. Aliased as ls.
    grm test [name]                        : Test response time for specific or all registries.
    grm add <name> <registry> [home]       : Add one custom registry.
    grm del <name>                         : Delete one custom registry by alias.
    grm version                            : Displays the current running version of grm. Aliased as v.
`

func Run() int {
	return runImpl(os.Args[1:])
}

func runImpl(args []string) int {

	if len(args) == 0 {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return helperText
		})
		return 0
	}

	for _, arg := range args {
		switch arg {
		case "-h", "--help", "help":
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return helperText
			})
			return 0
		case "-v", "--version", "version":
			logger.Success(internal.StringJoin("[Grm]: version", grmVersion))
			return 0
		}
	}

	source := registry.NewSource()
	return parserSourceForRun(args, source)
}

func parserSourceForRun(args []string, source registry.Source) int {

	act := action.NewAction(args)

	for _, arg := range args {
		switch arg {
		case "ls", "list":
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
