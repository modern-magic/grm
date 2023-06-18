package action

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/net"
	"github.com/modern-magic/grm/internal/shell"
	"github.com/modern-magic/grm/internal/source"
)

func printAlias(s string) source.S {
	return source.EnsureDefaultKey(s)
}

func isDefaultAlias(s string) (bool, string) {
	alias := printAlias(s)
	return alias != source.System, alias.String()
}

func isDefaultPath(s string) (bool, string) {
	alias, ok := source.DefaultSource[s]
	if ok {
		return true, alias.String()
	}
	return false, source.System.String()
}

func isDefault(i string, isPath bool) (s bool, alias, source string) {
	if isPath {
		s, alias = isDefaultPath(i)
		return s, alias, i
	}
	s, alias = isDefaultAlias(i)
	return s, alias, i
}

func isURL(p string) bool {
	_, err := url.Parse(p)
	return err == nil
}

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

func (action *actionImpl) currentPath() string {
	return action.conf.GetCurrentPath()
}

func (action *actionImpl) View(option ViewOptions) int {
	cp := action.currentPath()
	ok, alias, source := isDefault(cp, true)
	if !ok {
		sources, _ := action.conf.ScannerUserConf()
		if s, ok := sources[source]; ok {
			alias = s
		}
	}

	if !option.All {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Dim, alias, c.Reset)
		})
		return 0
	}

	for _, s := range action.conf.Paths {
		if s == alias {
			logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
				return fmt.Sprintf("* %s%s%s%s\n", c.Cyan, s, fmt.Sprintf(" %s%s%s", c.DimCyan, "default", c.Reset), c.Reset)
			})
			continue
		}
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("* %s%s%s\n", c.Dim, s, c.Reset)
		})
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
	name := action.args[0]
	ok, _, alias := isDefault(name, false)
	if ok {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s%s\n", c.Red, "error: can't remove default source", alias, c.Reset)
		})
		return 1
	}

	// user conf
	_, keys := action.conf.ScannerUserConf()
	if _, ok := keys[alias]; !ok {
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
	// // Set npm as default choose
	ok = action.conf.SetCurrentPath(source.DefaultKey[source.Npm])
	if ok {
		err := action.fs.OuputFile(action.conf.ConfPath, []byte(action.conf.GetCurrentConf()))
		if err != nil {
			logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
				return fmt.Sprintf("%s%s%s\n", c.Red, err, c.Reset)
			})
			return 1
		}
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
	name := action.args[0]
	path := action.args[1]
	ok, alias, _ := isDefault(name, false)
	if ok {
		logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Red, "error: can't be named the same as default", c.Reset)
		})
		return 1
	}
	// if
	// _, key := action.conf.ScannerUserConf()

	// _, userKey := action.conf.ScannerUserConf()
	// if _, ok := userKey[alias]; ok {
	// 	if !shell.MakeConfirm("The alias already exists. Do you want to modify it?") {
	// 		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
	// 			return fmt.Sprintf("%s%s%s\n", c.Dim, "process exit", c.Reset)
	// 		})
	// 		return 0
	// 	}
	// }
	// // verify path is a right url.
	// if !verifyURL(action.args[1]) {
	// 	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
	// 		return fmt.Sprintf("%s%s%s\n", c.Red, "invalid url", c.Reset)
	// 	})
	// 	return 1
	// }

	// file := path.Join(action.conf.BaseDir, alias)
	// err := action.fs.OuputFile(file, []byte(action.args[1]))
	// if err != nil {
	// 	logger.PrintTextWithColor(os.Stderr, func(c logger.Colors) string {
	// 		return fmt.Sprintf("%s%s%s\n", c.Red, err, c.Reset)
	// 	})
	// 	return 1
	// }
	// logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
	// 	return fmt.Sprintf("%s%s%s\n", c.Green, "update new conf success", c.Reset)
	// })
	// return 0
}

// Pick up should test path
func (action *actionImpl) Test() int {

	aliases := action.args
	_, keys := action.conf.ScannerUserConf()
	urls := make(map[string]string, len(keys)+len(source.DefaultSource))
	paths := action.conf.MergePaths(keys)

	if len(aliases) >= 1 {
		for _, alias := range aliases {
			if v, ok := paths[alias]; ok {
				urls[alias] = v
			}
		}
	} else {
		urls = paths
	}

	current := action.currentPath()
	net.MakeRequest(urls, func(message net.RequestMessage) {
		if message.Err != nil {
			fmt.Fprintf(os.Stderr, "failed with: %v\n", message.Err)
			return
		}
		var isDefault bool
		if message.Path == current {
			isDefault = true
		}

		if isDefault {
			message.Alias = fmt.Sprintf("%s%s%s", logger.TerminalColors.Cyan, message.Alias, logger.TerminalColors.Reset)
		}
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("* %s%s (%s)%s\n", c.Dim, message.Alias, message.Sec, c.Reset)
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
