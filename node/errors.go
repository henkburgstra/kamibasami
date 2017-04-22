package node

import (
	"fmt"
)

// NotFoundError indicates that a given Node does not exist
type NotFoundError struct {
	id string
}

// ID returns the id of the requested Node
func (e *NotFoundError) ID() string {
	return e.id
}

// Error is an implementation of the error interface
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Node '%s' not found", e.id)
}

// NewNotFoundError returns a NotFoundError instance
func NewNotFoundError(id string) *NotFoundError {
	e := new(NotFoundError)
	e.id = id
	return e
}

// InvalidPathError indicates that a given path does not exist or is invalid
type InvalidPathError struct {
	path string
}

// Path returns the requested path
func (e *InvalidPathError) Path() string {
	return e.path
}

// Error is an implementation of the error interface
func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("Path '%s' is invalid ", e.path)
}

// NewInvalidPathError returns a InvalidPath instance
func NewInvalidPathError(path string) *InvalidPathError {
	e := new(InvalidPathError)
	e.path = path
	return e
}
