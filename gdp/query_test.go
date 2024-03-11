package gdp

import (
	"testing"
)

func TestQuery(t *testing.T) {

	p := splitQuery("a")
	got := p[0]
	want := "a"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	p = splitQuery("a.link")
	got = p[0]
	want = "a.link"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestSplitQueries(t *testing.T) {
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

func TestAttr(t *testing.T) {
	q := "a[href='test'][class='aa']"
	qa := QueryAttr{q, 0, len(q)}
	attrs := qa.parseAttr()

	got := attrs["tag"]
	want := "a"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = attrs["href"]
	want = "test"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = attrs["class"]
	want = "aa"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}