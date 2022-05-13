# Grm

<p align="center">
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/modern-magic/grm"><img src="https://goreportcard.com/badge/github.com/modern-magic/grm?style=flat-square" /></a>
<a title="Doc for grm" target="_blank" href="https://pkg.go.dev/github.com/modern-magic/grm"><img src="https://pkg.go.dev/badge/github.com/modern-magic/grm.svg" /></a>
<a title="Codecov" target="_blank" href="https://codecov.io/gh/modern-magic/grm"><img src="https://img.shields.io/codecov/c/github/modern-magic/grm?style=flat-square&logo=codecov" /></a>
<a title="Release" target="_blank" href="https://github.com/modern-magic/grm/releases"><img src="https://img.shields.io/github/v/release/modern-magic/grm.svg?color=161823&style=flat-square&logo=smartthings" /></a>
</p>

A npm registry manger.

Use smaller dependencies than [nrm](https://github.com/Pana/nrm).

## Install

Using [homebrew](https://brew.sh/)

```bash

$ brew install modern-magic/tap/grm

```

Using [Go](https://golang.org/):

```bash

$ go install github.com/modern-magic/grm/cmd/grm@latest

```

Or you can download a [binary package or msi package](https://github.com/modern-magic/grm/releases/latest).

## Usage

```shell
Usage: Grm [options] [command]

Options:
  -v, --version                           output the version number
  -h, --help                              output usage information

Commands:
  ls                                      List all the registries
  current                                 Show current registry name
  use <registry>                          Change registry to registry
  add <name> <registry> [home]            Add one custom registry
  del <name>                              Delete one custom registry by alias
  test [name]                             Test registry or registries speed by alias
  help                                    Print this help

```

## Q & A

> What is grm?

A minimalist npm source manager.

> Why new nrm?

Installing `nrm` is too slow, and `grm` can be fast.

> Why not so comprehensive

Because we are only source manager.

> Why choose grm?

Compress to `nrm` and `nnrm` . `grm` has more advantages in installing. Benefit from golang cross platform.

> Why do I get a virus report when I use the msi installation package?

We use upx to pack grm. So in some antivirus software will report it have a virus. But you can use it with confidence. grm is a poison free program :)

## All Contributors

Thanks to the following friends for their contributions to Grm:

<a href="https://github.com/modern-magic/grm/graphs/contributors">
  <img src="https://opencollective.com/grm/contributors.svg?width=890&button=false" alt="contributors">
</a>

## Acknowledgements

Thanks to [JetBrains](https://www.jetbrains.com/) for allocating free open-source licences for IDEs.

<p align="right">
<img width="250px" height="250px"  src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand logo.">
</p>

## LICENSE

[MIT](./LICENSE)
