package model

import (
	"strings"
	// "database/sql"
	"strconv"
	// "fmt"
	"errors"
)

type MySQLBuilder struct {
	BaseBuilder
}

func (mb *MySQLBuilder) init() {
	mb.action = ""
	mb.selectColumns = make([]string, 0)
	mb.distinctColumns = make([]string, 0)
	mb.fromTable = ""
	mb.whereStatement = ""
	mb.havingStatement = ""
	mb.groupbyColumns = make([]string, 0)
	mb.orderbyStatement = make([]string, 0)

	mb.limitCount = 0
	mb.limitOffset = 0

	mb.countCol = make([]string, 0)

	// params
	mb.whereParams = make([]interface{}, 0)
	mb.havingParams = make([]interface{}, 0)
}

// ==== common functions =====
func (mb *MySQLBuilder) genSQL() string {
	sql := ""
	switch mb.action {
	case "SELECT":
		sql = mb.genSelectSQL()
	}

	return sql
}

// ==== Select functions ====

func (mb *MySQLBuilder) Select(columns ...string) BuilderInterface {
	mb.selectColumns = append(mb.selectColumns, columns...)

	return mb
}

func (mb *MySQLBuilder) Count(col ...string) BuilderInterface {
	if len(col) == 0 {
		mb.countCol = append(mb.countCol, "*")
	} else {
		mb.countCol = col
	}

	return mb
}

func (mb *MySQLBuilder) getColumns() string {
	columns := ""

	// count
	countArr := make([]string, 0)
	if len(mb.countCol) != 0 {
		for _, ccol := range mb.countCol {
			countArr = append(countArr, "COUNT("+ccol+")")
		}
	}

	if len(mb.selectColumns) == 0 {
		// no select columns but have count
		// then just select count
		if len(countArr) != 0 {
			columns += strings.Join(countArr, ", ")
		} else {
			// else select all columns
			columns += "*"
		}

	} else {
		// select columns
		columns += strings.Join(mb.selectColumns, ", ")

		// count columns
		if len(countArr) != 0 {
			columns += ", " + strings.Join(countArr, ", ")
		}

	}

	return columns
}

func (mb *MySQLBuilder) Distinct(columns ...string) BuilderInterface {
	mb.distinctColumns = append(mb.distinctColumns, columns...)

	return mb
}

func (mb *MySQLBuilder) getDistinctColumns() string {
	if len(mb.distinctColumns) == 0 {
		return ""
	}

	return "DISTINCT " + strings.Join(mb.distinctColumns, ", ")
}

func (mb *MySQLBuilder) From(tableName string) BuilderInterface {
	mb.fromTable = tableName

	return mb
}

func (mb *MySQLBuilder) Where(statement string, values ...interface{}) BuilderInterface {
	mb.whereStatement = statement
	mb.whereParams = append(mb.whereParams, values...)

	return mb
}

func (mb *MySQLBuilder) AndWhere(statement string, values ...interface{}) BuilderInterface {
	if mb.whereStatement == "" {
		mb.whereStatement = statement
	} else {
		mb.whereStatement += " AND " + statement
	}

	mb.whereParams = append(mb.whereParams, values...)

	return mb
}

func (mb *MySQLBuilder) OrWhere(statement string, values ...interface{}) BuilderInterface {
	if mb.whereStatement == "" {
		mb.whereStatement = statement
	} else {
		mb.whereStatement = "(" + mb.whereStatement + ")"
		mb.whereStatement += " OR " + statement
	}

	mb.whereParams = append(mb.whereParams, values...)

	return mb
}

func (mb *MySQLBuilder) genWhere() string {
	if mb.whereStatement != "" {
		return " WHERE " + mb.whereStatement
	} else {
		return ""
	}

}

func (mb *MySQLBuilder) Having(statement string, values ...interface{}) BuilderInterface {
	mb.havingStatement = statement
	mb.havingParams = append(mb.havingParams, values...)

	return mb
}

func (mb *MySQLBuilder) genHaving() string {
	if mb.havingStatement != "" {
		return " HAVING " + mb.havingStatement
	} else {
		return ""
	}

}

func (mb *MySQLBuilder) GroupBy(column ...string) BuilderInterface {
	mb.groupbyColumns = append(mb.groupbyColumns, column...)

	return mb
}

func (mb *MySQLBuilder) getGroupBy() string {
	if len(mb.groupbyColumns) != 0 {
		return " GROUP BY " + strings.Join(mb.groupbyColumns, ", ")
	}

	return ""
}

