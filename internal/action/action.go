package action

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/net"
	"github.com/modern-magic/grm/internal/shell"
	"github.com/modern-magic/grm/internal/source"
)

type ViewOptions struct {
	All bool
}

type actionImpl struct {
	fs   fs.FS
	args []string
	conf *source.GrmConfig
}

func NewAction(args []string) *actionImpl {
	action := &actionImpl{
		fs:   fs.NewFS(),
		args: args,
	}
	action.conf = source.NewGrmConf()
	action.conf.ListAllPath()
	return action
}

func verifyURL(s string) bool {
	_, err := url.Parse(s)
	return err == nil
}

func (action *actionImpl) currentPath() string {
	return action.conf.GetCurrentPath()
}

func (action *actionImpl) isDefaultAlias(alias source.S) bool {
	return alias != source.System
}

func (action *actionImpl) isAliasExists(alias string) bool {
	_, userKey := action.conf.ScannerUserConf()
	_, ok := userKey[alias]
	return ok
}

func (action *actionImpl) View(option ViewOptions) int {
	current := action.currentPath()
	alias := ""
	if s, ok := source.DefaultSource[current]; ok {
		alias = s.String()
	} else {
		userSource, _ := action.conf.ScannerUserConf()
		if s, ok := userSource[current]; ok {
			alias = s
		}
	}

	if !option.All {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Dim, alias, c.Reset)
		})
		return 0
	}

	for _, p := range action.conf.Paths {

		if p == alias {
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("* %s%s%s%s\n", c.Cyan, p, fmt.Sprintf(" %s%s%s", c.DimCyan, "default", c.Reset), c.Reset)
			})
		} else {
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("* %s%s%s\n", c.Dim, p, c.Reset)
			})
		}
	}

	return 0
}

func (action *actionImpl) Drop() int {
	if len(action.args) < 1 {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, "error: alias should be passed", c.Reset)
		})
		return 1
	}
	alias := action.args[0]
	s := source.EnsureDefaultKey(alias)
	if action.isDefaultAlias(s) {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, "error: can't remove default source", c.Reset)
		})
		return 1
	}

	if !action.isAliasExists(alias) {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, "error: can't found alias", c.Reset)
		})
		return 1
	}
	if !shell.MakeConfirm("Are you sure to remove the registry?") {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Dim, "process exit", c.Reset)
		})
		return 0
	}
	err := action.fs.Rm(path.Join(action.conf.BaseDir, alias))
	if err != nil {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, err, c.Reset)
		})
		return 1
	}
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Green, "remove registry success", c.Reset)
	})
	return 0
}

func (action *actionImpl) Join() int {
	if len(action.args) < 2 {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, "error: registry should be passed", c.Reset)
		})
		return 1
	}
	alias := action.args[0]
	s := source.EnsureDefaultKey(alias)
	if action.isDefaultAlias(s) {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, "error: can't be named the same as default", c.Reset)
		})
		return 1
	}

	_, userKey := action.conf.ScannerUserConf()
	if _, ok := userKey[alias]; ok {
		if !shell.MakeConfirm("The alias already exists. Do you want to modify it?") {
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("%s%s%s\n", c.Dim, "process exit", c.Reset)
			})
			return 0
		}
	}
	// verify path is a right url.

	if !verifyURL(action.args[1]) {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Dim, "invalid url", c.Reset)
		})
		return 1
	}

	file := path.Join(action.conf.BaseDir, alias)
	err := action.fs.OuputFile(file, []byte(action.args[1]))
	if err != nil {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, err, c.Reset)
		})
		return 1
	}
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Green, "update new conf success", c.Reset)
	})
	return 0
}

func (action *actionImpl) Test() int {

	userSoure, _ := action.conf.ScannerUserConf()
	urls := make([]string, 0, len(userSoure)+len(source.DefaultSource))
	for k := range source.DefaultSource {
		urls = append(urls, k)

	}
	for k := range userSoure {
		urls = append(urls, k)
	}

	current := action.currentPath()
	net.MakeRequest(urls, func(r string) {
		expr := strings.Split(r, "->") // url msg
		url := expr[0]
		msg := expr[1]
		name := ""
		var isDefault bool
		if current == url {
			isDefault = true
		}
		if alias, ok := source.DefaultSource[url]; ok {
			name = alias.String()
		} else {
			if alias, ok := userSoure[url]; ok {
				name = alias
			}
		}
		if isDefault {
			name = fmt.Sprintf("%s%s%s", logger.TerminalColors.Cyan, name, logger.TerminalColors.Reset)
		}
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("* %s%s%s%s\n", c.Dim, name, msg, c.Reset)
		})

	})
	return 0
}

func (action *actionImpl) Use() int {
	alias := action.args[0]
	s := source.EnsureDefaultKey(alias)
	url := ""
	if !action.isDefaultAlias(s) {
		_, userKey := action.conf.ScannerUserConf()
		if v, ok := userKey[alias]; ok {
			url = v
		} else {
			if !shell.MakeConfirm("This registry can't find. Do you want to add a new one?") {
				logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
					return fmt.Sprintf("%s%s%s\n", c.Dim, "process exit", c.Reset)
				})
				return 0
			}
			url, err := shell.MakePrompt(fmt.Sprintf("Enter registry address for %s: ", alias), func(input string) error {
				if len(input) == 0 {
					return errors.New("can't be empty")
				}
				if !verifyURL(input) {
					return errors.New("invalid url")
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update config file: %v\n", err)
				return 1
			}
			file := path.Join(action.conf.BaseDir, alias)
			err = action.fs.OuputFile(file, []byte(url))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update config file: %v\n", err)
				return 1
			}
		}
	} else {
		url = source.DefaultKey[s]
	}
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s%s%s%s\n", c.Dim, "using registry", fmt.Sprintf(" %s%s%s", c.Green, alias, c.Reset), c.Reset)
	})
	ok := action.conf.SetCurrentPath(url)
	if ok {
		err := action.fs.OuputFile(action.conf.ConfPath, []byte(action.conf.GetCurrentConf()))
		if err == nil {
			return 0
		}
	}
	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s%s%s\n", c.Red, "invalid error", c.Reset)
	})
	return 1

}
