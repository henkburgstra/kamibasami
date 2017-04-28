package node

type Webpage struct {
	INode
}

func NewWebpage(node INode) INode {
	w := new(Webpage)
	if node == nil {
		w.INode = NewNode()
	} else {
		w.INode = node
	}
	w.SetType("webpage")
	return w
}

func (n *Webpage) Fields() []Field {
	return []Field{
		Field{Name: "URL", Type: String},
	}
}

func init() {
	RegisterConstructor("webpage", NewWebpage)
}
