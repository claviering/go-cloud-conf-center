package router

import (
	"database/sql"
	"main/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Org(r *gin.RouterGroup, db *sql.DB) {
	r.GET("/orgTree", func(c *gin.Context) {
		res := controller.GetOrganizations(db)
		c.JSON(200, gin.H{
			"data":    res,
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/save", func(c *gin.Context) {
		type Request struct {
			DeptName     string `json:"deptName"`
			ParentDeptId int    `json:"parentDeptId"`
		}
		var req Request
		c.Bind(&req)
		controller.AddOrganization(db, req.DeptName, req.ParentDeptId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/deleteDept", func(c *gin.Context) {
		type Request struct {
			Id int `json:"id" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		db.Exec("UPDATE organization SET deleted = ? WHERE id = ?", 1, req.Id)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/updateDept", func(c *gin.Context) {
		type Request struct {
			Id            int    `json:"id"`
			DeptName      string `json:"deptName"`
			OrderNo       int    `json:"orderNo"`
			GroupParentId int    `json:"groupParentId"`
		}
		var req Request
		c.Bind(&req)
		controller.UpdateOrganization(db, req.Id, req.DeptName, req.OrderNo, req.GroupParentId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/updateList", func(c *gin.Context) {
		type Request struct {
			Id      int `json:"id" binding:"required"`
			OrderNo int `json:"orderNo" binding:"required"`
		}
		var req []Request
		c.BindJSON(&req)
		for _, v := range req {
			db.Exec("UPDATE organization SET sort = ? WHERE id = ?", v.OrderNo, v.Id)
		}
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/moveDept", func(c *gin.Context) {
		type Request struct {
			ID            int `json:"id" binding:"required"`
			GroupParentId int `json:"groupParentId" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		controller.MoveDept(db, req.ID, req.GroupParentId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
}
