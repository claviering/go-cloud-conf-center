package controller

import (
	"database/sql"
	"main/utils"

	_ "github.com/mattn/go-sqlite3"
)

func GetUserMaxSort(db *sql.DB) int {
	// 查询数据
	rows, err := db.Query("SELECT sort FROM oog_user ORDER BY sort DESC LIMIT 1")
	utils.CheckErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		utils.CheckErr(err)
	}
	return sort
}
