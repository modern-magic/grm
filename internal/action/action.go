package action

import (
	"fmt"
	"os"
	"strings"

	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/net"
	"github.com/modern-magic/grm/internal/source"
)

var MaxTestSpeedLimit = 5

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

	current := action.currentPath()

	alias := ""
	if s, ok := source.DefaultSource[current]; ok {
		alias = s.String()
	} else {
		userSource, userKey := action.conf.ScannerUserConf()
		if s, ok := userSource[current]; ok {
			alias = userKey[s]
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
	return 0

}

func (action *actionImpl) Join() int {
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

// Todo
func (action *actionImpl) Use() int {

	s := source.EnsureDefaultKey(action.args[1])
	if s == source.System {
		return 0
	}
	path := source.DefaultKey[s]

	logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
		return fmt.Sprintf("%s%s%s%s\n", c.Dim, "Using registry", fmt.Sprintf("%s%s%s", c.Green, s, c.Reset), c.Reset)
	})
	ok := action.conf.SetCurrentPath(path)
	if ok {
		err := action.fs.OuputFile(action.conf.ConfPath, []byte(action.conf.GetCurrentConf()))
		if err == nil {
			return 0
		}
	}
	return 1

}
