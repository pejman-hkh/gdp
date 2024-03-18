package gdp

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Tag struct {
	Tag       string
	Contents  string
	TagAttrs  Attr
	isEnd     bool
	Eq        int
	next      *Tag
	prev      *Tag
	parent    *Tag
	Childrens []*Tag
}

func (t *Tag) Print() {
	b1, _ := json.Marshal(t.Childrens)
	fmt.Println(string(b1))
}

func (tag *Tag) MakeHtml(content string) string {
	if tag.Tag == "script" {
		return "<script" + tag.TagAttrs.makeAttr() + ">" + tag.Contents + "</script>"
	} else if tag.Tag == "comment" {
		return ""
	}

	if isEndTag(tag) {
		return "<" + tag.Tag + "" + tag.TagAttrs.makeAttr() + " />"
	}

	return "<" + tag.Tag + tag.TagAttrs.makeAttr() + ">" + (content) + "</" + tag.Tag + ">"
}

func (tag *Tag) concatHtmls() string {
	children := tag.Childrens
	html := ""
	for _, child := range children {
		if child == nil {
			continue
		}
		if child.Tag == "empty" || child.Tag == "cdata" {
			html += child.Contents
		} else {
			content := ""
			if len(child.Childrens) > 0 {
				content = child.concatHtmls()
			}
			html += child.MakeHtml(content)
		}

	}
	return html
}

func (tag *Tag) concatTexts() string {
	html := ""
	children := tag.Childrens
	for _, child := range children {
		if child.Tag == "empty" {
			html += child.Contents
		}

		if len(child.Childrens) > 0 {
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
	return tag.TagAttrs.valueOf(key)
}

func (tag *Tag) Attrs() map[string]*string {
	return tag.TagAttrs.Attrs
}

func (tag *Tag) SetAttr(key string, value string) {
	tag.TagAttrs.setValue(key, value)
}

func (tag *Tag) RemoveClass(class string) {
	tag.TagAttrs.RemoveClass(class)
}

func (tag *Tag) AddClass(class string) {
	tag.TagAttrs.AddClass(class)
}

func (tag *Tag) HasClass(class string) bool {
	return tag.TagAttrs.HasClass(class)
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
	return tag.Tag
}

func (tag *Tag) Content() string {
	return tag.Contents
}

func (tag *Tag) Children() *NodeList {
	return &NodeList{tag.Childrens}
}

func (tag *Tag) Remove() {
	tag.parent.Childrens[tag.Eq] = nil
}

func (tag *Tag) SetHtml(html string) {
	document := Default(html)
	tag.Childrens = document.Childrens
	for _, child := range document.Childrens {
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

				if !tag.TagAttrs.HasClass(value) {
					f = false
				}
			} else {
				g := ""

				if attr == "tag" {
					g = tag.Tag
				} else {

					a := tag.TagAttrs.valueOf(attr)

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
		} else if len(tag.Childrens) > 0 {
			findChild = true
		}

		if findChild {
			found := tag.findAttr(attrs, tag.Childrens, sp, level+1)
			ret = append(ret, found...)
		}
	}

	return ret
}

func (tag *Tag) Find(mainQuery string) *NodeList {
	tags := tag.Childrens
	if mainQuery == "" {
		return &NodeList{tags}
	}

	var ret []*Tag
	splits := splitQueries(mainQuery)

	for _, split := range splits {

		found := tag.Childrens
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
