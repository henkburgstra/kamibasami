package node

import (
	"reflect"
	"testing"
)

func TestNewMockNodeRepo(t *testing.T) {
	r := NewMockNodeRepo()
	if r == nil {
		t.Error("NewMockNodeRepo returns nil")
	}
}

func TestNode_ID(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want string
	}{
		// TODO: Add test cases.
		{"id", &Node{id: "id"}, "id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ID(); got != tt.want {
				t.Errorf("Node.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_Name(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want string
	}{
		{"name", &Node{name: "name"}, "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Name(); got != tt.want {
				t.Errorf("Node.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_ParentID(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want string
	}{
		{"parentID", &Node{parentID: "parentID"}, "parentID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ParentID(); got != tt.want {
				t.Errorf("Node.ParentID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockNodeRepo_Put(t *testing.T) {
	node := new(Node)
	node.SetID("1")
	tests := []struct {
		name    string
		r       *MockNodeRepo
		node    INode
		wantErr bool
	}{
		{"put", NewMockNodeRepo(), node, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Put(tt.node); (err != nil) != tt.wantErr {
				t.Errorf("MockNodeRepo.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockNodeRepo_GetWithParent(t *testing.T) {
	type args struct {
		name   string
		parent string
	}
	repo := NewMockNodeRepo()
	node1 := &Node{id: "1", name: "parent"}
	repo.Put(node1)
	node2 := &Node{id: "2", name: "child", parentID: "1"}
	repo.Put(node2)
	node3 := &Node{id: "3", name: "grandchild", parentID: "2"}
	repo.Put(node3)
	tests := []struct {
		name     string
		r        *MockNodeRepo
		args     args
		wantNode INode
		wantErr  bool
	}{
		{"een niveau", repo, args{"child", "1"}, node2, false},
		{"twee niveaus", repo, args{"grandchild", "2"}, node3, false},
		{"verkeerde parent", repo, args{"grandchild", "1"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNode, err := tt.r.GetWithParent(tt.args.name, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockNodeRepo.GetWithParent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("MockNodeRepo.GetWithParent() = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}

func TestGetWithPath(t *testing.T) {
	type args struct {
		r    INodeRepo
		path string
	}
	repo := NewMockNodeRepo()
	node1 := &Node{id: "1", name: "parent"}
	repo.Put(node1)
	node2 := &Node{id: "2", name: "child", parentID: "1"}
	repo.Put(node2)
	node3 := &Node{id: "3", name: "grandchild", parentID: "2"}
	repo.Put(node3)
	tests := []struct {
		name     string
		args     args
		wantNode INode
		wantErr  bool
	}{
		{"valid path", args{repo, "parent/child/grandchild"}, node3, false},
		{"valid path2", args{repo, "parent/child"}, node2, false},
		{"invalid path", args{repo, "parent/child/something"}, nil, true},
		{"backslashes", args{repo, "parent\\child\\grandchild"}, node3, false},
		{"leading slash", args{repo, "/parent/child/grandchild"}, node3, false},
		{"trailing slash", args{repo, "parent/child/grandchild/"}, node3, false},
		{"leading and trailing slash", args{repo, "/parent/child/grandchild/"}, node3, false},
		{"leading backslash", args{repo, "\\parent\\child\\grandchild"}, node3, false},
		{"trailing backslash", args{repo, "parent\\child\\grandchild\\"}, node3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNode, err := GetWithPath(tt.args.r, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWithPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("GetWithPath() = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}

func TestGetPath(t *testing.T) {
	type args struct {
		r    INodeRepo
		node INode
	}
	repo := NewMockNodeRepo()
	node1 := &Node{id: "1", name: "parent"}
	repo.Put(node1)
	node2 := &Node{id: "2", name: "child", parentID: "1"}
	repo.Put(node2)
	node3 := &Node{id: "3", name: "grandchild", parentID: "2"}
	repo.Put(node3)
	tests := []struct {
		name     string
		args     args
		wantPath string
	}{
		{"grandchild", args{repo, node3}, "parent/child/grandchild"},
		{"child", args{repo, node2}, "parent/child"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPath := GetPath(tt.args.r, tt.args.node); gotPath != tt.wantPath {
				t.Errorf("GetPath() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
