package gdp

import (
	"testing"
)

func TestTag(t *testing.T) {

	p := splitQuery("a")
	got := p[0]
	want := "a"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestSplitQuery(t *testing.T) {
	a := splitQueries(".test a, .test1 a")
	got := a[0]
	want := ".test a"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = a[1]
	want = ".test1 a"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
