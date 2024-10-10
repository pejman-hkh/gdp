package gdp

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Tag struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Attrs    Attr   `json:"attrs"`
	Eq       int    `json:"eq"`
	Children []*Tag `json:"children"`
	next     *Tag
	prev     *Tag
	parent   *Tag
}

func (t *Tag) Print() {
	b1, _ := json.MarshalIndent(t.Children, " ", " ")
	fmt.Println(string(b1))
}

func (tag *Tag) MakeHtml(content string) string {
	if tag.Name == "script" {
		return "<script" + tag.Attrs.makeAttr() + ">" + tag.Content + "</script>"
	} else if tag.Name == "comment" {
		return ""
	}

	if inStringArray(tag.Name, noEndTags[:]) {
		return "<" + tag.Name + "" + tag.Attrs.makeAttr() + " />"
	}

	return "<" + tag.Name + tag.Attrs.makeAttr() + ">" + (content) + "</" + tag.Name + ">"
}

func (tag *Tag) concatHtmls() string {
	children := tag.Children
	html := ""
	if len(children) > 0 {

		for _, child := range children {
			if child == nil {
				continue
			}
			if child.Name == "empty" || child.Name == "cdata" {
				html += child.Content
			} else {
				content := ""
				if len(child.Children) > 0 {
					content = child.concatHtmls()
				}
				html += child.MakeHtml(content)
			}

		}
	}
	return html
}

func (tag *Tag) concatTexts() string {
	html := ""
	children := tag.Children
	for _, child := range children {
		if child.Name == "empty" {
			html += child.Content
		}

		if len(child.Children) > 0 {
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
	return tag.Attrs.valueOf(key)
}

func (tag *Tag) GetAttrs() map[string]*string {
	return tag.Attrs.Attrs
}

func (tag *Tag) SetAttr(key string, value string) {
	tag.Attrs.setValue(key, value)
}

func (tag *Tag) RemoveClass(class string) {
	tag.Attrs.RemoveClass(class)
}

func (tag *Tag) AddClass(class string) {
	tag.Attrs.AddClass(class)
}

func (tag *Tag) HasClass(class string) bool {
	return tag.Attrs.HasClass(class)
}

func (tag *Tag) GetElementById(id string) *Tag {
	return tag.Find("#" + id).Eq(0)
}

func (tag *Tag) Parent() *Tag {
	if tag.parent == nil {
		return &Tag{}
	}
	return tag.parent
}

func (tag *Tag) Prev() *Tag {
	if tag.prev == nil {
		return &Tag{}
	}

	return tag.prev
}

func (tag *Tag) Next() *Tag {
	if tag.next == nil {
		return &Tag{}
	}
	return tag.next
}

func (tag *Tag) TagName() string {
	return tag.Name
}

func (tag *Tag) GetContent() string {
	return tag.Content
}

func (tag *Tag) GetChildren() *NodeList {
	return &NodeList{tag.Children}
}

func (tag *Tag) Remove() {
	tag.parent.Children[tag.Eq] = nil
}

func (tag *Tag) SetHtml(html string) {
	document := Default(html)
	tag.Children = document.Children
	for _, child := range document.Children {
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

				if !tag.Attrs.HasClass(value) {
					f = false
				}
			} else {
				g := ""

				if attr == "tag" {
					g = tag.Name
				} else {

					a := tag.Attrs.valueOf(attr)

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
		} else if len(tag.Children) > 0 {
			findChild = true
		}

		if findChild {
			found := tag.findAttr(attrs, tag.Children, sp, level+1)
			ret = append(ret, found...)
		}
	}

	return ret
}

func (tag *Tag) Find(mainQuery string) *NodeList {
	tags := tag.Children
	if mainQuery == "" {
		return &NodeList{tags}
	}

	var ret []*Tag
	splits := splitQueries(mainQuery)

	for _, split := range splits {

		found := tag.Children
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
