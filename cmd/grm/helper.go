package main

import (
	"fmt"
	"os"

	"github.com/modern-magic/grm/internal/action"
	"github.com/modern-magic/grm/internal/logger"
)

var helperText = `
` + `Running version ` + grmVersion + `.` + `

Usage:

    grm current                            : Display active registry name.
    grm use <name>                         : Use this registry.
    grm list                               : List the registry info. Aliased as ls.
    grm test [name]                        : Test response time for specific or all registries.
    grm add <name> <registry>              : Add one custom registry.
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
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("%s%s%s\n", c.Blue, grmVersion, c.Reset)
			})
			return 0
		}
	}

	return parserSourceForRun(args[0], args[1:])
}

func parserSourceForRun(command string, args []string) int {

	act := action.NewAction(args)
	switch command {
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
	return 0
}
