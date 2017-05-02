package node

import (
	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/service"
)

func Home(svc *service.Service, c *gin.Context) {
	c.String(200, "Home")
}
