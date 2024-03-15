package gdp

import (
	"fmt"
	"strconv"
)

type Tag struct {
	tag      string
	content  string
	attrs    Attr
	isEnd    bool
	eq       int
	next     *Tag
	prev     *Tag
	parent   *Tag
	children []*Tag
}

func (t *Tag) Print() {
	for _, tag := range t.children {
		fmt.Printf("{tag:%s, content:%s,", tag.tag, tag.content)
		fmt.Printf("attr:%s,", tag.attrs.makeAttr())
		if len(tag.children) > 0 {
			fmt.Printf("children:")
			tag.Print()
			fmt.Printf(",")
		}
		fmt.Printf("}\n")
	}
}

func (tag *Tag) MakeHtml(content string) string {
	if tag.tag == "script" {
		return "<script" + tag.attrs.makeAttr() + ">" + tag.content + "</script>"
	} else if tag.tag == "comment" {
		return ""
	}

	if isEndTag(tag) {
		return "<" + tag.tag + "" + tag.attrs.makeAttr() + " />"
	}

	return "<" + tag.tag + tag.attrs.makeAttr() + ">" + (content) + "</" + tag.tag + ">"
}

func (tag *Tag) concatHtmls() string {
	children := tag.children
	html := ""
	for _, child := range children {
		if child == nil {
			continue
		}
		if child.tag == "empty" || child.tag == "cdata" {
			html += child.content
		} else {
			content := ""
			if len(child.children) > 0 {
				content = child.concatHtmls()
			}
			html += child.MakeHtml(content)
		}

	}
	return html
}

func (tag *Tag) concatTexts() string {
	html := ""
	children := tag.children
	for _, child := range children {
		if child.tag == "empty" {
			html += child.content
		}

		if len(child.children) > 0 {
			html += child.concatTexts()
		}
	}
	return html
}

func (tag *Tag) Text() string {
	return tag.concatTexts()
}

func (tag *Tag) Html() string {
	return tag.concatHtmls()
}

func (tag *Tag) OuterHtml() string {
	content := tag.concatHtmls()
	return tag.MakeHtml(content)
}

func (tag *Tag) Attr(key string) string {
	return tag.attrs.valueOf(key)
}

func (tag *Tag) Attrs() map[string]*string {
	return tag.attrs.attrs
}

func (tag *Tag) SetAttr(key string, value string) {
	tag.attrs.setValue(key, value)
}

func (tag *Tag) RemoveClass(class string) {
	tag.attrs.RemoveClass(class)
}

func (tag *Tag) AddClass(class string) {
	tag.attrs.AddClass(class)
}

func (tag *Tag) HasClass(class string) bool {
	return tag.attrs.HasClass(class)
}

func (tag *Tag) GetElementById(id string) *Tag {
	return tag.Find("#" + id).Eq(0)
}

func (tag *Tag) Parent() *Tag {
	return tag.parent
}

func (tag *Tag) Prev() *Tag {
	return tag.prev
}

func (tag *Tag) Next() *Tag {
	return tag.next
}

func (tag *Tag) TagName() string {
	return tag.tag
}

func (tag *Tag) Content() string {
	return tag.content
}

func (tag *Tag) Children() *NodeList {
	return &NodeList{tag.children}
}

func (tag *Tag) Remove() {
	tag.parent.children[tag.eq] = nil
}

func (tag *Tag) SetHtml(html string) {
	document := Default(html)
	tag.children = document.children
	for _, child := range document.children {
		child.parent = tag
	}
}

func (mtag *Tag) findAttr(attrs map[string]string, tags []*Tag, sp string, level int) []*Tag {

	var ret []*Tag
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		f := true
		for attr, value := range attrs {
			if attr == "eq" {
				continue
			}

			if attr == "class" {

				if !tag.attrs.HasClass(value) {
					f = false
				}
			} else {
				g := ""

				if attr == "tag" {
					g = tag.tag
				} else {

					a := tag.attrs.valueOf(attr)

					if a != "" {
						g = a
					}

				}

				if g != value {
					f = false
				}
			}

		}

		if f {
			ret = append(ret, tag)
		}
		findChild := false
		if sp == ">" {
			if level < 1 {
				findChild = true
			}
		} else if len(tag.children) > 0 {
			findChild = true
		}

		if findChild {
			found := tag.findAttr(attrs, tag.children, sp, level+1)
			ret = append(ret, found...)
		}
	}

	return ret
}

func (tag *Tag) Find(mainQuery string) *NodeList {
	tags := tag.children
	if mainQuery == "" {
		return &NodeList{tags}
	}

	var ret []*Tag
	splits := splitQueries(mainQuery)

	for _, split := range splits {

		found := tag.children
		query := splitQuery(split)
		sp := ""
		for _, q := range query {
			if q == " " || q == ">" || q == "|" || q == "~" || q == "+" {
				sp = q
				continue
			}

			qa := queryAttr{q, 0, len(q)}
			attrs := qa.parseAttr()
			found = tag.findAttr(attrs, found, sp, 0)
			if _, ok := attrs["first"]; ok {
				if len(found) > 0 {
					tmp := []*Tag{}
					tmp = append(tmp, found[0])
					found = tmp

				}
			} else if _, ok := attrs["last"]; ok {
				if len(found) > 0 {
					tmp := []*Tag{}
					tmp = append(tmp, found[len(found)-1])
					found = tmp

				}
			} else if val, ok := attrs["eq"]; ok {
				if len(found) > 0 {
					i, _ := strconv.Atoi(val)
					if i < len(found) {
						tmp := []*Tag{}
						tmp = append(tmp, found[i])
						found = tmp
					} else {
						found = []*Tag{}
					}

				}
			}

		}

		ret = append(ret, found...)
	}

	return &NodeList{ret}
}
