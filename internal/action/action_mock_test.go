package action

import (
	"testing"

	"github.com/modern-magic/grm/internal/registry"
)

func TestMockAction(t *testing.T) {
    action :=MockAction(registry.PresetRegistry)
    s :=  action.View(ViewOptions{
       All: true,
    })
    if s !=2 {
    	t.Fatal("should printf all source")
    }
    s = action.View(ViewOptions{
       All: false,
    })

    if s !=1{
        t.Fatal("should printf current source")
    }

}