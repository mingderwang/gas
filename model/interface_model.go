package model

import "github.com/go-gas/SQLBuilder"

type ModelInterface interface {
	// Connect (protocol, hostname, port, username, password, dbname, params string) error
	// ConnectWithConfig (map[string]string) error
	// Close ()

	//SetDB(SlimDbInterface)
	//GetDB() SlimDbInterface

	// orm functions
	Insert(interface{}) (int64, error) // Save struct data
	MultiInsert(...interface{}) ([]int64, error)

	// builder
	SetBuilder(SQLBuilder.BuilderInterface)
	Builder() SQLBuilder.BuilderInterface

	// TestConn () // Test func

	// private functions
	//parseTable(interface{}) TableInfo // parse table from struct
	//parseField() FieldInfo

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

	b  SQLBuilder.BuilderInterface // model has builder
}

func (m *Model) SetBuilder(b SQLBuilder.BuilderInterface) {
	m.b = b
}

func (m *Model) Builder() SQLBuilder.BuilderInterface {
	return m.b
}
