package controller

import (
	"database/sql"
	"main/utils"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id            int    `json:"id"`
	Mobile        string `json:"mobile"`
	Nickname      string `json:"nickname"`
	CreateTime    int64  `json:"createTime"`
	UserName      string `json:"userName"`
	Email         string `json:"email"`
	GroupName     string `json:"groupName"`
	MainGroupName string `json:"mainGroupName"`
	DeptId        int    `json:"deptId"`
	OrderNo       int    `json:"orderNo"`
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

func List(db *sql.DB, email string, deptId string, pageNum string, pageSize string) ([]User, int) {
	// convert pageNum and pageSize to int
	pageNumInt, _ := strconv.Atoi(pageNum)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	queryAll := "SELECT count(pm_user.id) FROM pm_user"
	queryList := "SELECT pm_user.id,mobile,nickname,pm_user.create_time,user_name,email FROM pm_user"
	if deptId != "" {
		queryAll = queryAll + " JOIN pm_user_dept ON pm_user.id = pm_user_dept.user_id WHERE pm_user_dept.dept_id = " + deptId
		queryList = queryList + " JOIN pm_user_dept ON pm_user.id = pm_user_dept.user_id WHERE pm_user_dept.dept_id = " + deptId
	}
	if email != "" {
		if deptId != "" {
			queryAll = queryAll + " AND"
			queryList = queryList + " AND"
		} else {
			queryAll = queryAll + " WHERE"
			queryList = queryList + " WHERE"
		}
		queryAll = queryAll + " email like '%" + email + "%'"
		queryList = queryList + " email like '%" + email + "%'"
	}
	queryList = queryList + " LIMIT " + strconv.Itoa(pageSizeInt) + " OFFSET " + strconv.Itoa((pageNumInt-1)*pageSizeInt)

	total, _ := db.Query(queryAll)
	var totalCount int
	for total.Next() {
		total.Scan(&totalCount)
	}
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
			Id:            id,
			Mobile:        mobile,
			Nickname:      nickname,
			CreateTime:    createTime,
			UserName:      userName,
			Email:         email,
			GroupName:     "组织名称",
			MainGroupName: "主组织名称",
			DeptId:        5,
			OrderNo:       8,
		})
	}

	return user, totalCount
}

func ResetPassword(db *sql.DB, userId string) {

}

func DisableUser(db *sql.DB, userId string) {
	_, err := db.Exec("UPDATE pm_user SET user_status = 0 WHERE id = ?", userId)
	utils.CheckErr(err)
}
func EnableUser(db *sql.DB, userId string) {
	_, err := db.Exec("UPDATE pm_user SET user_status = 1 WHERE id = ?", userId)
	utils.CheckErr(err)
}

type UpdateUserList struct {
	Id     int `json:"id"`
	DeptId int `json:"deptId"`
}

func UpdateUser(db *sql.DB, list []UpdateUserList) {
	// update dept_id in pm_user where id in UpdateUserList
	var userList []string
	for _, v := range list {
		// convert id to string
		userList = append(userList, strconv.Itoa(v.Id))
	}
	users := "(" + strings.Join(userList, ",") + ")"
	_, err := db.Exec("UPDATE pm_user_dept SET dept_id = ? WHERE user_id in "+users, list[0].DeptId)
	utils.CheckErr(err)
}
