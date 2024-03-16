package main

import (
	"fmt"
	"os"

	"github.com/pejman-hkh/gdp/gdp"
	"github.com/pejman-hkh/gdp/gox"
)

func main() {
	goxFile, _ := os.ReadFile("test.gox")
	document := gdp.Default(string(goxFile))

	out := gox.ToGo(&document)
	fmt.Print(out)
}
