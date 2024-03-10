package gdp

import "strings"

type Attr struct {
	name  string
	value string
}

func getAttr(attrs []Attr, index string) Attr {
	var ret Attr
	for _, attr := range attrs {
		if attr.name == index {
			ret = attr
			break
		}
	}
	return ret
}

func (a *Attr) inClass(v string) bool {
	split := strings.Split(a.value, " ")
	for _, s := range split {
		if s == v {
			return true
		}
	}
	return false
}
