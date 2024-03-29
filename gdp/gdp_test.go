package gdp

import (
	"testing"
)

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		html := `<div class="test">this is for benchmark</div>`
		Default(html)
	}
}

func TestTagWithoutClosing(t *testing.T) {
	html := `<div class="test">this is for test`
	document := Default(html)

	got := document.Html()
	want := `<div class="test">this is for test</div>`

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestContent(t *testing.T) {
	html := `<div class="test">this < is for test</div>`
	document := Default(html)

	got := document.Html()

	want := `<div class="test">this < is for test</div>`

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestNoEndTag(t *testing.T) {
	html := `<div class="test"><br /> Test</div>`
	document := Default(html)

	got := document.Html()
	want := `<div class="test"><br /> Test</div>`

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
