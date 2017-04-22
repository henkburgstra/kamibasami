package node

type Webpage struct {
	*Node
}

func NewWebpage(id string, name string, parentID string) *Webpage {
	w := new(Webpage)
	w.Node = NewNode(id, name, parentID)
	return w
}

func (n *Webpage) Fields() []Field {
	return []Field{
		Field{Name: "URL"},
	}
}
