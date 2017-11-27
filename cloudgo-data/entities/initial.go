package entities

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //for real
)

var mydb *sql.DB

func init() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}
	mydb = db
}

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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
