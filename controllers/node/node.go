package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/service"
	"github.com/henkburgstra/kamibasami/viewmodels"
)

func Home(svc *service.Service, c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl.html", viewmodels.Viewmodel{Title: "test"})
}
