package node

import "testing"

func TestFieldType_String(t *testing.T) {
	tests := []struct {
		name string
		t    FieldType
		want string
	}{
		{"String", String, "String"},
		{"Int", Int, "Int"},
		{"Date", Date, "Date"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("FieldType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
