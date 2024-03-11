package gdp

import (
	"fmt"
)

type Tag struct {
	tag      string
	content  string
	attrs    []*Attr
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
		fmt.Printf("attr:%s,", makeAttr(tag.attrs))
		if len(tag.children) > 0 {
			fmt.Printf("children:")
			tag.Print()
			fmt.Printf(",")
		}
		fmt.Printf("}\n")
	}
}

func (tag *Tag) makeHtml(content string) string {
	if tag.tag == "script" {
		return "<script" + makeAttr(tag.attrs) + ">" + tag.content + "</script>"
	} else if tag.tag == "comment" {
		return ""
	}

	if isEndTag(tag) {
		return "<" + tag.tag + "" + makeAttr(tag.attrs) + " />"
	}

	return "<" + tag.tag + makeAttr(tag.attrs) + ">" + (content) + "</" + tag.tag + ">"
}

func (tag *Tag) concatHtmls() string {
	children := tag.children
	html := ""
	for _, child := range children {
		if child.tag == "empty" || child.tag == "cdata" {
			html += child.content
		} else {
			content := ""
			if len(child.children) > 0 {
				content = child.concatHtmls()
			}
			html += child.makeHtml(content)
		}

	}
	return html
}

func (tag *Tag) Html() string {
	return tag.concatHtmls()
}

func (tag *Tag) Attr(key string) string {
	attr := getAttr(tag.attrs, key)
	if attr != nil {
		return attr.value
	}
	return ""
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

func (tag *Tag) Children() *NodeList {
	return &NodeList{tag.children}
}

func (mtag *Tag) findAttr(attrs map[string]string, tags []*Tag) []*Tag {

	var ret []*Tag
	for _, tag := range tags {
		f := true
		for attr, value := range attrs {
			if attr == "class" {
				classAttr := getAttr(tag.attrs, "class")

				if !classAttr.inClass(value) {
					f = false
				}
			} else {
				g := ""

				if attr == "tag" {
					g = tag.tag
				} else {
					a := getAttr(tag.attrs, attr)
					if a != nil {
						g = a.value
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

		if len(tag.children) > 0 {
			found := tag.findAttr(attrs, tag.children)
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

		for _, q := range query {
			if q == " " {
				continue
			}

			qa := QueryAttr{q, 0, len(q)}
			attrs := qa.parseAttr()
			found = tag.findAttr(attrs, found)
		}

		ret = append(ret, found...)
	}

	return &NodeList{ret}
}
