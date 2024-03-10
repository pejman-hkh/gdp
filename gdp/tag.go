package gdp

type Tag struct {
	tag      string
	content  string
	attrs    []Attr
	isEnd    bool
	eq       int
	next     *Tag
	prev     *Tag
	parent   *Tag
	children []Tag
}

func (tag *Tag) Attr(key string) string {
	attr := getAttrIndex(tag.attrs, key)
	return attr.value
}

func (tag *Tag) findAttr(attrs map[string]string, tags []Tag) []Tag {
	var ret []Tag
	for _, tag := range tags {
		f := true
		for attr, value := range attrs {
			if attr == "class" {
				classAttr := getAttrIndex(tag.attrs, "class")

				if !classAttr.inClass(value) {
					f = false
				}
			} else {
				var g string
				if attr == "tag" {
					g = tag.tag
				} else {
					a := getAttrIndex(tag.attrs, attr)
					g = a.value
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

func (tag *Tag) Find(query string) []Tag {
	tags := tag.children
	if query == "" {
		return tags
	}

	q := Query{query, 0, len(query)}
	var ret []Tag
	for {
		var splited string
		c := q.parseQuery(&splited)
		if !c {
			break
		}

		if splited == "," {
			qa := QueryAttr{splited, 0, len(splited)}
			attrs := qa.parseAttr()
			found := tag.findAttr(attrs, tags)
			ret = append(ret, found...)
		} else if splited != " " {
			qa := QueryAttr{splited, 0, len(splited)}
			attrs := qa.parseAttr()
			ret = tag.findAttr(attrs, tags)
		}
	}

	return ret
}
