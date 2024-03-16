package main

import (
	"fmt"
	"github.com/pejman-hkh/gdp/gox"
)

type React struct {}

func (r React) Link(props map[string]string, childrens string) string {
	return react.Run("a", map[string]string{`href` :`{{props["to"]}}`}, []string{``,childrens})
}

func (r React) SideNav(props map[string]string, childrens string) string {
	return react.Run("nav", map[string]string{}, []string{``,childrens})
}

func (r React) Content(props map[string]string, childrens string) string {
	return react.Run("main", map[string]string{}, []string{``,childrens})
}

func (r React) Header(props map[string]string, childrens string) string {
	return react.Run("f", map[string]string{}, []string{`
	`, react.Run("header", map[string]string{}, []string{`
	`, react.Run("nav", map[string]string{}, []string{`
		`, react.Run("ul", map[string]string{}, []string{`
			`, react.Run("li", map[string]string{}, []string{`
			`, react.Run("a", map[string]string{`href` :`/`}, []string{`Home`}), `
			`}), `
		`}), `
	`}), `	`,childrens}), `
	`, react.Run("h1", map[string]string{}, []string{``,props["title"]}), `
	`})
}

func (r React) Footer(props map[string]string, childrens string) string {
	return react.Run("footer", map[string]string{}, []string{`
	`, react.Run("Link", map[string]string{`to` :`https://www.github.com/pejman-hkh/gdp`}, []string{`https://www.github.com/pejman-hkh/gdp`}), `	`,childrens})
}

var react gox.Gox = gox.Gox{React{}}
func layout() string {
	return react.Run("f", map[string]string{}, []string{`
	`, react.Run("html", map[string]string{}, []string{`
	`, react.Run("head", map[string]string{}, []string{`
	`}), `
	`, react.Run("Header", map[string]string{`title` :`test`}, []string{`test 
		`, react.Run("Link", map[string]string{`to` :`https://www.google.com/`}, []string{`Google`}), `
	`}), `
	`, react.Run("SideNav", map[string]string{}, []string{`
		`, react.Run("li", map[string]string{}, []string{react.Run("a", map[string]string{`href` :`https://github.com`}, []string{`Test`})}), `
	`}), `
	`, react.Run("Content", map[string]string{}, []string{`
		main content
	`}), `
	`, react.Run("Footer", map[string]string{`title` :`test`}, []string{`test`}), `
	`}), `
	`})
}

func main() {
	fmt.Print(layout())
}
