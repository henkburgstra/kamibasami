package node

import (
	"strings"

	"database/sql"

	"github.com/satori/go.uuid"
)

func checkerr(e error) {
	if e != nil {
		panic(e)
	}
}

type Constructor func(node INode) INode

var constructors map[string]Constructor = make(map[string]Constructor)

func RegisterConstructor(nodeType string, constructor Constructor) {
	constructors[nodeType] = constructor
}

func GetConstructor(nodeType string) Constructor {
	return constructors[nodeType]
}

func Transform(node INode) INode {
	constructor, ok := constructors[node.Type()]
	if ok {
		return constructor(node)
	}
	return node
}

// INode is the interface for Node types.
type INode interface {
	ID() string
	SetID(string)
	Type() string
	SetType(string)
	Name() string
	SetName(string)
	ParentID() string
	SetParentID(string)
	Fields() []Field
	Value(name string) interface{}
	SetValue(name string, value interface{})
}

// INodeRepo is the interface for Node repositories.
type INodeRepo interface {
	Get(id string) (node INode, err error)
	GetWithParent(name string, parent string) (node INode, err error)
	GetChildren(id string) (nodes []INode, err error)
	Put(node INode) (err error)
}

// Node is the basic implementation of a Node structure.
type Node struct {
	id       string
	nodeType string
	name     string
	parentID string
	values   map[string]interface{}
}

// ID returns the ID of a Node.
func (n *Node) ID() string {
	return n.id
}

// SetID sets the id of a Node.
func (n *Node) SetID(id string) {
	n.id = id
}

func (n *Node) Type() string {
	return n.nodeType
}

func (n *Node) SetType(nodeType string) {
	n.nodeType = nodeType
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

func (n *Node) Value(name string) interface{} {
	value, ok := n.values[name]
	if ok {
		return value
	}
	return nil
}

func (n *Node) SetValue(name string, value interface{}) {
	n.values[name] = value
}

func NewNode(name string) *Node {
	return &Node{name: name, values: make(map[string]interface{})}
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

func (r *MockNodeRepo) Get(id string) (INode, error) {
	node, ok := r.nodesByID[id]
	if !ok {
		return nil, NewNotFoundError(id)
	}
	return Transform(node), nil
}

func (r *MockNodeRepo) GetChildren(id string) (nodes []INode, err error) {
	nodes = make([]INode, 0)
	return
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

// Put implements INodeRepo.Put.
func (r *MockNodeRepo) Put(node INode) (err error) {
	if node.ID() == "" {
		node.SetID(uuid.NewV4().String())
	}
	r.nodesByID[node.ID()] = node
	nodes := r.nodesByParent[node.Name()]
	if nodes == nil {
		nodes = make(map[string]INode)
		r.nodesByParent[node.Name()] = nodes
	}
	nodes[node.ParentID()] = node
	return
}

// NormalizePath makes sure that path delimiters are slashes and not backslashes
// and that path doesn't start or end with a slash
func NormalizePath(p string) (path string) {
	path = strings.Replace(p, "\\", "/", -1)
	path = strings.Trim(path, "\\/")
	return
}

// GetWithPath returns the Node associated with a path.
func GetWithPath(r INodeRepo, path string) (INode, error) {
	var node INode
	path = NormalizePath(path)
	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		return nil, NewInvalidPathError(path)
	}

	var parentID string

	for i := 0; i < len(parts); i++ {
		name := parts[i]
		var err error
		node, err = r.GetWithParent(name, parentID)
		if err != nil {
			return nil, NewInvalidPathError(path)
		}
		parentID = node.ID()
	}

	if node == nil {
		return nil, NewNotFoundError("")
	}
	return Transform(node), nil
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
// and returns the last Node
func CreatePath(r INodeRepo, path string) (INode, error) {
	var node INode
	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		return nil, NewInvalidPathError(path)
	}

	var parentID string

	for i := 0; i < len(parts); i++ {
		name := parts[i]
		var err error
		node, err = r.GetWithParent(name, parentID)
		if err != nil {
			node = NewNode(name)
			node.SetParentID(parentID)
			err = r.Put(node)
			if err != nil {
				return nil, err
			}
		}
		parentID = node.ID()
	}

	if node == nil {
		return nil, NewNotFoundError("")
	}
	return Transform(node), nil
}

type DBNodeRepo struct {
	db     *sql.DB
	dbType string
}

func NewDBNodeRepo(db *sql.DB, dbType string) *DBNodeRepo {
	return &DBNodeRepo{db: db, dbType: dbType}
}

func (r *DBNodeRepo) Get(id string) (INode, error) {
	row := r.db.QueryRow(`SELECT node_id, node_type, node_name, parent_id, node_values
		FROM node
		WHERE node_id = ?`, id)
	nodeID := ""
	nodeType := ""
	nodeName := ""
	parentID := ""
	nodeValues := ""
	err := row.Scan(&nodeID, &nodeType, &nodeName, &parentID, &nodeValues)
	if err != nil {
		return nil, err
	}
	node := NewNode(nodeName)
	node.SetID(nodeID)
	node.SetParentID(parentID)
	node.SetType(nodeType)
	return Transform(node), nil
}

func (r *DBNodeRepo) Put(node INode) (err error) {
	// TODO: node.Values()
	if node.ID() == "" {
		node.SetID(uuid.NewV4().String())
		_, err = r.db.Exec(`INSERT INTO node
		(node_id, node_type, node_name, parent_id, node_values)
		VALUES
		(?, ?, ?, ?, ?)	`, node.ID(), node.Type(), node.Name(), node.ParentID(), "")
	} else {
		_, err = r.db.Exec(`UPDATE node
		SET node_type = ?, node_name = ?, parent_id = ?, node_values = ?
		WHERE node_id = ?`, node.Type(), node.Name(), node.ParentID(), "", node.ID())
	}
	return
}

func (r *DBNodeRepo) GetChildren(id string) (nodes []INode, err error) {
	nodes = make([]INode, 0)
	rows, err := r.db.Query(`SELECT node_id, node_type, node_name, node_values
		FROM node
		WHERE node_id = ?`, id)
	if err != nil {
		return
	}
	defer rows.Close()

	nodeID := ""
	nodeType := ""
	nodeName := ""
	nodeValues := ""

	for rows.Next() {
		if err = rows.Scan(&nodeID, &nodeType, &nodeName, &nodeValues); err != nil {
			return
		}
		node := NewNode(nodeName)
		node.SetID(nodeID)
		node.SetParentID(id)
		node.SetType(nodeType)
		nodes = append(nodes, Transform(node))
	}

	return
}

func (r *DBNodeRepo) GetWithParent(name string, parent string) (INode, error) {
	row := r.db.QueryRow(`SELECT node_id, node_type, node_name, node_values
		FROM node
		WHERE node_name = ?
		AND parent_id = ?`, name, parent)
	nodeID := ""
	nodeType := ""
	nodeName := ""
	nodeValues := ""
	err := row.Scan(&nodeID, &nodeType, &nodeName, &nodeValues)
	if err != nil {
		return nil, err
	}
	node := NewNode(nodeName)
	node.SetID(nodeID)
	node.SetParentID(parent)
	node.SetType(nodeType)
	return Transform(node), nil
}
