package gdp

type NodeList struct {
	list []*Tag
}

func (n *NodeList) Eq(index int) *Tag {
	for i, tag := range n.list {
		if i == index {
			return tag
		}
	}
	return &Tag{}
}

func (n *NodeList) Each(callback func(int, *Tag)) {
	for index, tag := range n.list {
		callback(index, tag)
	}
}

func (n *NodeList) Html() string {
	ret := ""
	n.Each(func(i int, t *Tag) {
		ret += t.OuterHtml()
	})
	return ret
}
