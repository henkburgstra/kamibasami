package main

import (
	"strongfit/nevo"

	"github.com/henkburgstra/kamibasami/node"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type NevoDataResonse struct {
	FoodgroupCount int
	FoodCount      int
	Foodgroups     []*nevo.Foodgroup
	Foods          []*nevo.Food
}

func apiNevoImport(c *gin.Context) {
	nevo.NevoImport()
}

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

func apiNevoData(c *gin.Context, db *sqlx.DB) {
	var respons NevoDataResonse
	respons.Foodgroups = nevo.AllFoodgroups(db)
	respons.Foods = nevo.AllFoods(db)
	respons.FoodgroupCount = len(respons.Foodgroups)
	respons.FoodCount = len(respons.Foods)
	c.JSON(200, respons)
}

func main() {
	db, err := sqlx.Open("sqlite3", "strongfit.db")
	if err != nil {
		fmt.Println("Fout bij het openen van strongfit.db")
		fmt.Println(err)
		return
	}
	router := gin.Default()
	router.GET("/api/nevo/import", apiNevoImport)
	router.GET("/api/nevo/data", func(c *gin.Context) {
		apiNevoData(c, db)
	})
	router.GET("/test", testNode)
	router.Run()
}
