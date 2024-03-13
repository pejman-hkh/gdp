package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
	mainTitle := telm.Text()
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

func routes(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	route := strings.Split(path, "/")

	if route[1] == "api" {
		url := "https://www.imdb.com/title/" + route[3] + "/"

		w.Header().Set("Content-Type", "application/json")
		os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1088")
		os.Setenv("HTTP_PROXY", "socks5://127.0.0.1:1088")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// handle err
		}
		req.Header.Set("Authority", "www.imdb.com")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")
		req.Header.Set("Cache-Control", "max-age=0")
		req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		if res != nil {
			defer res.Body.Close()
			b, _ := io.ReadAll(res.Body)

			api := imdbApi(string(b))
			marshal, _ := json.Marshal(api)
			fmt.Fprint(w, string(marshal[:]))
		}
	}
}

func main() {
	//Just open the URL: http://localhost:8090/api/title/tt0137523 in your browser to see the API. Also, you can cache this result in the database
	//Unfortunately, I am in Iran, and I set a proxy in the client request here. You can remove it for your testing.
	http.HandleFunc("/", routes)
	http.ListenAndServe(":8090", nil)
}
