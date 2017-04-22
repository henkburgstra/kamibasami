package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/service"
)

// Test is a test webpage controller
func Test(svc *service.Service, c *gin.Context) {
	c.String(200, "test")
}

func init() {
	RegisterController("webpage", "GET", "/webpage/test", Test)
}
