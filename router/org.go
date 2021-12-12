package router

import (
	"database/sql"
	"main/controller"
	"main/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Org(r *gin.Engine, db *sql.DB) {
	r.GET("/org/get", func(c *gin.Context) {
		res := controller.GetOrganizations(db)
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
		controller.AddOrganization(db, req.NAME, req.PARENT_CODE, req.LEVEL)
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
		controller.UpdateOrganization(db, req.ID, req.ORG_NAME, req.SORT, req.PARENT_CODE)
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
		utils.CheckErr(err)
		var level int
		rows.Scan(&level)
		res, err := db.Exec("UPDATE organization SET parentCode = ?,level = ? WHERE id = ?", req.PARENT_CODE, level+1, req.ID)
		utils.CheckErr(err)
		c.JSON(200, gin.H{
			"message": res,
		})
	})
}
