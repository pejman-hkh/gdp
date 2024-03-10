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
	fmt.Printf("%+v", found[0].Attr("src"))

	document = gdp.Default(`<div class="test">test</div><div class="test1">test1</div>`)
	found = document.Find(".test,.test1")
	fmt.Printf("%+v", found)
}
