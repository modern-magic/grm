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
		switch arg {
		case "-h", "--help", "help":
			fmt.Printf(internal.AnsiColor.Color(internal.TipColor), helpText())
		case "-v", "--version", "version":
			fmt.Printf("grm: %s", version)
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
		}

	}
	os.Exit(0)

}
