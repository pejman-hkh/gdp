package gdp

import "strings"

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
			ct := ""
			if len(child.children) > 0 {
				ct = child.concatHtmls()
			}
			html += child.makeHtml(ct)
		}

	}
	return html
}

func (tag *Tag) Html() string {
	return tag.concatHtmls()
}

func (tag *Tag) Attr(key string) string {
	attr := getAttr(tag.attrs, key)
	return attr.value
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
				var g string
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

func (tag *Tag) Find(query string) *NodeList {
	tags := tag.children
	if query == "" {
		return &NodeList{tags}
	}

	q := Query{query, 0, len(query)}
	ret := tag.children
	for {
		var splited string
		c := q.parseQuery(&splited)
		if !c {
			break
		}

		if splited == "," {
			q.parseQuery(&splited)

			qa := QueryAttr{splited, 0, len(splited)}
			attrs := qa.parseAttr()
			found := tag.findAttr(attrs, tags)
			ret = append(ret, found...)
		} else if strings.TrimSpace(splited) != "" {
			qa := QueryAttr{splited, 0, len(splited)}
			attrs := qa.parseAttr()
			ret = tag.findAttr(attrs, ret)

		}
	}

	return &NodeList{ret}
}
