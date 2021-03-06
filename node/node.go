package node

import (
	"fmt"
	"strings"

	"database/sql"
	"encoding/json"

	"github.com/blevesearch/bleve"

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
	Index(index bleve.Index)
}

// INodeRepo is the interface for Node repositories.
type INodeRepo interface {
	Get(id string) (node INode, err error)
	GetWithParent(name string, parent string) (node INode, err error)
	GetChildren(id string) (nodes []INode, err error)
	Put(node INode) (err error)
	Delete(id string) error
	SetTags(id string, tags ...string) error
	UnsetTags(id string, tags ...string) error
	DeleteTag(tag string) error
	UpdateTagsWithPath(path string, node INode) error
}

// Node is the basic implementation of a Node structure.
type Node struct {
	id       string
	nodeType string
	name     string
	parentID string
	fields   []Field
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
	return n.fields
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

func (n *Node) Index(index bleve.Index) {
	if index == nil {
		// TODO: logging
		return
	}
	index.Index(n.ID(),
		struct {
			Type string
			Name string
		}{
			Type: n.Type(),
			Name: n.Name(),
		})
	fields := make(map[string]string)
	for _, field := range n.Fields() {
		v := Value{value: n.Value(field.Name)}.String()
		if v != "" && v != "NULL" {
			fields[field.Name] = v
		}
		if len(fields) > 0 {
			index.Index(n.ID(), fields)
		}
	}
}

func NewNode() *Node {
	return &Node{fields: make([]Field, 0), values: make(map[string]interface{})}
}

// MockNodeRepo mocks the INodeRepo interface.
type MockNodeRepo struct {
	nodesByID     map[string]INode
	nodesByParent map[string]map[string]INode
	tags          map[string]bool
	tagsByNode    map[string]map[string]bool
}

// NewMockNodeRepo creates a new MockNodeRepo instance.
func NewMockNodeRepo() *MockNodeRepo {
	r := new(MockNodeRepo)
	r.nodesByID = make(map[string]INode)
	r.nodesByParent = make(map[string]map[string]INode)
	r.tags = make(map[string]bool)
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

func (r *MockNodeRepo) Delete(id string) error {
	node, ok := r.nodesByID[id]
	if !ok {
		return nil
	}
	delete(r.nodesByID, id)
	delete(r.nodesByParent[node.Name()], node.ParentID())
	return nil
}

func (r *MockNodeRepo) SetTags(id string, tags ...string) error {
	for _, tag := range tags {
		r.tags[tag] = true
		t, ok := r.tagsByNode[id]
		if !ok {
			t = make(map[string]bool)
			r.tagsByNode[id] = t
		}
		t[tag] = true
	}
	return nil
}

func (r *MockNodeRepo) UnsetTags(id string, tags ...string) error {
	for _, tag := range tags {
		t, ok := r.tagsByNode[id]
		if ok {
			delete(t, tag)
		}
	}
	return nil
}

func (r *MockNodeRepo) DeleteTag(tag string) error {
	for _, t := range r.tagsByNode {
		delete(t, tag)
	}
	delete(r.tags, tag)
	return nil
}

func (r *MockNodeRepo) UpdateTagsWithPath(path string, node INode) error {
	return UpdateTagsWithPath(r, path, node.ID())
}

// NormalizePath makes sure that path delimiters are slashes and not backslashes
// and that path doesn't start or end with a slash
func NormalizePath(p string) (path string) {
	path = strings.Replace(p, "\\", "/", -1)
	path = strings.Trim(path, "\\/")
	return
}

func Values2Json(node INode) (string, error) {
	n := Transform(node)
	values := make(map[string]interface{})
	for _, field := range n.Fields() {
		values[field.Name] = n.Value(field.Name)
	}
	b, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
			node = NewNode()
			node.SetType("folder")
			node.SetName(name)
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
func UpdateTagsWithPath(repo INodeRepo, path string, id string) error {
	var err error
	tags := strings.Split(path, "/")
	if len(tags) > 0 {
		err = repo.SetTags(id, tags...)
	}
	return err
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
	node := NewNode()
	node.SetName(nodeName)
	node.SetID(nodeID)
	node.SetParentID(parentID)
	node.SetType(nodeType)
	v := make(map[string]interface{})
	err = json.Unmarshal([]byte(nodeValues), v)
	if err != nil {
		// TODO: logging
	} else {
		for key, value := range v {
			node.SetValue(key, value)
		}
	}
	return Transform(node), nil
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
		node := NewNode()
		node.SetName(nodeName)
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
	node := NewNode()
	node.SetName(nodeName)
	node.SetID(nodeID)
	node.SetParentID(parent)
	node.SetType(nodeType)
	return Transform(node), nil
}

func (r *DBNodeRepo) Put(node INode) (err error) {
	// TODO: node.Values()
	values, err := Values2Json(node)
	if err != nil {
		// TODO: logging
	}
	if node.ID() == "" {
		node.SetID(uuid.NewV4().String())
		_, err = r.db.Exec(`INSERT INTO node
		(node_id, node_type, node_name, parent_id, node_values)
		VALUES
		(?, ?, ?, ?, ?)	`, node.ID(), node.Type(), node.Name(), node.ParentID(), values)
	} else {
		_, err = r.db.Exec(`UPDATE node
		SET node_type = ?, node_name = ?, parent_id = ?, node_values = ?
		WHERE node_id = ?`, node.Type(), node.Name(), node.ParentID(), values, node.ID())
	}
	return
}

func (r *DBNodeRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM node WHERE node_id = ?`, id)
	return err
}

func (r *DBNodeRepo) SetTags(id string, tags ...string) error {
	for _, tag := range tags {
		var x int
		err := r.db.QueryRow("SELECT 1 FROM tag WHERE tag_name = ?", tag).Scan(&x)
		switch {
		case err == sql.ErrNoRows:
			_, err = r.db.Exec(`INSERT INTO tag (tag_name) VALUES (?)`, tag)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		case err != nil:
			fmt.Println(err.Error())
			return err
		}
		var nodeTagsID int
		err = r.db.QueryRow("SELECT node_tags_id FROM node_tags WHERE node_id = ? AND tag_name = ?", id, tag).Scan(&nodeTagsID)
		switch {
		case err == sql.ErrNoRows:
			_, err = r.db.Exec(`INSERT INTO node_tags (node_id, tag_name) VALUES (?, ?)`, id, tag)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		case err != nil:
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func (r *DBNodeRepo) UnsetTags(id string, tags ...string) error {
	for _, tag := range tags {
		_, err := r.db.Exec(`DELETE FROM node_tags WHERE node_id = ? AND tag_name = ?`, id, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DBNodeRepo) DeleteTag(tag string) error {
	_, err := r.db.Exec(`DELETE FROM node_tags WHERE tag_name = ?`, tag)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM tag WHERE tag_name = ?`, tag)
	return err
}

func (r *DBNodeRepo) UpdateTagsWithPath(path string, node INode) error {
	return UpdateTagsWithPath(r, path, node.ID())
}
