package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/henkburgstra/kamibasami/controllers/node"
	"github.com/henkburgstra/kamibasami/controllers/webpage"
	"github.com/henkburgstra/kamibasami/service"
)

type ControllerFunc func(service *service.Service, c *gin.Context)

type Controller struct {
	Method  string
	URI     string
	Handler ControllerFunc
}

var ctrls = make([]Controller, 0)
var ctrlIndex = make(map[string]string)

func Get() []Controller {
	return ctrls
}

// Register is used to register a controller
func Register(module string, method string, uri string, handler ControllerFunc) (err error) {
	var key = fmt.Sprintf("%s:%s", method, uri)
	if module, ok := ctrlIndex[key]; ok {
		return fmt.Errorf("Controller with method '%s' and uri '%s' already registered by module '%s'",
			module, method, uri)
	}
	ctrlIndex[key] = module
	ctrls = append(ctrls, Controller{method, uri, handler})
	return
}

func init() {
	Register("node", "GET", "/", node.Home)
	Register("webpage", "GET", "/webpage/test", webpage.Test)
	Register("webpage", "POST", "/api/webpage", webpage.PostWebpage)
}
