package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
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

// RegisterController is used to register a controller
func RegisterController(module string, method string, uri string, handler ControllerFunc) (err error) {
	var key = fmt.Sprintf("%s:%s", method, uri)
	if module, ok := ctrlIndex[key]; ok {
		return fmt.Errorf("Controller with method '%s' and uri '%s' already registered by module '%s'",
			module, method, uri)
	}
	ctrlIndex[key] = module
	ctrls = append(ctrls, Controller{method, uri, handler})
	return
}
