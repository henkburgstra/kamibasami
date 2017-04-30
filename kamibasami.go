package main

import (
	"fmt"

	"github.com/henkburgstra/kamibasami/controllers"
	"github.com/henkburgstra/kamibasami/node"
	"github.com/henkburgstra/kamibasami/service"

	"database/sql"

	"github.com/blevesearch/bleve"
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

	mapping := bleve.NewIndexMapping()
	index, err := bleve.Open("kamibasami.bleve")
	if err != nil {
		index, err = bleve.New("kamibasami.bleve", mapping)
	}
	if err != nil {
		fmt.Println("Error opening Bleve index")
		fmt.Println(err)
		return
	}
	svc := service.NewService()
	svc.SetIndex(index)
	svc.SetNodeRepo(node.NewDBNodeRepo(db, "sqlite"))

	router := gin.Default()
	router.Static("/wwwroot", "./wwwroot")

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
