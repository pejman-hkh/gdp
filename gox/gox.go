package gox

import (
	"fmt"
	"reflect"
	"regexp"

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
		r := regexp.MustCompile(`{{([^{}]*)}}`)
		matches := r.FindAllStringSubmatch(*value, -1)

		if len(matches) > 0 && len(matches[0]) > 1 {
			attrs += fmt.Sprintf("%s`%s` :%s", pre, key, matches[0][1])
		} else {
			attrs += fmt.Sprintf("%s`%s` :`%s`", pre, key, *value)
		}
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
				content := t.Content()

				r := regexp.MustCompile(`(.*?){([^{}]*)}(.*?)`)
				matches := r.FindAllStringSubmatch(content, -1)
				if len(matches) > 0 {
					ra := ""
					for _, v := range matches {
						ra += "`" + v[1] + "`"
						if v[2] != "" {
							ra += "," + v[2]
						}
						if v[3] != "" {
							ra += "," + v[3]
						}
					}
					if ra != "" {
						ret += pre + ra
					}
				} else {
					ret += pre + "`" + content + "`"
				}
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
