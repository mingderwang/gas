package model

import (
// "database/sql"
// "reflect"
// "fmt"
)

var ()

type ModelInterface interface {
	// Connect (protocol, hostname, port, username, password, dbname, params string) error
	// ConnectWithConfig (map[string]string) error
	// Close ()

	SetDB(SlimDbInterface)
	GetDB() SlimDbInterface

	// orm functions
	Insert(interface{}) (int64, error) // Save struct data
	MultiInsert(...interface{}) ([]int64, error)

	// builder
	SetBuilder(BuilderInterface)
	Builder() BuilderInterface

	// TestConn () // Test func

	// private functions
	parseTable(interface{}) TableInfo // parse table from struct
	parseField() FieldInfo

	// transaction
	TransactionStart() error
	TransactionCommit() error
	TransactionRollback() error
}

type BuilderWraperInterface interface {
	Select(columns ...string) ModelInterface
	From(tableName string) ModelInterface
	Where(columnName, operation, value string) ModelInterface
	Get() map[string]string
}

type Model struct {
	ModelInterface // implements ModelInterface

	db SlimDbInterface
	b  BuilderInterface // model has builder

}

func (m *Model) SetDB(db SlimDbInterface) {
	m.db = db
}

func (m *Model) GetDB() SlimDbInterface {
	return m.db
}

func (m *Model) SetBuilder(b BuilderInterface) {
	m.b = b
}

func (m *Model) Builder() BuilderInterface {
	return m.b
}
