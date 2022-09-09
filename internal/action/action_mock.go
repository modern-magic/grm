package action

import (
	"github.com/modern-magic/grm/internal/fs"
	"github.com/modern-magic/grm/internal/registry"
)

type mockAction struct {
	source map[string] registry.RegsitryInfo
	fs     fs.FS
}



func MockAction(source map[string] registry.RegsitryInfo, ) Action {
	return &mockAction{
		source: source,
        fs: fs.NewFS(),
	}
}


func (action *mockAction) View(option  ViewOptions) int {
    if(option.All){
        return 2
    }
	return 1
}
func (action *mockAction) Drop() int {

	return 1
}

func (action *mockAction) Join() int {
	return 1
}

func (action *mockAction) Test() int {
	return 1
}

func (action *mockAction) Use() int {
	return 1
}
