package gox

import (
	"fmt"
	"reflect"

	"github.com/pejman-hkh/gdp/gdp"
)

func isUpperCase(c byte) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	}
	return false
}

func Invoke(obj any, name string, args ...any) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	return reflect.ValueOf(obj).MethodByName(name).Call(inputs)
}

type Gox struct {
	Obj any
}

func (g *Gox) Run(name string, props map[string]string, childs []string) string {

	html := ""
	for _, v := range childs {
		html += v
	}

	attrs := ""
	for k, v := range props {
		attrs += fmt.Sprintf(" %s=\"%s\"", k, v)
	}
	if name == "f" {
		return html
	} else if !isUpperCase(name[0]) {
		return `<` + name + `` + attrs + `>` + html + `</` + name + `>`
	} else {
		rf := Invoke(g.Obj, name, props, html)
		return rf[0].Interface().(string)
	}
}

func convertToGoxFunc(tag *gdp.Tag, child string) string {
	attrs := ""
	pre := ""
	for key, value := range tag.Attrs() {
		attrs += fmt.Sprintf("%s`%s` :`%s`", pre, key, *value)
		pre = ","
	}

	return fmt.Sprintf(`react.Run("%s", map[string]string{`+attrs+`}, %s)`, tag.TagName(), child)
}

func ToGo(tag *gdp.Tag) string {
	ret := ""
	pre := ""
	tag.Children().Each(func(i int, t *gdp.Tag) {
		if t.TagName() == "empty" {
			if t.Parent().TagName() == "document" {
				ret += t.Content()
			} else {
				ret += pre + "`" + t.Content() + "`"
			}
		} else {
			childs := `[]string{`

			if t.Children().Length() > 0 {
				childs += ToGo(t)
			}
			childs += `}`
			ret += pre + convertToGoxFunc(t, childs)
		}
		if t.Parent().TagName() != "document" {
			pre = ", "
		}
	})
	return ret
}
