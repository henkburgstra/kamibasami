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

func storePage(svc *service.Service, url string, path string) (node.INode, error) {
	resp, err := soup.Get(url)
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(resp)
	title := url
	titleNode := doc.Find("title")
	if titleNode.Pointer != nil {
		title = titleNode.Text()
	}
	parent, err := node.CreatePath(svc.NodeRepo(), path)
	if err != nil {
		return nil, err
	}
	page := node.NewWebpage(nil)
	page.SetName(title)
	page.SetParentID(parent.ID())
	page.SetValue("URL", url)
	svc.NodeRepo().Put(page)
	page.Index(svc.Index())
	return page, nil
}

func init() {
	RegisterController("webpage", "GET", "/webpage/test", Test)
}
