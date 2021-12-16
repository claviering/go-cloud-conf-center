package controller

import (
	"database/sql"
	"fmt"
	"main/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Org struct {
	Id             int    `json:"id"`
	ORDER_NO       int    `json:"orderNo"`
	Level          int    `json:"level"`
	DEPT_PARENT_ID int32  `json:"deptParentId"`
	DEPT_CODE      string `json:"deptCode"`
	DEPT_NAME      string `json:"deptName"`
}

type Organization struct {
	Org
	CHILDREN []*Organization `json:"children"`
}

type OrgValue struct {
	Org
	CHILDREN []OrgValue `json:"children"`
}

type Response struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []Organization `json:"data"`
}

func getOrgValue(list []*Organization) []OrgValue {
	var res []OrgValue
	for _, item := range list {
		var children []OrgValue
		if len(item.CHILDREN) > 0 {
			children = getOrgValue(item.CHILDREN)
		}
		tmp := OrgValue{Org: Org{Id: item.Id, Level: item.Level, ORDER_NO: item.ORDER_NO, DEPT_CODE: item.DEPT_CODE, DEPT_NAME: item.DEPT_NAME, DEPT_PARENT_ID: item.DEPT_PARENT_ID}, CHILDREN: children}
		res = append(res, tmp)
	}
	return res
}

func buildTree(list []Organization) []OrgValue {
	if len(list) == 0 {
		return []OrgValue{}
	}
	treeMap := make(map[int][]*Organization)
	parentMap := make(map[int32][]*Organization)
	for _, item := range list {
		i := item
		_, ok := treeMap[item.Level]
		if !ok {
			treeMap[item.Level] = make([]*Organization, 0)
		}
		treeMap[item.Level] = append(treeMap[item.Level], &i)
		_, ok = parentMap[item.DEPT_PARENT_ID]
		if item.DEPT_PARENT_ID != 0 {
			if !ok {
				parentMap[item.DEPT_PARENT_ID] = make([]*Organization, 0)
			}
			parentMap[item.DEPT_PARENT_ID] = append(parentMap[item.DEPT_PARENT_ID], &i)
		}
	}
	for i := len(treeMap) - 2; i >= 0; i-- {
		for j := 0; j < len(treeMap[i]); j++ {
			id := (*treeMap[i][j]).Id
			fmt.Printf("id: %d", id)
			treeMap[i][j].CHILDREN = parentMap[int32(id)]
		}
	}
	return getOrgValue(treeMap[0])
}

func GetOrganizations(db *sql.DB) []OrgValue {
	// 查询数据
	rows, err := db.Query("SELECT id,level,sort,orgCode,orgName,deptParentId FROM organization WHERE deleted IS NULL ORDER BY sort ASC")
	utils.CheckErr(err)
	var orgs []Organization
	for rows.Next() {
		var id int
		var level int
		var sort int
		var orgCode string
		var orgName string
		var deptParentId sql.NullInt32
		// 顺序要和查询回来的字段顺序一致
		err = rows.Scan(&id, &level, &sort, &orgCode, &orgName, &deptParentId)
		utils.CheckErr(err)
		org := Organization{Org: Org{Id: id, Level: level, DEPT_CODE: orgCode, DEPT_NAME: orgName, DEPT_PARENT_ID: deptParentId.Int32, ORDER_NO: sort}}
		orgs = append(orgs, org)
	}
	return buildTree(orgs)
}

func getMaxSort(db *sql.DB, deptParentId int) int {
	// 查询数据
	rows, err := db.Query("SELECT sort FROM organization WHERE deptParentId = ? ORDER BY sort DESC LIMIT 1", deptParentId)
	utils.CheckErr(err)
	var sort int
	for rows.Next() {
		err = rows.Scan(&sort)
		utils.CheckErr(err)
	}
	return sort
}

func AddOrganization(db *sql.DB, name string, deptParentId int) sql.Result {
	sort := getMaxSort(db, deptParentId) + 1
	level := 0
	if deptParentId != 0 {
		level = getLevelByid(db, deptParentId) + 1
	}
	// 插入数据
	stmt, err := db.Prepare("INSERT INTO organization(level, orgCode, orgName, deptParentId, sort) values(?,?,?,?,?)")
	utils.CheckErr(err)
	res, err := stmt.Exec(level, name, name, deptParentId, sort)
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
func getLevelByid(db *sql.DB, id int) int {
	// 查询数据
	rows, err := db.Query("SELECT level FROM organization WHERE deptParentId = ?", id)
	utils.CheckErr(err)
	var level int
	for rows.Next() {
		err = rows.Scan(&level)
		utils.CheckErr(err)
	}
	return level
}

func UpdateOrganization(db *sql.DB, id int, orgName string, sort int, deptParentId int) sql.Result {
	oldSort := getSortByid(db, id)
	if oldSort > sort {
		// 调小排序
		_, _ = db.Exec("UPDATE organization SET sort = sort + 1 WHERE deptParentId = ? AND sort >= ? AND sort < ?", deptParentId, sort, oldSort)
	} else if oldSort < sort {
		// 调大排序
		_, _ = db.Exec("UPDATE organization SET sort = sort - 1 WHERE deptParentId = ? AND sort > ? AND sort <= ?", deptParentId, oldSort, sort)
	}
	res, err := db.Exec("UPDATE organization SET orgName = ?, sort = ? WHERE id = ?", orgName, sort, id)
	utils.CheckErr(err)
	return res
}

func MoveDept(db *sql.DB, id int, deptParentId int) {
	// 查询数据
	rows, err := db.Query("SELECT level FROM organization WHERE deptParentId = ?", deptParentId)
	utils.CheckErr(err)
	var level int
	rows.Scan(&level)
	sort := getMaxSort(db, deptParentId) + 1
	db.Exec("UPDATE organization SET deptParentId = ?,level = ?, sort = ? WHERE id = ?", deptParentId, level+1, sort, id)
}
