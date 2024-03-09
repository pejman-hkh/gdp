package main

import (
	"fmt"
	"os"

	"github.com/pejman-hkh/gdp/gdp"
)

func main() {

	fileContent, _ := os.ReadFile("fightclub.html")
	var document gdp.Tag = gdp.Default(string(fileContent))
	fmt.Printf("%+v", document)
}
