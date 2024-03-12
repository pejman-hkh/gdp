# GDP
GoLang Dom Parser

# API

### type Attr
### type NodeList
- func (n *NodeList) Each(callback func(int, *Tag))
- func (n *NodeList) Eq(index int) *Tag
### type Parser
- func (p *Parser) Parse(parent *Tag) []*Tag
### type QueryAttr
### type Tag
- func Default(html string) Tag
- func (tag *Tag) Attr(key string) string
- func (tag *Tag) Children() *NodeList
- func (tag *Tag) Find(mainQuery string) *NodeList
- func (tag *Tag) GetElementById(id string) *Tag
- func (tag *Tag) Html() string
- func (tag *Tag) Next() *Tag
- func (tag *Tag) Parent() *Tag
- func (tag *Tag) Prev() *Tag
- func (t *Tag) Print()
- func (tag *Tag) Remove()
- func (tag *Tag) SetAttr(key string, value string)
- func (tag *Tag) SetHtml(html string)

# Usage
```go
package main

import (
	"fmt"
	"os"

	"github.com/pejman-hkh/gdp/gdp"
)

func main() {

	fileContent, _ := os.ReadFile("fightclub.html")
	document := gdp.Default(string(fileContent))
	found := document.Find(".ipc-image")
	fmt.Printf("%+v\n", found.Eq(0).Attr("src"))

	document = gdp.Default(`<div class="parent"><div class="prev">test</div><div class="middle" id="middle">test1</div><span class="next"></span></div>`)

	fmt.Printf("%+v\n", document)

	document.Find(".prev,.middle,.next").Each(func(index int, tag *gdp.Tag) {
		fmt.Printf("%+v\n", tag)
	})

	middle := document.GetElementById("middle")
	fmt.Println(middle.Parent().Attr("class"))
	fmt.Println(middle.Prev().Attr("class"))
	fmt.Println(middle.Next().Attr("class"))
	fmt.Println(document.Html())

}
```