package router

import (
	"database/sql"
	"main/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// /pmUser
func UsersRouter(r *gin.RouterGroup, db *sql.DB) {
	r.POST("/save", func(c *gin.Context) {
		type Request struct {
			DeptId   int    `json:"deptId" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Mobile   string `json:"mobile" binding:"required"`
			Nickname string `json:"nickname" binding:"required"`
			UserName string `json:"userName"`
		}
		var req Request
		c.BindJSON(&req)
		controller.SaveUser(db, req.DeptId, req.Email, req.Mobile, req.Nickname, req.UserName)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.GET("/list/:pageNum/:pageSize", func(c *gin.Context) {
		pageNum := c.Param("pageNum")
		pageSize := c.Param("pageSize")
		email := c.Query("email")
		deptId := c.Query("deptId")
		res, total := controller.List(db, email, deptId, pageNum, pageSize)
		c.JSON(200, gin.H{
			"data": gin.H{
				"data":  res,
				"total": total,
			},
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/resetPassword/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		controller.ResetPassword(db, userId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/delete/:id", func(c *gin.Context) {
		userId := c.Param("id")
		controller.ResetPassword(db, userId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/disableUser/:id", func(c *gin.Context) {
		userId := c.Param("id")
		controller.DisableUser(db, userId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/enableUser/:id", func(c *gin.Context) {
		userId := c.Param("id")
		controller.EnableUser(db, userId)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
	r.POST("/updateUser", func(c *gin.Context) {
		var req []controller.UpdateUserList
		c.BindJSON(&req)
		controller.UpdateUser(db, req)
		c.JSON(200, gin.H{
			"data":    "ok",
			"code":    200,
			"msg":     "success",
			"success": true,
		})
	})
}
