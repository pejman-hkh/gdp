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

func TestHtml(t *testing.T) {
	document := Default(`<div class="parent"><div class="prev">test</div><div class="middle" id="middle">test1</div><span class="next"></span></div>`)
	tag := document.Find(".parent").Eq(0)
	tag.SetHtml("<span>changed html</span>")
	got := tag.Html()
	want := `<span>changed html</span>`
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestFindFirst(t *testing.T) {
	document := Default(`<span>first</span><span>second</span><span>third</span>`)
	first := document.Find("span:first")

	got1 := len(first.list)
	want1 := 1
	if got1 != want1 {
		t.Errorf("got %q, wanted %q", got1, want1)
	}

	got := first.Eq(0).Html()
	want := "first"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestFindLast(t *testing.T) {
	document := Default(`<span>first</span><span>second</span><span>third</span>`)
	last := document.Find("span:last")

	got1 := len(last.list)
	want1 := 1
	if got1 != want1 {
		t.Errorf("got %q, wanted %q", got1, want1)
	}

	got := last.Eq(0).Html()
	want := "third"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestEq(t *testing.T) {
	document := Default(`<span>first</span><span>second</span><span>third</span>`)
	eq := document.Find("span:eq(1)")

	got := eq.Eq(0).Html()
	want := "second"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}
