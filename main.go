package main

import (
	"fmt"
	"gonrm/internal"
	"os"
)

var (
	version = "V0.0.0"
)

var helpText = func() string {

	return `Usage: gonrm [options] [command]

Options:
  -v, --version                           output the version number
  -h, --help                              output usage information

Commands:
  ls                                      List all the registries
  current                                 Show current registry name
  use <registry>                          Change registry to registry
  add <name> <registry> [home]             Add one custom registry
  del <name>                              Delete one custom registry by alias
  help                                    Print this help
`

}

func main() {
	osArgs := os.Args[1:]
	registries := internal.Regis
	registries.InitlizeRegistries()
	for _, arg := range osArgs {
		switch {
		case arg == "-h", arg == "--help", arg == "help":
			fmt.Printf(internal.AnsiColor.Color(internal.TipColor), helpText())
			os.Exit(0)
		case arg == "-v", arg == "--version":
			fmt.Printf("gonrm: %s", version)
			os.Exit(0)
		case arg == "ls":
			internal.ShowRegistries(registries)
			os.Exit(0)
		case arg == "current":
			internal.ShowCurrentRegistry()
			os.Exit(0)
		case arg == "use":
			internal.SetUsageRegistry(osArgs[1:], registries)
			os.Exit(0)
		case arg == "add":
			internal.AddRegistry(osArgs[1:])
			os.Exit(0)
		case arg == "del":
			internal.DelRegistry(osArgs[1:], registries.NrmRegistriesKeys)
			os.Exit(0)
		}
	}

}
