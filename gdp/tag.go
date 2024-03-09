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
