# Grm

[![Go Reference](https://pkg.go.dev/badge/github.com/XeryYue/grm.svg)](https://pkg.go.dev/github.com/XeryYue/grm)

A npm registry manger.

Use smaller dependencies than [nrm](https://github.com/Pana/nrm).

### Install

-   Using [Go](https://golang.org/):

```shell

$ go install github.com/XeryYue/grm@latest

```

Or download a [binary package](https://github.com/XeryYue/grm/releases/latest).

### Usage

```shell
Usage: grm [options] [command]

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

### Q & A

> Why new nrm?

-   Installing `nrm` is too slow, and `grm` can be fast.

> Why not so comprehensive

-   Because we are only source manager.
