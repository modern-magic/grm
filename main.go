package main

import (
	"fmt"
	"nrm/internal"
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
  add <registry> <url> [home]             Add one custom registry
  set-auth [options] <registry> [value]   Set authorize information for a custom registry with a base64 encoded string or username and pasword
  set-email <registry> <value>            Set email for a custom registry
  set-hosted-repo <registry> <value>      Set hosted npm repository for a custom registry to publish packages
  del <registry>                          Delete one custom registry
  home <registry> [browser]               Open the homepage of registry with optional browser
  publish [options] [<tarball>|<folder>]  Publish package to current registry if current registry is a custom registry.
   if you're not using custom registry, this command will run npm publish directly
  test [registry]                         Show response time for specific or all registries
  help                                    Print this help
`

}

func main() {
	osArgs := os.Args[1:]
	for _, arg := range osArgs {
		switch {
		case arg == "-h", arg == "--help":
			fmt.Printf(helpText())
			os.Exit(0)
		case arg == "-v", arg == "--version":
			fmt.Printf("gonrm %s \n", version)
			os.Exit(0)
		case arg == "ls":
			internal.ShowList()
			os.Exit(0)
		}
	}
}
