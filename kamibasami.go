package main

import (
	"fmt"

	"github.com/henkburgstra/kamibasami/controllers"
	"github.com/henkburgstra/kamibasami/node"
	"github.com/henkburgstra/kamibasami/service"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "kamibasami.db")
	if err != nil {
		fmt.Println("Fout bij het openen van strongfit.db")
		fmt.Println(err)
		return
	}
	node.CreateTables(db)
	svc := service.NewService()
	svc.SetNodeRepo(node.NewDBNodeRepo(db, "sqlite"))

	router := gin.Default()
	router.GET("/test", testNode)
	for _, controller := range controllers.Get() {
		f := router.GET
		switch controller.Method {
		case "HEAD":
			f = router.HEAD
		case "POST":
			f = router.POST
		case "PUT":
			f = router.PUT
		case "DELETE":
			f = router.DELETE
		}
		f(controller.URI, func(c *gin.Context) {
			controller.Handler(svc, c)
		})
	}
	router.Run()
}
