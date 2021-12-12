package router

import (
	"database/sql"
	"main/controller"
	"main/utils"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func UsersRouter(r *gin.Engine, db *sql.DB) {
	r.POST("/user/query", func(c *gin.Context) {
		type Request struct {
			PAGE_SIZE int    `json:"pageSize" `
			USERNAME  string `json:"username"`
		}
		var req Request
		c.BindJSON(&req)
		res, err := db.Exec("SELECT * FROM org_user")
		utils.CheckErr(err)
		c.JSON(200, gin.H{
			"message": res,
		})
	})
	r.POST("/user/add", func(c *gin.Context) {
		type Request struct {
			USERID   string `json:"userid" binding:"required"`
			MOBILE   string `json:"mobile" binding:"required"`
			USERNAME string `json:"username" binding:"required"`
			ORG_CODE string `json:"orgCode" binding:"required"`
		}
		var req Request
		c.BindJSON(&req)
		sort := controller.GetUserMaxSort(db) + 1
		// 时间戳
		timeStamp := time.Now().Unix()
		res, err := db.Exec("INSERT INTO org_user(userid,username,orgCode,mobile,createtime,sort) VALUES (?,?,?,?,?,?)", req.USERID, req.USERNAME, req.ORG_CODE, req.MOBILE, timeStamp, sort)
		utils.CheckErr(err)
		c.JSON(200, gin.H{
			"message": res,
		})
	})
}
