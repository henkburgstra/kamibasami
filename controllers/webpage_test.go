package controllers

import (
	"reflect"
	"testing"

	"github.com/henkburgstra/kamibasami/node"
	"github.com/henkburgstra/kamibasami/service"
)

func Test_storePage(t *testing.T) {
	type args struct {
		svc  *service.Service
		url  string
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantPage *node.Webpage
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, err := storePage(tt.args.svc, tt.args.url, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("storePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPage, tt.wantPage) {
				t.Errorf("storePage() = %v, want %v", gotPage, tt.wantPage)
			}
		})
	}
}
