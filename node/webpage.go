package node

import "github.com/blevesearch/bleve"

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

func (n *Webpage) Index(index bleve.Index) {
	if index == nil {
		// TODO: logging
		return
	}
	index.Index(n.ID(),
		struct {
			Name string
			URL  string
		}{
			Name: n.Name(),
			URL:  Value{value: n.Value("URL")}.String(),
		})
}

func init() {
	RegisterConstructor("webpage", NewWebpage)
}
