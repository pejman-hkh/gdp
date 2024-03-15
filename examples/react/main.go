package main

import (
	"fmt"
	"reflect"

	"github.com/pejman-hkh/gdp/gdp"
)

type React struct{}

func (r React) Link(props map[string]*string, childrens string) string {
	return Render(`<a href="` + *props["to"] + `">` + childrens + `</a>`)
}

func (r React) Header(props map[string]*string, childrens string) string {
	return Render(`<header><nav><ul><li><a href="/">Home</a></li></ul></nav>` + childrens + `</header><h1>` + *props["title"] + `</h1>`)
}

func (r React) Footer(props map[string]*string, childrens string) string {
	return Render(`<footer><Link to="https://www.github.com/pejman-hkh/gdp">https://www.github.com/pejman-hkh/gdp</Link> ` + childrens + `</footer>`)
}

func isUpperCase(c byte) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	}
	return false
}

func Invoke(obj any, name string, args ...any) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	return reflect.ValueOf(obj).MethodByName(name).Call(inputs)
}

func Render(html string) string {
	document := gdp.Default(html)
	ret := ``
	document.Children().Each(func(i int, child *gdp.Tag) {
		tagName := child.TagName()

		if isUpperCase(tagName[0]) {

			rf := Invoke(React{}, tagName, child.Attrs(), Render(child.Html()))
			ret += rf[0].Interface().(string)
		} else {
			if tagName == `empty` {

				ret += child.Content()
			} else {
				ret += child.MakeHtml(Render(child.Html()))
			}
		}

	})

	return ret
}

func main() {

	fmt.Print(Render(`<Header title="test">test <Link to="https://www.google.com/">Google</Link></Header>
	<Footer title="test">test</Footer>`))
}
