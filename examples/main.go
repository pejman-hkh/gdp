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

	document = gdp.Default(`<div class="parent"><div class="prev">test</div><div class="middle" id="middle">test1</div><span class="next"></span></div>`)

	fmt.Printf("%+v", document)

	document.Find(".prev,.middle,.next").Each(func(index int, tag *gdp.Tag) {
		fmt.Println(tag)
	})

	middle := document.GetElementById("middle")
	fmt.Println(middle.Parent().Attr("class"))
	fmt.Println(middle.Prev().Attr("class"))
	fmt.Println(middle.Next().Attr("class"))

}
