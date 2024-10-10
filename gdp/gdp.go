package gdp

import (
	"bytes"
	"unicode"
)

type Parser struct {
	i    int
	len  int
	html string
}

func Default(html string) Tag {
	var p Parser
	p.html = html
	p.len = len(html)

	var document Tag
	document.Name = "document"
	document.Children = p.parse(&document, &document)
	document.Attrs = Attr{nil}

	return document
}

var noEndTags = [18]string{"comment", "empty", "!DOCTYPE", "area", "base", "col", "embed", "param", "source", "track", "meta", "link", "br", "input", "hr", "img", "path"}

func (p *Parser) parseEndTag() Tag {

	var buffer bytes.Buffer
	for p.i < p.len {
		tok := p.html[p.i]
		p.i++
		if tok == '>' {
			break
		}
		buffer.WriteByte(tok)
	}
	name := buffer.String()

	tag := Tag{}
	tag.Name = name
	return tag
}

func (p *Parser) skipSpace() {
	for p.i < p.len {

		tok := p.html[p.i]

		if tok == ' ' || tok == '\n' || tok == '\t' {
			p.i++
		} else {
			break
		}
	}
}

func (p *Parser) parseAttrs() Attr {

	attrs := make(map[string]*string)
	for {
		isThereValue := false
		var buffer bytes.Buffer
		p.skipSpace()
		for p.i < p.len {

			tok := p.html[p.i]
			if tok == '>' || tok == '=' {
				if tok == '=' {
					p.i++
					isThereValue = true
				}

				break
			}
			p.i++
			buffer.WriteByte(tok)
		}

		name := buffer.String()

		var buffer1 bytes.Buffer
		if isThereValue {
			tok := p.html[p.i]
			var t byte = 0
			if tok == '"' || tok == '\'' {
				t = tok
				p.i++
			}

			for p.i < p.len {

				tok := p.html[p.i]
				if tok == t {
					break
				}

				if t == 0 && tok == ' ' {
					break
				}

				if t == 0 && tok == '>' {
					break
				}
				p.i++
				buffer1.WriteByte(tok)
			}
		}

		value := buffer1.String()

		if len(name) > 0 && name[0] != '/' && name[0] != ' ' {
			attrs[name] = &value
		}

		tok := p.html[p.i]
		if tok == '>' {
			p.i++
			break
		}
		p.i++
	}

	return Attr{attrs}
}

func (p *Parser) getTag(tag *Tag) {
	var buffer bytes.Buffer
	mapAttr := make(map[string]*string)
	attrs := Attr{mapAttr}

	for p.i < p.len {
		tok := p.html[p.i]
		if tok == '>' {
			p.i++
			break
		}

		if tok == ' ' || tok == '\n' || tok == '\t' {
			attrs = p.parseAttrs()
			break
		}

		p.i++
		buffer.WriteByte(tok)
	}

	name := buffer.String()

	tag.Name = name
	tag.Attrs = attrs
}

func inStringArray(str string, array []string) bool {
	for i := 0; i < len(array); i++ {
		if str == array[i] {
			return true
		}
	}
	return false
}

func (p *Parser) parseTag(tag *Tag) {
	p.getTag(tag)
	if tag.Name == "script" {
		p.parseScript(tag)
		return
	}

	if inStringArray(tag.Name, noEndTags[:]) {
		return
	}

	var etag Tag
	var tags []*Tag
	for p.i < p.len {
		tags = p.parse(tag, &etag)

		if tag.Name == etag.Name {
			tag.Children = tags
			return
		}
	}
	tag.Children = tags
	return
}

func (p *Parser) isEqual(text string) bool {
	textLen := len(text)
	if p.i+textLen > p.len {
		return false
	}

	html := p.html[p.i : p.i+textLen]
	return html == text
}

func (p *Parser) getUntil(until string) string {
	var buffer bytes.Buffer

	for p.i < p.len {

		tok := p.html[p.i]

		if tok == until[0] && p.isEqual(until) {
			break
		}
		p.i++

		buffer.WriteByte(tok)
	}

	return buffer.String()
}

func (p *Parser) parseComment(tag *Tag) {

	p.i += 3
	content := p.getUntil("-->")
	p.i += 3

	tag.Name = "comment"
	tag.Content = content
}

func (p *Parser) parseContent(tag *Tag) {

	content := p.getUntil("<")

	tag.Name = "empty"
	tag.Content = content
}

func (p *Parser) parseCData(tag *Tag) {
	p.i += 8
	content := p.getUntil("]]>")
	p.i += 3

	tag.Name = "cdata"
	tag.Content = content
}

func (p *Parser) parseScript(tag *Tag) {
	content := p.getUntil("</script")
	p.i += 9

	tag.Name = "script"
	tag.Content = content
}

func (p *Parser) parse(parent *Tag, etag *Tag) []*Tag {
	var tags []*Tag

	prev := &Tag{}
	eq := 0
	for p.i < p.len {
		tok := p.html[p.i]
		var tag Tag = Tag{}
		tag.parent = &Tag{}
		tag.next = &Tag{}
		tag.prev = &Tag{}
		tag.Children = []*Tag{}
		if tok == '<' {
			p.i++
			next := p.html[p.i]

			if next == '!' {
				if p.isEqual("![CDATA[") {

					p.parseCData(&tag)
				} else if p.isEqual("!--") {
					p.parseComment(&tag)
				} else {
					p.parseContent(&tag)
				}
			} else if next == '/' {

				p.i++
				tag = p.parseEndTag()
				etag.Name = tag.Name

				break
			} else if !unicode.IsLetter(rune(next)) {
				p.parseContent(&tag)
				tag.Content = "<" + tag.Content
			} else {

				p.parseTag(&tag)
			}
		} else {
			p.parseContent(&tag)
		}

		tag.Eq = eq
		eq++
		tag.prev = prev
		tag.parent = parent
		prev.next = &tag

		tags = append(tags, &tag)
		prev = &tag

	}

	return tags
}
