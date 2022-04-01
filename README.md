# GoNrm

### Install

```shell

$ git clone git@github.com:XeryYue/go-nrm.git

$ go get

$ go build

$ go install

```

### Usage

```shell
Usage: gonrm [options] [command]

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

```
