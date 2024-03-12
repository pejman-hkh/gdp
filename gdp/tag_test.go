package gdp

import (
	"testing"
)

func TestRemove(t *testing.T) {
	document := Default(`<div id="test">test</div><div id="test1">test1</div>`)
	document.Find("#test").Eq(0).Remove()
	got := document.Html()
	want := `<div id="test1">test1</div>`
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestSib(t *testing.T) {
	document := Default(`<div class="parent"><div class="prev">test</div><div class="middle" id="middle">test1</div><span class="next"></span></div>`)

	middle := document.GetElementById("middle")
	got := middle.Parent().Attr("class")
	want := "parent"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
	got = middle.Prev().Attr("class")
	want = "prev"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = middle.Next().Attr("class")
	want = "next"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
