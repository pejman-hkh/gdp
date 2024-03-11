package gdp

import (
	"strings"
)

type Attr struct {
	name  string
	value string
}

func makeAttr(attrs []*Attr) string {
	ret := ""
	pre := ""

	for _, attr := range attrs {
		value := attr.value
		name := attr.name

		if name == "class" && value == "" {
			continue
		}

		ret += pre + name + `="` + value + `"`
		pre = " "

	}
	if ret != "" {
		ret = " " + ret
	}
	return ret
}

func getAttr(attrs []*Attr, index string) *Attr {
	var ret *Attr
	for _, attr := range attrs {

		if attr.name == index {
			ret = attr
			break
		}
	}

	return ret
}

func (a *Attr) inClass(v string) bool {
	if a == nil {
		return false
	}

	split := strings.Split(a.value, " ")
	for _, s := range split {
		if s == v {
			return true
		}
	}
	return false
}
