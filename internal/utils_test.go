package internal_test

import (
	"testing"

	"github.com/modern-magic/grm/internal"
)

func TestStringJoin(t *testing.T) {
	var s []string
	var expected string
	var actual string

	s = []string{"a", "b", "c"}
	expected = "a b c"
	actual = internal.StringJoin(s...)
	if actual != expected {
		t.Errorf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestIsUri(t *testing.T) {
	var uri string
	var expected bool
	var actual bool

	uri = "http://www.example.com"
	expected = true
	actual = internal.IsUri(uri)
	if actual != expected {
		t.Errorf("Expected: %t, Actual: %t", expected, actual)
	}

}

func TestPickArgs(t *testing.T) {
	args := []string{"a", "b", "c"}
	cursor := 1
	v := internal.PickArgs(args, cursor)
	if v != "b" {
		t.Errorf("Expected: %s, Actual: %s", "b", v)
	}

}
