package internal_test

import (
	"testing"

	"github.com/modern-magic/grm/internal"
)

func TestFetch(t *testing.T) {
	var uri string
	var expected internal.FetchContext
	var actual internal.FetchContext

	uri = "https://npmmirror.com"
	expected = internal.FetchContext{
		Status:     "200 OK",
		StatusCode: 200,
		Time:       0.0,
		IsTimeout:  false,
	}
	actual = internal.Fetch(uri)
	// the `Time` field is not accurate, so we can't compare it
	if actual.StatusCode != expected.StatusCode {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}
