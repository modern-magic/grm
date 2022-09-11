package action

import (
	"testing"

	"github.com/modern-magic/grm/internal/registry"
)

func TestMockActionView(t *testing.T) {
	act := MockAction(registry.PresetRegistry, []string{})
	s := act.View(ViewOptions{
		All: true,
	})

	if s != 2 {
		t.Errorf("should printf all source")
	}
	s = act.View(ViewOptions{})
	if s != 1 {
		t.Errorf("should printf current source")
	}
}

func TestMockActionDrop(t *testing.T) {
	act := MockAction(registry.PresetRegistry, []string{"test"})
	s := act.Drop()
	if s != 0 {
		t.Errorf("drop action failed")
	}
}

func TestMockActionJoin(t *testing.T) {
	act := MockAction(registry.PresetRegistry, []string{"test"})
	s := act.Join()
	if s != 0 {
		t.Errorf("join action failed")
	}
}

func TestMockActionUse(t *testing.T) {
	act := MockAction(registry.PresetRegistry, []string{"npm"})
	s := act.Use()
	if s != 0 {
		t.Errorf("use action failed")
	}
}
