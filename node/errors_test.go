package node

import (
	"reflect"
	"testing"
)

func TestNotFoundError_Id(t *testing.T) {
	tests := []struct {
		name string
		e    *NotFoundError
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ID(); got != tt.want {
				t.Errorf("NodeNotFoundError.Id() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *NotFoundError
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("NodeNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want *NotFoundError
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundError(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidPathError_Path(t *testing.T) {
	tests := []struct {
		name string
		e    *InvalidPathError
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Path(); got != tt.want {
				t.Errorf("InvalidPathError.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidPathError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *InvalidPathError
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("InvalidPathError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInvalidPathError(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *InvalidPathError
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInvalidPathError(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvalidPathError() = %v, want %v", got, tt.want)
			}
		})
	}
}
