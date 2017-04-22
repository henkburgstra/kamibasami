package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var ctrls = make([]gin.HandlerFunc, 0)
var ctrlIndex = make(map[string]string)

// RegisterController is used to register a controller
func RegisterController(module string, method string, uri string, handler gin.HandlerFunc) (err error) {
	var key = fmt.Sprintf("%s:%s", method, uri)
	if module, ok := ctrlIndex[key]; ok {
		return fmt.Errorf("Controller with method '%s' and uri '%s' already registered by module '%s'",
			module, method, uri)
	}
	ctrlIndex[key] = module
	ctrls = append(ctrls, handler)
	return
}
