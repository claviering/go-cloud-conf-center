package controller

import (
	"database/sql"
	"main/utils"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id         int    `json:"id"`
	Mobile     string `json:"mobile"`
	Nickname   string `json:"nickname"`
	CreateTime int64  `json:"createTime"`
	UserName   string `json:"userName"`
	Email      string `json:"email"`
}

func getUserMaxSort(db *sql.DB) int {
	// 查询数据
	rows, err := db.Query("SELECT order_no FROM pm_user_dept ORDER BY order_no DESC LIMIT 1")
	utils.CheckErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		utils.CheckErr(err)
	}
	return sort
}

func SaveUser(db *sql.DB, deptId int, email string, mobile string, nickname string, userName string) {
	// get timestamp in milliseconds
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	rows, err := db.Exec("INSERT INTO pm_user(email, mobile, nickname, create_time, user_name) values(?,?,?,?,?)", email, mobile, nickname, timeStamp, userName)
	utils.CheckErr(err)
	userId, _ := rows.LastInsertId()
	_, err = db.Exec("INSERT INTO pm_user_dept(dept_id, user_id, order_no) values(?,?,?)", deptId, userId, getUserMaxSort(db)+1)
	utils.CheckErr(err)
}

func List(db *sql.DB, email string, pageNum string, pageSize string) ([]User, int) {
	// convert pageNum and pageSize to int
	pageNumInt, _ := strconv.Atoi(pageNum)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	queryAll := "SELECT count(id) FROM pm_user"
	if email != "" {
		queryAll = queryAll + " WHERE email like '%" + email + "%'"
	}
	total, _ := db.Query(queryAll)
	var totalCount int
	for total.Next() {
		total.Scan(&totalCount)
	}
	queryList := "SELECT id,mobile,nickname,create_time,user_name,email FROM pm_user"
	if email != "" {
		queryList = queryList + " WHERE email like '%" + email + "%'"
	}
	queryList = queryList + " LIMIT " + strconv.Itoa(pageSizeInt) + " OFFSET " + strconv.Itoa((pageNumInt-1)*pageSizeInt)
	rows, err := db.Query(queryList)
	utils.CheckErr(err)
	var user []User
	for rows.Next() {
		var id int
		var mobile string
		var nickname string
		var createTime int64
		var userName string
		var email string
		err = rows.Scan(&id, &mobile, &nickname, &createTime, &userName, &email)
		utils.CheckErr(err)
		user = append(user, User{
			Id:         id,
			Mobile:     mobile,
			Nickname:   nickname,
			CreateTime: createTime,
			UserName:   userName,
			Email:      email,
		})
	}

	return user, totalCount
}
