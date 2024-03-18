package gdp

import (
	"bytes"
)

var hasNoEndTags = [18]string{"comment", "php", "empty", "!DOCTYPE", "area", "base", "col", "embed", "param", "source", "track", "meta", "link", "br", "input", "hr", "img", "path"}

func Default(html string) Tag {
	var p Parser
	p.html = html
	p.len = len(p.html)

	var document Tag
	document.Tag = "document"
	document.Childrens = p.parse(&document)
	document.TagAttrs = Attr{nil}
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

	for p.i < p.len {

		c := p.html[p.i]
		p.i++
		if c == until[0] && p.isEqual(until[1:]) {
			break
		}

		buffer.WriteByte(c)
	}
	return buffer.String()
}

func (p *Parser) skipSpace() {
	for p.i < p.len {

		c1 := p.html[p.i]

		if c1 == ' ' || c1 == '\n' || c1 == '\t' {
			p.i++
		} else {
			break
		}
	}
}

func (p *Parser) parseAttr() Attr {

	attrs := make(map[string]*string)
	for {
		isThereValue := false
		var buffer bytes.Buffer
		p.skipSpace()
		for p.i < p.len {

			c1 := p.html[p.i]
			p.i++
			if c1 == '>' || c1 == '=' {
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

			for p.i < p.len {

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
			attrs[name] = &value
		}

		c1 := p.html[p.i]
		p.i++
		if c1 == '>' {
			break
		}
	}

	return Attr{attrs}
}

func (p *Parser) parseTag(tag *Tag) bool {
	if p.isEqual("![CDATA[") {
		p.i += 8
		tag.Tag = "cdata"
		return true
	}

	if p.html[p.i+1] == '/' {
		p.i++
	}

	var buffer bytes.Buffer
	mapAttr := make(map[string]*string)
	attrs := Attr{mapAttr}
	for p.i < p.len {

		c1 := p.html[p.i]
		p.i++

		if c1 == '>' {
			break
		}

		if c1 == ' ' || c1 == '\n' || c1 == '\t' {
			attrs = p.parseAttr()
			break
		}
		buffer.WriteByte(c1)

	}
	name := buffer.String()

	tag.Tag = name
	tag.TagAttrs = attrs
	tag.isEnd = false

	if name[0] == '/' {
		tag.isEnd = true
		tag.Tag = name[1:]
	}

	if name[len(name)-1] == '/' {
		name = name[0 : len(name)-1]
		tag.Tag = name
	}
	return true
}

func (p *Parser) parseContent(first byte, tag *Tag) bool {
	p.i--
	content := p.getUntil("<", first)
	p.i--

	tag.Tag = "empty"
	tag.Contents = content

	return true
}

func (p *Parser) parseComment(tag *Tag) bool {
	p.i += 3

	content := p.getUntil("-->", 0)

	p.i += 2

	tag.Tag = "comment"
	tag.Contents = content

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
	return html == text
}

func isEndTag(tag *Tag) bool {
	for i := 0; i < 18; i++ {
		if tag.Tag == hasNoEndTags[i] {
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

	if p.i >= p.len {
		return false
	}

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

	if tag.Tag == "cdata" {
		tag.Contents = p.parseCdata()
		return true
	}
	name := tag.Tag
	if len(name) >= 4 && name[0:4] == "?xml" {
		p.isXml = true
		return true
	}

	if p.isXml {
		hasNoEndTags[11] = ""
	}

	if isEndTag(tag) || tag.isEnd {
		return true
	}

	if tag.Tag == "script" {
		tag.Contents = p.parseScript()
	} else {
		tag.Childrens = p.parse(tag)
	}

	if tag.Tag == p.current.Tag {
		return true
	}

	var etag Tag

	for p.next(&etag) {
		if tag.Tag == etag.Tag {
			break
		}
	}

	return true
}

func (p *Parser) parse(parent *Tag) []*Tag {
	var tags []*Tag
	eq := 0
	stag := &Tag{}
	for p.i < p.len {

		var tag Tag = Tag{}
		tag.next = &Tag{}
		tag.prev = &Tag{}
		if !p.getTag(&tag) {
			break
		}

		if tag.isEnd && parent.Tag == tag.Tag {
			break
		}

		if !tag.isEnd {

			tag.Eq = eq
			eq++
			tag.prev = stag
			tag.parent = parent
			stag.next = &tag

			tags = append(tags, &tag)
			stag = &tag
		}

	}
	return tags
}
