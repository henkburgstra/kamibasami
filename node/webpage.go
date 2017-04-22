package node

type Webpage struct {
	*Node
}

func (n *Webpage) Fields() []Field {
	return []Field{
		Field{Name: "URL"},
	}
}