func (mb *MySQLBuilder) OrderBy(orderstatement string) BuilderInterface {
	mb.orderbyStatement = append(mb.orderbyStatement, orderstatement)

	return mb
}

func (mb *MySQLBuilder) Desc(column ...string) BuilderInterface {
	o := strings.Join(column, ",") + " DESC"

	mb.orderbyStatement = append(mb.orderbyStatement, o)

	return mb
}

func (mb *MySQLBuilder) Asc(column ...string) BuilderInterface {
	o := strings.Join(column, ",") + " ASC"

	mb.orderbyStatement = append(mb.orderbyStatement, o)

	return mb
}

func (mb *MySQLBuilder) genOrderby() string {
	if len(mb.orderbyStatement) != 0 {
		return " ORDER BY " + strings.Join(mb.orderbyStatement, ", ")
	}

	return ""
}

func (mb *MySQLBuilder) Limit(count int, offset ...int) BuilderInterface {
	mb.limitCount = count
	if len(offset) != 0 {
		mb.limitOffset = offset[0]
	}

	return mb
}

func (mb *MySQLBuilder) genLimit() string {
	s := ""

	if mb.limitCount != 0 {
		s += " LIMIT "

		if mb.limitOffset != 0 {
			s += strconv.Itoa(mb.limitOffset) + ", " + strconv.Itoa(mb.limitCount)
		} else {
			s += strconv.Itoa(mb.limitCount)
		}
	}

	return s
}

func (mb *MySQLBuilder) Get(tableStruct interface{}) ([]map[string]string, error) {
	table := mb.parseStruct(tableStruct)
	// set from table if fromTable is empty
	if mb.fromTable == "" {
		mb.fromTable = table.TableName
	}

	// default action is SELECT
	if mb.action == "" {
		mb.action = "SELECT"
	}

	sqlstr := mb.genSQL()
	// println(sqlstr)

	// param
	stmt, err := mb.GetDB().Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var p []interface{}
	if len(mb.havingParams) != 0 {
		if len(mb.whereParams) != 0 {
			p = append(mb.whereParams, mb.havingParams)
		} else {
			p = mb.havingParams
		}
	} else {
		p = mb.whereParams
	}

	res, err := mb.GetDB().Query(stmt, p...)

	// reset builder
	mb.init()

	return res, err
}

func (mb *MySQLBuilder) genSelectSQL() string {
	sql := "SELECT "

	// distinct
	disCol := mb.getDistinctColumns()

	// columns
	if disCol != "" {
		sql += disCol
	} else {
		sql += mb.getColumns()
	}

	// from
	sql += " FROM " + mb.fromTable

	// where
	sql += mb.genWhere()

	// group by
	sql += mb.getGroupBy()

	// having
	sql += mb.genHaving()

	// order by
	sql += mb.genOrderby()

	// limit
	sql += mb.genLimit()

	mb.lastSQL = sql

	return sql
}

// ==== Insert functions ====
func (mb *MySQLBuilder) Insert(tableStructWithValue interface{}) (int64, error) {
	table := mb.parseStruct(tableStructWithValue)

	// if value is valid
	if !table.InsertValid() {
		// throw error
		return 0, errors.New("Value does not valid before insert " + table.TableName + " data")
	}

	// final sql: INSERT INTO table.TableName (fieldName1, fieldName2) VALUES (?, ?)
	insertSQL := "INSERT INTO " + table.TableName

	fieldName := make([]string, 0)
	values := make([]interface{}, 0)

	for fname, finfo := range table.Fields {
		if finfo.GetReflectValue().CanInterface() {
			fieldName = append(fieldName, fname)
			values = append(values, finfo.String())
		}
	}

	insertSQL += " ( " + strings.Join(fieldName, ", ") + " ) "
	qmark := make([]string, 0)
	for i := 0; i < len(fieldName); i++ {
		qmark = append(qmark, "?")
	}
	insertSQL += "VALUES ( " + strings.Join(qmark, ", ") + " )"

	stmt, err := mb.GetDB().Prepare(insertSQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := mb.GetDB().Exec(stmt, values...)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// reset builder
	mb.init()

	mb.lastSQL = insertSQL

	return lastID, nil
}

// for debug function
func (mb *MySQLBuilder) getLastSQL() string {
	return mb.lastSQL
}
