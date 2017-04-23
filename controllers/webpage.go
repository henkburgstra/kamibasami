package controllers

import (
	"github.com/anaskhan96/soup"
	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/node"
	"github.com/henkburgstra/kamibasami/service"
)

// Test is a test webpage controller
func Test(svc *service.Service, c *gin.Context) {
	c.String(200, "test")
}

func storePage(svc *service.Service, url string, path string) (err error) {
	resp, err := soup.Get(url)
	if err != nil {
		return
	}
	doc := soup.HTMLParse(resp)
	title := url
	titleNode := doc.Find("title")
	if titleNode.Pointer != nil {
		title = titleNode.Text()
	}
	parent := node.CreatePath(path)
	// TODO uuid
	page := node.NewWebpage("id", title, parent.ID())
	svc.NodeRepo().Put(page)
	return
}

func init() {
	RegisterController("webpage", "GET", "/webpage/test", Test)
}
