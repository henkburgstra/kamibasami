package webpage

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/node"
	"github.com/henkburgstra/kamibasami/service"
)

func startServer() {
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.String(200, `
		<html>
			<head><title>Testonderwerp</title></head>
			<body>
				<h1>Eerste hoofdstuk</h1>
				<p>Tekst</p>
				<h2>Tweede hoofdstuk</h2>
				<p>Meer tekst</p>
			</body>
		</html>`)
	})
	router.Run()
}

func Test_storePage(t *testing.T) {
	svc := service.NewService()
	svc.SetNodeRepo(node.NewMockNodeRepo())

	go startServer()
	time.Sleep(3 * time.Second)

	type args struct {
		svc  *service.Service
		url  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid path", args{svc: svc, url: "http://localhost:8080/test", path: "a/test/path"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, err := storePage(tt.args.svc, tt.args.url, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("storePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPage.Name() != "Testonderwerp" {
				t.Errorf("storePage() page name = '%s', want 'Testonderwerp'", gotPage.Name())
				return
			}
		})
	}
}
