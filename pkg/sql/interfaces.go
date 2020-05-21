package sql

import "database/sql"

type Querier interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Executer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type ExecQueryer interface {
	Executer
	Querier
}
