package node

type Webpage struct {
	INode
}

func NewWebpage(id string, name string, parentID string) *Webpage {
	w := new(Webpage)
	w.INode = NewNode(id, name, parentID)
	w.SetType("webpage")
	return w
}

func NewWebNode(node INode) INode {
	w := new(Webpage)
	w.INode = node
	w.SetType("webpage")
	return w
}

func (n *Webpage) Fields() []Field {
	return []Field{
		Field{Name: "URL", Type: String},
	}
}

func init() {
	RegisterConstructor("webpage", NewWebNode)
}
