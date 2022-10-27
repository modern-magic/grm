## Grm Contributing Guide

### Ready to start

We welcome everyone to join in the construction of the project.For basic operation of Git,you can refer to [GitHub's help documentation](https://help.github.com/en/github/using-git).

1. [Fork this repository](https://help.github.com/en/github/getting-started-with-github/fork-a-repo) to your own account and then clone it.
2. Create a new branch for your changes: `git checkout -b {BRANCH_NAME}`.
3. Donwload go mod.
4. Develop on your dev tools

At any time, you think it's ok, you can start the following steps to submit your amazing works:

We are using `golangci-lint`, So you need change your lint utils in your local dev tools.

1. Run `git commit -m '{YOUR_MESSAGE}'` to commit changes. Commit info should be formatted by the [rules](https://github.com/conventional-changelog/commitlint/blob/master/%40commitlint/config-conventional/README.md).
2. Push code to your own repo and [create PullRequest](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests) at GitHub.

### Q & A

> Are we have project style?

- Yes, We follow the project style of golang for development, If you don't know you can view [style Guide](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)

> How can I build and test different platform of a package locally?

- We offer a Makefile to build it. You can view the different platform package after you build it.

> How to use Meakfile in windows system?

- In windows, you need install `GCC` to your local environment. And then find you installed directory and can fork `mingw32-make.exe` and rename it as `make.exe`. Then you can use make command.

> Why my Makefile execution error

- Before you run `make` in your shell. You need check your local environment already has `upx`. If you're winodows user you should need install extra [wixtools](https://wixtoolset.org/).

> Why i use makefile can't generator msi?

- Unfortunately, We try to use wix in linux system. But it can't support the new syntax. So we only allowed build msi with windows system.
