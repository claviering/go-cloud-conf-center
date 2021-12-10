package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type MyString string

const MyStringNull MyString = "\x00"

type Response struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []Organization `json:"data"`
}

type Organization struct {
	Id         int    `json:"id"`
	Level      int    `json:"level"`
	Sort       int    `json:"sort"`
	OrgCode    string `json:"orgCode"`
	OrgName    string `json:"orgName"`
	ParentCode string `json:"parentCode"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getOrganizations(db *sql.DB) Response {
	// 查询数据
	rows, err := db.Query("SELECT id,level,sort,orgCode,orgName,parentCode FROM organization WHERE deleted IS NULL ORDER BY sort ASC")
	checkErr(err)
	var orgs []Organization
	for rows.Next() {
		var id int
		var level int
		var sort int
		var orgCode string
		var orgName string
		var parentCode sql.NullString
		// 顺序要和查询回来的字段顺序一致
		err = rows.Scan(&id, &level, &sort, &orgCode, &orgName, &parentCode)
		checkErr(err)
		org := Organization{Id: id, Level: level, Sort: sort, OrgCode: orgCode, OrgName: orgName, ParentCode: parentCode.String}
		orgs = append(orgs, org)
	}
	data := Response{Status: "success", Message: "", Data: orgs}
	return data
}

func getMaxSort(db *sql.DB, parentCode string) int {
	// 查询数据
	rows, err := db.Query("SELECT sort FROM organization WHERE parentCode = ? ORDER BY sort DESC LIMIT 1", parentCode)
	checkErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		checkErr(err)
	}
	return sort
}

func getSortByid(db *sql.DB, id int) int {
	// 查询数据
	rows, err := db.Query("SELECT sort FROM organization WHERE id = ?", id)
	checkErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		checkErr(err)
	}
	return sort
}

func addOrganization(db *sql.DB, name string, parentCode string, level int) {
	fmt.Printf("parentCode is %s\n", parentCode)
	sort := getMaxSort(db, parentCode) + 1
	// 插入数据
	stmt, err := db.Prepare("INSERT INTO organization(level, orgCode, orgName, parentCode, sort) values(?,?,?,?,?)")
	checkErr(err)
	res, err := stmt.Exec(level, name, name, parentCode, sort)
	checkErr(err)
	fmt.Print(res)
}

func updateOrganization(db *sql.DB, id int, orgName string, sort int, parentCode string) sql.Result {
	oldSort := getSortByid(db, id)
	fmt.Printf("oldSort is %d\n", oldSort)
	fmt.Printf("sort is %d\n", sort)
	if oldSort > sort {
		// 调小排序
		_, _ = db.Exec("UPDATE organization SET sort = sort + 1 WHERE parentCode = ? AND sort >= ? AND sort < ?", parentCode, sort, oldSort)
	} else if oldSort < sort {
		// 调大排序
		_, _ = db.Exec("UPDATE organization SET sort = sort - 1 WHERE parentCode = ? AND sort > ? AND sort <= ?", parentCode, oldSort, sort)
	}
	res, err := db.Exec("UPDATE organization SET orgName = ?, sort = ? WHERE id = ?", orgName, sort, id)
	checkErr(err)
	return res
}

func main() {
	db, err := sql.Open("sqlite3", "/Users/weiye/clondconfconter.db")
	checkErr(err)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "I'm a cook server.")
	})
	r.GET("/org/get", func(c *gin.Context) {
		res := getOrganizations(db)
		c.JSON(200, gin.H{
			"message": res,
		})
	})
	r.POST("/org/add", func(c *gin.Context) {
		type Request struct {
			NAME        string `json:"name" binding:"required"`
			PARENT_CODE string `json:"parentCode" binding:"required"`
			LEVEL       int    `json:"level" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		addOrganization(db, req.NAME, req.PARENT_CODE, req.LEVEL)
		c.JSON(200, gin.H{
			"message": req.NAME,
		})
	})
	r.POST("/org/del", func(c *gin.Context) {
		type Request struct {
			ID int `json:"id" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		db.Exec("UPDATE organization SET deleted = ? WHERE id = ?", 1, req.ID)
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
	r.POST("/org/update", func(c *gin.Context) {
		type Request struct {
			ID          int    `json:"id" binding:"required"`
			ORG_NAME    string `json:"orgName" binding:"required"`
			SORT        int    `json:"sort" binding:"required"`
			PARENT_CODE string `json:"parentCode" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		updateOrganization(db, req.ID, req.ORG_NAME, req.SORT, req.PARENT_CODE)
		c.JSON(200, gin.H{
			"message": req.ID,
		})
	})
	r.POST("/org/updateList", func(c *gin.Context) {
		type Request struct {
			ID   int `json:"id" binding:"required"`
			SORT int `json:"sort" binding:"required"`
		}
		var req []Request
		c.BindJSON(&req)
		for _, v := range req {
			db.Exec("UPDATE organization SET sort = ? WHERE id = ?", v.SORT, v.ID)
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
	r.POST("/org/move", func(c *gin.Context) {
		type Request struct {
			ID          int    `json:"id" binding:"required"`
			PARENT_CODE string `json:"parentCode" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		// 查询数据
		rows, err := db.Query("SELECT level FROM organization WHERE orgCode = ?", req.PARENT_CODE)
		checkErr(err)
		var level int
		rows.Scan(&level)
		res, err := db.Exec("UPDATE organization SET parentCode = ?,level = ? WHERE id = ?", req.PARENT_CODE, level+1, req.ID)
		checkErr(err)
		c.JSON(200, gin.H{
			"message": res,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "http://localhost:8080")
}
