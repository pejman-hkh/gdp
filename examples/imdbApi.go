package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pejman-hkh/gdp/gdp"
)

type Casts map[int]map[string]string

func inArray(key string, array []string) bool {
	for _, v := range array {
		if v == key {
			return true
		}
	}
	return false
}

func imdbApi(content string) map[string]interface{} {

	document := gdp.Default(content)
	epic := document.Find(".ipc-image").Eq(0)
	mainPic := []string{epic.Attr("src"), epic.Attr("srcSet")}

	selm := document.Find("div[data-testid=hero-rating-bar__aggregate-rating__score]").Eq(0)
	sselm := selm.Find("span").Eq(0)
	rate := sselm.Html()

	rated := selm.Next().Next().Html()

	telm := document.Find("h1[data-testid=hero__pageTitle]").Eq(0)
	mainTitle := telm.Html()
	info := []string{}
	telm.Next().Find("li").Each(func(i int, tag *gdp.Tag) {
		info = append(info, tag.Html())
	})

	plot := document.Find("p[data-testid=plot] span").Eq(0).Html()

	arr := []string{"Director", "Writers", "Stars", "Directors", "Writer", "Star"}

	casts := make(map[string]Casts)
	document.Find(".ipc-inline-list").Each(func(ii int, t *gdp.Tag) {

		title := t.Parent().Prev().Html()
		if inArray(title, arr) {
			casts[title] = Casts{}
			castArray := make(Casts)
			i := 0
			t.Find("a").Each(func(iii int, a *gdp.Tag) {
				cast := map[string]string{"name": a.Html(), "link": a.Attr("href")}
				castArray[i] = cast
				i++
			})
			casts[title] = castArray
		}
	})

	topCast := make(Casts)
	i := 0
	document.Find("div[data-testid=title-cast-item]").Each(func(ii int, cast *gdp.Tag) {
		pic := cast.Find("img").Eq(0)
		link := cast.Find("a").Eq(0)
		castMap := map[string]string{"name": link.Attr("aria-label"), "link": link.Attr("href"), "pic": pic.Attr("src"), "pics": pic.Attr("srcSet")}
		topCast[i] = castMap
		i++
	})

	eret := make(map[string]interface{})
	eret["title"] = mainTitle
	eret["info"] = info
	eret["rate"] = rate
	eret["rated"] = rated
	eret["pic"] = mainPic
	eret["topcast"] = topCast
	eret["casts"] = casts
	eret["plot"] = plot

	return eret
}

func main() {
	fileContent, _ := os.ReadFile("fightclub.html")
	api := imdbApi(string(fileContent))
	marshal, _ := json.Marshal(api)
	fmt.Print(string(marshal[:]))
}
