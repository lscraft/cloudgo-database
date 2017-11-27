package sqlt

import (
	"database/sql"
)

// SQLExecer .
type SQLExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// RowScanner .
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// RowScannerWithColumnInfo .
type RowScannerWithColumnInfo interface {
	RowScanner
	Columns() ([]string, error)
}

// RowMapperCallback .
type RowMapperCallback func(rscan RowScanner) error

// SQLTemplate for CRUD
type SQLTemplate struct {
	// SQLExecer could be sql.DB or sql.TX
	// if DB, each statement execute sql with random conn.
	// if Tx, all statements use the same conn as the Tx's connection
	SQLExecer
}

// NewSQLTemplate .
func NewSQLTemplate(sqlExecer SQLExecer) SQLTemplate {
	tpl := SQLTemplate{sqlExecer}
	return tpl
}

//Create Operation for Template

//Insert .
func (sqlt *SQLTemplate) Insert(insertQuery string, id *int, args ...interface{}) error {
	res, err := sqlt.Exec(insertQuery, args...)
	if err != nil {
		return err
	}

	if id != nil {
		newID, _ := res.LastInsertId()
		*id = int(newID)
	}

	return nil
}

//Select .
func (sqlt *SQLTemplate) Select(selectQuery string, rowMapper RowMapperCallback, args ...interface{}) error {

	// https://stackoverflow.com/questions/24878264/how-can-i-build-varidics-of-type-interface-in-go
	rows, err := sqlt.Query(selectQuery, args...)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rowMapper(rows)
		if err != nil {
			return err
		}
	}

	return nil
}

//SelectOne .
func (sqlt *SQLTemplate) SelectOne(selectQuery string, rowMapper RowMapperCallback, args ...interface{}) error {
	row := sqlt.QueryRow(selectQuery, args...)
	err := rowMapper(row)
	if err != nil {
		return err
	}
	return nil
}

//Update .
func (sqlt *SQLTemplate) Update(updateQuery string, args ...interface{}) (int, error) {
	res, err := sqlt.Exec(updateQuery, args...)
	if err != nil {
		return 0, err
	}

	af, err := res.RowsAffected()
	return int(af), err
}

//Delete .
func (sqlt *SQLTemplate) Delete(deleteQuery string, args ...interface{}) (int, error) {
	res, err := sqlt.Exec(deleteQuery, args)
	if err != nil {
		return 0, err
	}

	af, err := res.RowsAffected()
	return int(af), err
}
