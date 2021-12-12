package controller

import (
	"database/sql"
	"main/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Organization struct {
	Id         int    `json:"id"`
	Level      int    `json:"level"`
	Sort       int    `json:"sort"`
	OrgCode    string `json:"orgCode"`
	OrgName    string `json:"orgName"`
	ParentCode string `json:"parentCode"`
}

type Response struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []Organization `json:"data"`
}

func GetOrganizations(db *sql.DB) Response {
	// 查询数据
	rows, err := db.Query("SELECT id,level,sort,orgCode,orgName,parentCode FROM organization WHERE deleted IS NULL ORDER BY sort ASC")
	utils.CheckErr(err)
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
		utils.CheckErr(err)
		org := Organization{Id: id, Level: level, Sort: sort, OrgCode: orgCode, OrgName: orgName, ParentCode: parentCode.String}
		orgs = append(orgs, org)
	}
	data := Response{Status: "success", Message: "", Data: orgs}
	return data
}

func getMaxSort(db *sql.DB, parentCode string) int {
	// 查询数据
	rows, err := db.Query("SELECT sort FROM organization WHERE parentCode = ? ORDER BY sort DESC LIMIT 1", parentCode)
	utils.CheckErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		utils.CheckErr(err)
	}
	return sort
}

func AddOrganization(db *sql.DB, name string, parentCode string, level int) sql.Result {
	sort := getMaxSort(db, parentCode) + 1
	// 插入数据
	stmt, err := db.Prepare("INSERT INTO organization(level, orgCode, orgName, parentCode, sort) values(?,?,?,?,?)")
	utils.CheckErr(err)
	res, err := stmt.Exec(level, name, name, parentCode, sort)
	utils.CheckErr(err)
	return res
}

func getSortByid(db *sql.DB, id int) int {
	// 查询数据
	rows, err := db.Query("SELECT sort FROM organization WHERE id = ?", id)
	utils.CheckErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		utils.CheckErr(err)
	}
	return sort
}

func UpdateOrganization(db *sql.DB, id int, orgName string, sort int, parentCode string) sql.Result {
	oldSort := getSortByid(db, id)
	if oldSort > sort {
		// 调小排序
		_, _ = db.Exec("UPDATE organization SET sort = sort + 1 WHERE parentCode = ? AND sort >= ? AND sort < ?", parentCode, sort, oldSort)
	} else if oldSort < sort {
		// 调大排序
		_, _ = db.Exec("UPDATE organization SET sort = sort - 1 WHERE parentCode = ? AND sort > ? AND sort <= ?", parentCode, oldSort, sort)
	}
	res, err := db.Exec("UPDATE organization SET orgName = ?, sort = ? WHERE id = ?", orgName, sort, id)
	utils.CheckErr(err)
	return res
}
