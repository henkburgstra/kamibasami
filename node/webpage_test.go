package node

import (
	"reflect"
	"testing"
)

func TestWebpage_Fields(t *testing.T) {
	type fields struct {
		Node *Node
	}
	tests := []struct {
		name   string
		fields fields
		want   []Field
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Webpage{
				INode: tt.fields.Node,
			}
			if got := n.Fields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Webpage.Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}
