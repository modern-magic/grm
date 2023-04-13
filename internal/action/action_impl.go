package action

import (
	"fmt"
	"os"

	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/logger"
	"github.com/modern-magic/grm/internal/source"
)

var MaxTestSpeedLimit = 5

type actionImpl struct {
	fs   fs.FS
	args []string
	conf *source.GrmConfig
}

func NewAction(args []string) Action {
	action := &actionImpl{
		fs:   fs.NewFS(),
		args: args,
	}
	action.conf = source.NewGrmConf()
	return action
}

func (action *actionImpl) currentPath() string {
	return action.conf.GetCurrentPath()
}

func (action *actionImpl) View(option ViewOptions) int {

	current := action.currentPath()
	paths := action.conf.ListAllPath()

	alias := ""
	if s, ok := source.DefaultSource[current]; ok {
		alias = s.String()
	}

	if !option.All {
		logger.PrintTextWithColor(os.Stdout, func(c logger.Colors) string {
			return fmt.Sprintf("%s%s%s\n", c.Dim, alias, c.Reset)
		})
		return 0
	}

	for _, p := range paths {

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
	return 0
}

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
