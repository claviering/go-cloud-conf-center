package main

import (
	"database/sql"
	"main/router"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type MyString string

const MyStringNull MyString = "\x00"

func main() {
	db, err := sql.Open("sqlite3", "/Users/weiye/clondconfconter.db")
	utils.CheckErr(err)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "I'm a cook server.")
	})
	pmDeptRouterGroup := r.Group("/pmDept")
	router.Org(pmDeptRouterGroup, db)

	userRouterGroup := r.Group("/user")
	router.UsersRouter(userRouterGroup, db)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "http://localhost:8080")
}
