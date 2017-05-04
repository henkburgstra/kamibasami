package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/service"
)

func Home(svc *service.Service, c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl.html", gin.H{
		"title": "Posts",
	})
}
