package main

import (
	"fmt"
	"os"

	"github.com/pejman-hkh/gdp/gdp"
)

func main() {

	fileContent, _ := os.ReadFile("fightclub.html")
	var document gdp.Tag = gdp.Default(string(fileContent))

	found := document.Find(".ipc-image")

	fmt.Printf("%+v", found[0].Attr("src"))

}
