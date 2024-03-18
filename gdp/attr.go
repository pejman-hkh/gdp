package gdp

import (
	"strings"
)

type Attr struct {
	Attrs map[string]*string
}

func (a *Attr) makeAttr() string {
	ret := ""
	pre := ""
	if a.Attrs == nil {
		return ""
	}

	for name, value := range a.Attrs {

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
	if a.Attrs != nil {
		(a.Attrs)[key] = &value
	}
}

func (a *Attr) valueOf(key string) string {
	if a.Attrs != nil {
		v := (a.Attrs)[key]
		if v != nil {
			return *v
		}
	}
	return ""
}

func (a *Attr) RemoveClass(class string) {
	if a.Attrs == nil {
		return
	}

	split := strings.Split(a.valueOf("class"), " ")
	cls := ""
	pre := ""
	for _, s := range split {
		if s != class {
			cls += pre + s
			pre = " "
		}
	}
	a.setValue("class", cls)
}

func (a *Attr) AddClass(class string) {
	if a.Attrs == nil {
		return
	}

	if !a.HasClass(class) {
		classes := a.valueOf("class")
		if classes != "" {
			classes += " "
		}
		a.setValue("class", classes+class)
	}
}

func (a *Attr) HasClass(v string) bool {

	if a.Attrs == nil {
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
