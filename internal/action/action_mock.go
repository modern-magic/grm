package action

import (
	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/registry"
)

type mockAction struct {
	source map[string]registry.RegsitryInfo
	fs     fs.FS
	args   []string
}

func MockAction(source map[string]registry.RegsitryInfo, args []string) Action {
	return &mockAction{
		source: source,
		fs:     fs.NewFS(),
		args:   args,
	}
}

func (action *mockAction) View(option ViewOptions) int {
	if option.All {
		return 2
	}
	return 1
}
func (action *mockAction) Drop() int {
	if _, ok := action.source[action.args[0]]; ok {
		return 0
	}
	return 1
}

func (action *mockAction) Join() int {
	if _, ok := action.source[action.args[0]]; ok {
		return 1
	}
	return 0
}

func (action *mockAction) Test() int {
	return 1
}

func (action *mockAction) Use() int {
	if _, ok := action.source[action.args[0]]; ok {
		return 0
	}
	return 1
}
