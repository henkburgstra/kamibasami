package node

import (
	"strings"
)

// INode is the interface for Node types.
type INode interface {
	ID() string
	SetID(string)
	Name() string
	SetName(string)
	ParentID() string
	SetParentID(string)
	Fields() []Field
}

// INodeRepo is the interface for Node repositories.
type INodeRepo interface {
	Get(id string) (node INode, err error)
	GetWithParent(name string, parent string) (node INode, err error)
	GetChildren(id string) []INode
	Put(node INode) (err error)
}

// Node is the basic implementation of a Node structure.
type Node struct {
	id       string
	name     string
	parentID string
}

// ID returns the ID of a Node.
func (n *Node) ID() string {
	return n.id
}

// SetID sets the id of a Node.
func (n *Node) SetID(id string) {
	n.id = id
}

// Name returns the name of a Node.
func (n *Node) Name() string {
	return n.name
}

// SetName sets the name of a Node.
func (n *Node) SetName(name string) {
	n.name = name
}

// ParentID returns the parent id of a Node.
func (n *Node) ParentID() string {
	return n.parentID
}

// SetParentID sets the parent id of a Node.
func (n *Node) SetParentID(id string) {
	n.parentID = id
}

func (n *Node) Fields() []Field {
	return make([]Field, 0)
}

func NewNode(id string, name string, parentID string) *Node {
	return &Node{id: id, name: name, parentID: parentID}
}

// MockNodeRepo mocks the INodeRepo interface.
type MockNodeRepo struct {
	nodesByID     map[string]INode
	nodesByParent map[string]map[string]INode
}

// NewMockNodeRepo creates a new MockNodeRepo instance.
func NewMockNodeRepo() *MockNodeRepo {
	r := new(MockNodeRepo)
	r.nodesByID = make(map[string]INode)
	r.nodesByParent = make(map[string]map[string]INode)
	return r
}

// Put implements INodeRepo.Put.
func (r *MockNodeRepo) Put(node INode) (err error) {
	r.nodesByID[node.ID()] = node
	nodes := r.nodesByParent[node.Name()]
	if nodes == nil {
		nodes = make(map[string]INode)
		r.nodesByParent[node.Name()] = nodes
	}
	nodes[node.ParentID()] = node
	return
}

func (r *MockNodeRepo) Get(id string) (node INode, err error) {
	node, ok := r.nodesByID[id]
	if !ok {
		err = NewNotFoundError(id)
	}
	return
}

func (r *MockNodeRepo) GetChildren(id string) []INode {
	c := make([]INode, 0)
	return c
}

// GetWithParent implements INodeRepo.GetWithParent.
func (r *MockNodeRepo) GetWithParent(name string, parent string) (node INode, err error) {
	nodes := r.nodesByParent[name]
	if nodes == nil {
		err = NewNotFoundError(name)
		return
	}
	node = nodes[parent]
	if node == nil {
		err = NewNotFoundError(name)
	}
	return
}

func NormalizePath(path string) (path string) {
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Trim(path, "\\/")
	return
}

// GetWithPath returns the Node associated with a path.
func GetWithPath(r INodeRepo, path string) (node INode, err error) {
	path = NormalizePath(path)
	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		err = NewInvalidPathError(path)
		return
	}

	var parentID string

	for i := 0; i < len(parts); i++ {
		name := parts[i]
		node, err = r.GetWithParent(name, parentID)
		if err != nil {
			err = NewInvalidPathError(path)
			return
		}
		parentID = node.ID()
	}

	return
}

// GetPath returns the path of a Node.
func GetPath(r INodeRepo, node INode) (path string) {
	l := make([]string, 0)
	for {
		l = append(l, node.Name())
		if node.ParentID() == "" {
			break
		}
		var err error
		node, err = r.Get(node.ParentID())
		if err != nil {
			break
		}
	}
	// reverse the path slice
	for i := len(l)/2 - 1; i >= 0; i-- {
		opp := len(l) - 1 - i
		l[i], l[opp] = l[opp], l[i]
	}
	path = strings.Join(l, "/")
	return
}

// CreatePath creates all missing Nodes in a path
// returns the last Node
func CreatePath(r INodeRepo, path string) (node INode, err error) {
	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		err = NewInvalidPathError(path)
		return
	}

	var parentID string

	for i := 0; i < len(parts); i++ {
		name := parts[i]
		node, err = r.GetWithParent(name, parentID)
		if err != nil {
			// TODO: new uuid
			node = NewNode("id", name, parentID)
			err = r.Put(node)
			if err != nil {
				return
			}
		}
		parentID = node.ID()
	}

	return
}
