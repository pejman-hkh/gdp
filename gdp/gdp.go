package gdp

import (
	"bytes"
)

var hasNoEndTags = [17]string{"comment", "php", "empty", "!DOCTYPE", "area", "base", "col", "embed", "param", "source", "track", "meta", "link", "br", "input", "hr", "img"}

func Default(html string) Tag {
	var p Parser
	p.html = html
	p.len = len(p.html)

	var document Tag
	document.tag = "document"
	document.children = p.Parse(&document)

	p.current = &document

	return document
}

type Parser struct {
	html    string
	len     int
	i       int
	current *Tag
	isXml   bool
}

func (p *Parser) getUntil(until string, first byte) string {
	var buffer bytes.Buffer
	if first != 0 {
		buffer.WriteByte(first)
	}

	for {
		c := p.html[p.i]
		p.i++
		if c == until[0] && p.isEqual(until[1:]) {
			break
		}

		buffer.WriteByte(c)
	}
	return buffer.String()
}

func (p *Parser) parseAttr() []Attr {
	var attrs []Attr

	for {
		isThereValue := false
		var buffer bytes.Buffer
		for {
			c1 := p.html[p.i]
			p.i++
			if c1 == ' ' || c1 == '>' || c1 == '=' {
				if c1 == '=' {
					isThereValue = true
				}

				if c1 == '>' {
					p.i--
				}

				break
			}

			buffer.WriteByte(c1)
		}

		name := buffer.String()
		var buffer1 bytes.Buffer
		if isThereValue {
			g := p.html[p.i]
			var t byte = 0
			if g == '"' || g == '\'' {
				t = g
				p.i++
			}

			for {
				c1 := p.html[p.i]
				p.i++
				if c1 == t {
					break
				}

				if t == 0 && c1 == ' ' {
					break
				}

				if t == 0 && c1 == '>' {
					p.i--
					break
				}

				buffer1.WriteByte(c1)
			}
		}
		value := buffer1.String()

		if len(name) > 0 && name[0] != '/' && name[0] != ' ' {
			var attr Attr
			attr.name = name
			attr.value = value

			attrs = append(attrs, attr)
		}

		c1 := p.html[p.i]
		p.i++
		if c1 == '>' {
			break
		}
	}

	return attrs
}

func (p *Parser) parseTag(tag *Tag) bool {
	if p.isEqual("![CDATA[") {
		p.i += 8
		tag.tag = "cdata"
		return true
	}

	if p.html[p.i+1] == '/' {
		p.i++
	}

	var buffer bytes.Buffer
	var attrs []Attr
	for {
		if p.i == p.len {
			break
		}

		c1 := p.html[p.i]
		p.i++

		if c1 == '>' {
			break
		}

		if c1 == ' ' {
			attrs = p.parseAttr()
			break
		}
		buffer.WriteByte(c1)

	}

	var name string = buffer.String()
	tag.tag = name
	tag.attrs = attrs
	tag.isEnd = false

	if name[0] == '/' {
		tag.isEnd = true
		tag.tag = name[1:]
	}

	if name[len(name)-1] == '/' {
		name = name[0 : len(name)-1]
		tag.tag = name
	}
	return true
}

func (p *Parser) parseContent(first byte, tag *Tag) bool {
	p.i--
	content := p.getUntil("<", first)
	p.i--

	tag.tag = "empty"
	tag.content = content

	return true
}

func (p *Parser) parseComment(tag *Tag) bool {
	p.i += 3

	content := p.getUntil("-->", 0)

	p.i += 2

	tag.tag = "comment"
	tag.content = content

	return true
}

func (p *Parser) parseScript() string {
	content := p.getUntil("</script", 0)
	p.i += 8
	return content
}

func (p *Parser) parseCdata() string {
	content := p.getUntil("]]>", 0)
	p.i += 2
	return content
}

func (p *Parser) isEqual(text string) bool {
	textLen := len(text)
	if p.i+textLen >= p.len {
		return false
	}

	html := p.html[p.i : p.i+textLen]

	if html == text {
		return true
	}

	return false
}

func (p *Parser) isEndTag(tag *Tag) bool {
	for i := 0; i < 17; i++ {
		if tag.tag == hasNoEndTags[i] {
			return true
		}
	}
	return false
}

func (p *Parser) next1(tag *Tag) bool {
	if p.i >= p.len {
		return false
	}

	c := p.html[p.i]
	p.i++

	if c == '<' {
		if p.isEqual("!--") {
			return p.parseComment(tag)
		}

		if p.html[p.i] == ' ' {
			p.i++
			return p.parseContent('<', tag)
		}

		return p.parseTag(tag)
	} else {
		return p.parseContent(0, tag)
	}

}

func (p *Parser) next(tag *Tag) bool {
	ret := p.next1(tag)
	if ret {
		p.current = tag
	}

	return ret
}

func (p *Parser) getTag(tag *Tag) bool {
	ret := p.next(tag)
	if !ret {
		return false
	}

	if tag.tag == "cdata" {
		tag.content = p.parseCdata()
		return true
	}
	name := tag.tag
	if len(name) >= 4 && name[0:4] == "?xml" {
		p.isXml = true
		return true
	}

	if p.isXml {
		hasNoEndTags[11] = ""
	}

	if p.isEndTag(tag) || tag.isEnd {
		return true
	}

	if tag.tag == "script" {
		tag.content = p.parseScript()
	} else {
		tag.children = p.Parse(tag)
	}

	if tag.tag == p.current.tag {
		return true
	}

	var etag Tag

	for p.next(&etag) {
		if tag.tag == etag.tag {
			break
		}
	}

	return true
}

func (p *Parser) Parse(parent *Tag) []Tag {
	var tags []Tag
	var eq int = 0
	var stag *Tag = &Tag{}
	for {
		var tag Tag

		ret := p.getTag(&tag)
		if !ret {
			break
		}

		if tag.isEnd && parent.tag == tag.tag {
			break
		}

		if !tag.isEnd {

			tag.eq = eq
			eq++
			tag.prev = stag
			tag.parent = parent
			stag.next = &tag

			tags = append(tags, tag)
		}

		stag = &tag
	}
	return tags
}