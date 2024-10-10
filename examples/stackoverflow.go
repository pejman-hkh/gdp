package main

import (
	"fmt"
	"os"

	"gdp/gdp"
)

func main() {
	html, _ := os.ReadFile("./data/stackoverflow.html")

	document := gdp.Default(string(html))

	document.Find("#tags-browser .grid--item").Each(func(index int, tag *gdp.Tag) {
		link := tag.Find("a").Eq(0)
		fmt.Printf("%+v\n", link.Attr("href"))
	})
}
