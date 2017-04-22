package node

import (
	"fmt"
)

type NodeNotFoundError struct {
	id string
}

func (e *NodeNotFoundError) Id() string {
	return e.id
}

func (e *NodeNotFoundError) Error() string {
	return fmt.Sprintf("Node '%s' not found", e.id)
}

func NewNodeNotFoundError(id string) *NodeNotFoundError {
	e := new(NodeNotFoundError)
	e.id = id
	return e
}

type InvalidPathError struct {
	path string
}

func (e *InvalidPathError) Path() string {
	return e.path
}

func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("Path '%s' is invalid ", e.path)
}

func NewInvalidPathError(path string) *InvalidPathError {
	e := new(InvalidPathError)
	e.path = path
	return e
}
