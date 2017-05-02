package webpage

import (
	"net/http"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/node"
	"github.com/henkburgstra/kamibasami/service"
)

type Webpage struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type Result struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func NewResult() *Result {
	r := new(Result)
	r.Status = 200
	return r
}

func NewDataResult(data interface{}) *Result {
	r := NewResult()
	r.Data = data
	return r
}

func NewErrorResult(err int, msg string) *Result {
	r := new(Result)
	r.Status = err
	r.Error = msg
	return r
}

// Test is a test webpage controller
func Test(svc *service.Service, c *gin.Context) {
	c.String(200, "test")
}

func PostWebpage(svc *service.Service, c *gin.Context) {
	var webpage Webpage
	if c.BindJSON(&webpage) == nil {
		_, err := storePage(svc, webpage.URL, webpage.Path)
		if err != nil {
			c.JSON(http.StatusOK, NewErrorResult(http.StatusInternalServerError, err.Error()))
		} else {
			c.JSON(http.StatusOK, NewResult())
		}
	} else {
		c.JSON(http.StatusOK, NewErrorResult(http.StatusInternalServerError, "Invalid data"))
	}
}

func storePage(svc *service.Service, url string, path string) (node.INode, error) {
	path = node.NormalizePath(path)
	resp, err := soup.Get(url)
	if err != nil {
		return nil, err
	}
	doc := soup.HTMLParse(resp)
	title := url
	titleNode := doc.Find("title")
	if titleNode.Pointer != nil {
		title = strings.Trim(titleNode.Text(), "\n\t ")
	}
	parent, err := node.CreatePath(svc.NodeRepo(), path)
	if err != nil {
		return nil, err
	}
	repo := svc.NodeRepo()
	page, err := repo.GetWithParent(title, parent.ID())
	if err != nil {
		page = node.NewWebpage(nil)
		page.SetName(title)
		page.SetParentID(parent.ID())
	}
	page.SetValue("URL", url)
	repo.Put(page)
	page.Index(svc.Index())
	tags := strings.Split(path, "/")
	if len(tags) > 0 {
		err = repo.SetTags(page.ID(), tags...)
	}
	return page, err
}
