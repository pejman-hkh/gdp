# gdp
GoLang Dom Parser

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
	fmt.Printf("%+v", found.Eq(0).Attr("src"))

	document = gdp.Default(`<div class="test">test</div><div class="test1">test1</div>`)
	document.Find(".test,.test1").Each(func(index int, tag *gdp.Tag) {
		fmt.Print(tag)
	})
	fmt.Printf("%+v", found)
}


```