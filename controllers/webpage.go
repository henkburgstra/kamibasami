package controllers

import (
	"github.com/gin-gonic/gin"
)

// Test is a test webpage controller
func Test(c *gin.Context) {
	c.String(200, "test")
}

func init() {
	RegisterController("webpage", "GET", "/webpage/test", Test)
}
