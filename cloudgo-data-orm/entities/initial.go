package entities

import (
	_ "github.com/go-sql-driver/mysql" //for real
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

//var mydb *sql.DB
var engine *xorm.Engine

func init() {
	en, err := xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/test2?charset=utf8&parseTime=true")
	checkErr(err)
	engine = en
	engine.SetMapper(core.SameMapper{})
	u := &UserInfo{}
	exist, err2 := engine.IsTableExist(u)
	checkErr(err2)
	if !exist {
		err3 := engine.CreateTables(u)
		checkErr(err3)
	}
}

/*
// SQLExecer interface for supporting sql.DB and sql.Tx to do sql statement
type SQLExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// DaoSource Data Access Object Source
type DaoSource struct {
	// if DB, each statement execute sql with random conn.
	// if Tx, all statements use the same conn as the Tx's connection
	SQLExecer
}
*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
