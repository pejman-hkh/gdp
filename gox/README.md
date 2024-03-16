# GOX
Just like jsx you can write gox

```gox
package main

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/pejman-hkh/gdp/gox"
)

type React struct {}

func (r React) Link(props map[string]string, childrens string) string {
	return <a href={{props["to"]}}>{childrens}</a>
}

func (r React) SideNav(props map[string]string, childrens string) string {
	return <nav>{childrens}</nav>
}

func (r React) Content(props map[string]string, childrens string) string {
	return <main>{childrens}</main>
}

func (r React) Header(props map[string]string, childrens string) string {
	return <f>
	<header>
	<nav>
		<ul>
			<li>
			<a href="/">Home</a>
			</li>
		</ul>
	</nav>
	{childrens}
	</header>
	<h1>{props["title"]}</h1>
	</f>
}

func (r React) Footer(props map[string]string, childrens string) string {
	return <footer>
	<Link to="https://www.github.com/pejman-hkh/gdp">https://www.github.com/pejman-hkh/gdp</Link> 
	{childrens}
	</footer>
}

var react gox.Gox = gox.Gox{React{}}
func (r React) Layout(props map[string]string, childrens string) string {
	return <f>
	<html>
	<head>
	</head>
	<Header title="test">test 
		<Link to="/about">About</Link>
		<Link to="/contact">About</Link>
	</Header>
	<SideNav>
		<li><a href="/contact">Contact</a></li>
	</SideNav>
	<Content>
		{childrens}
	</Content>
	<Footer title="test">test</Footer>
	</html>
	</f>
}

func routes(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	route := strings.Split(path, "/")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if route[1] == "about" {
		fmt.Fprint(w, <Layout>
		About
		</Layout>)
	} else if route[1] == "contact" {
		fmt.Fprint(w, <Layout>
		Contact
		</Layout>)
	} else {
		fmt.Fprint(w, <Layout>
		Home Page
		</Layout>)
	}
}

func main() {

	http.HandleFunc("/", routes)
	http.ListenAndServe(":8090", nil)

}
```

This will convert to :
```go
package main

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/pejman-hkh/gdp/gox"
)

type React struct {}

func (r React) Link(props map[string]string, childrens string) string {
	return react.Run("a", map[string]string{`href` :props["to"]}, []string{``,childrens})
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
func (r React) Layout(props map[string]string, childrens string) string {
	return react.Run("f", map[string]string{}, []string{`
	`, react.Run("html", map[string]string{}, []string{`
	`, react.Run("head", map[string]string{}, []string{`
	`}), `
	`, react.Run("Header", map[string]string{`title` :`test`}, []string{`test 
		`, react.Run("Link", map[string]string{`to` :`/about`}, []string{`About`}), `
		`, react.Run("Link", map[string]string{`to` :`/contact`}, []string{`About`}), `
	`}), `
	`, react.Run("SideNav", map[string]string{}, []string{`
		`, react.Run("li", map[string]string{}, []string{react.Run("a", map[string]string{`href` :`/contact`}, []string{`Contact`})}), `
	`}), `
	`, react.Run("Content", map[string]string{}, []string{`		`,childrens}), `
	`, react.Run("Footer", map[string]string{`title` :`test`}, []string{`test`}), `
	`}), `
	`})
}

func routes(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	route := strings.Split(path, "/")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if route[1] == "about" {
		fmt.Fprint(w, react.Run("Layout", map[string]string{}, []string{`
		About
		`}))
	} else if route[1] == "contact" {
		fmt.Fprint(w, react.Run("Layout", map[string]string{}, []string{`
		Contact
		`}))
	} else {
		fmt.Fprint(w, react.Run("Layout", map[string]string{}, []string{`
		Home Page
		`}))
	}
}

func main() {

	http.HandleFunc("/", routes)
	http.ListenAndServe(":8090", nil)

}

```