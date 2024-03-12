package gdp

import (
	"strings"
)

type Attr struct {
	attrs *map[string]*string
}

func (a *Attr) makeAttr() string {
	ret := ""
	pre := ""
	if a.attrs == nil {
		return ""
	}

	for name, value := range *a.attrs {

		if name == "class" && *value == "" {
			continue
		}

		ret += pre + name + `="` + *value + `"`
		pre = " "

	}
	if ret != "" {
		ret = " " + ret
	}
	return ret
}

func (a *Attr) setValue(key string, value string) {
	if a.attrs != nil {
		(*a.attrs)[key] = &value
	}
}

func (a *Attr) valueOf(key string) string {
	if a.attrs != nil {
		v := (*a.attrs)[key]
		if v != nil {
			return *v
		}
	}
	return ""
}

func (a *Attr) inClass(v string) bool {

	if a.attrs == nil {
		return false
	}

	split := strings.Split(a.valueOf("class"), " ")
	for _, s := range split {
		if s == v {
			return true
		}
	}
	return false
}
