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

func testNode(c *gin.Context) {
	repo := node.NewMockNodeRepo()
	node1 := node.NewNode("1", "parent", "")
	repo.Put(node1)
	node2 := node.NewNode("2", "child", "1")
	repo.Put(node2)
	node3 := node.NewNode("3", "grandchild", "2")
	repo.Put(node3)
	path := node.GetPath(repo, node3)
	c.String(200, path)
}

// func apiNevoData(c *gin.Context, db *sqlx.DB) {
// 	var respons NevoDataResonse
// 	respons.Foodgroups = nevo.AllFoodgroups(db)
// 	respons.Foods = nevo.AllFoods(db)
// 	respons.FoodgroupCount = len(respons.Foodgroups)
// 	respons.FoodCount = len(respons.Foods)
// 	c.JSON(200, respons)
// }

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
