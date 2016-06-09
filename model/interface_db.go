package model

import (
	"database/sql"
)

type SlimDbInterface interface {
	Connect(protocal, hostname, port, username, password, dbname, params string) error
	ConnectWithConfig(interface{}) error
	parseConfig(interface{}) (protocal, hostname, port, username, password, dbname, params, charset string)
	ConnectWithDefault(username, password, dbname string) error

	Prepare(string) (*sql.Stmt, error)
	Exec(*sql.Stmt, ...interface{}) (sql.Result, error)
	Query(*sql.Stmt, ...interface{}) ([]map[string]string, error)

	Begin() error
	Commit() error
	Rollback() error

	Close() error
}

type SlimDb struct {
	SlimDbInterface // implements SlimDbInterface

	Conn    *sql.DB // Connection object
	isUseTX bool    // use transaction
	Tx      *sql.Tx // transaction obj
}
