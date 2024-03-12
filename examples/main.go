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
	middle.Parent().SetAttr("class", "parent1")
	fmt.Println(middle.Parent().Attr("class"))
	fmt.Println(middle.Prev().Attr("class"))
	fmt.Println(middle.Next().Attr("class"))
	fmt.Println(document.Html())

	document = gdp.Default(`<div id="test">test</div><div id="test1">test1</div>`)
	document.Find("#test").Eq(0).Remove()
	fmt.Print(document.Html())
	fmt.Print(document.Find("#test").Eq(0))

	document = gdp.Default(`<div class="parent"><div class="prev">test</div><div class="middle" id="middle">test1</div><span class="next"></span></div>`)
	tag := document.Find(".parent").Eq(0)
	if tag.HasClass("parent") {
		fmt.Println("Parent has class parent")
	}
	tag.SetHtml("<span>changed html</span>")
	tag.Find("span").Eq(0).AddClass("test")
	fmt.Print(tag.Html())
	tag.Find("span").Eq(0).RemoveClass("test")
	fmt.Print(tag.Html())
}
