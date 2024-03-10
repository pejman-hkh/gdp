package main

import (
	"fmt"
	"os"

	"github.com/pejman-hkh/gdp/gdp"
)

func main() {
	html, _ := os.ReadFile("stackoverflow.html")

	document := gdp.Default(string(html))
	//fmt.Print(document.Html())
	//document.Print()

	document.Find("#tags-browser .grid--item").Each(func(index int, tag *gdp.Tag) {
		link := tag.Find("a").Eq(0)
		fmt.Printf("%+v\n", link.Attr("href"))
	})
}
