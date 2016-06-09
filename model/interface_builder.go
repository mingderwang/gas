package model

import (
	"reflect"
	"strings"
)

var (
// parsedTable map[string]*TableInfo = make(map[string]*TableInfo, 0) // store parsed table info parsedTable[tablename]TableInfo{}
)

type BuilderInterface interface {
	// db functions
	SetDB(SlimDbInterface)
	GetDB() SlimDbInterface

	// build sql functions
	Select(columns ...string) BuilderInterface
	Distinct(columns ...string) BuilderInterface

	From(tableName string) BuilderInterface

	Where(statement string, values ...interface{}) BuilderInterface
	AndWhere(statement string, values ...interface{}) BuilderInterface
	OrWhere(statement string, values ...interface{}) BuilderInterface

	GroupBy(column ...string) BuilderInterface

	Having(statement string, values ...interface{}) BuilderInterface

	OrderBy(orderstatement string) BuilderInterface
	Desc(column ...string) BuilderInterface
	Asc(column ...string) BuilderInterface

	Limit(count int, offset ...int) BuilderInterface
	// Offset()

	Count(col ...string) BuilderInterface

	// action functions
	Get(interface{}) ([]map[string]string, error)
	Insert(tableStructWithValue interface{}) (int64, error) // return last insert id, default is 0
	Update()
	Delete()

	// generate functions
	genWhere() string
	genSQL() string

	// for debug function
	getLastSQL() string
}

type BaseBuilder struct {
	BuilderInterface

	db SlimDbInterface

	action           string
	selectColumns    []string
	distinctColumns  []string
	fromTable        string
	whereStatement   string
	havingStatement  string
	groupbyColumns   []string
	orderbyStatement []string

	whereParams  []interface{}
	havingParams []interface{}

	limitOffset int
	limitCount  int

	countCol []string

	lastSQL string
}

func (bb *BaseBuilder) SetDB(db SlimDbInterface) {
	bb.db = db
}

func (bb *BaseBuilder) GetDB() SlimDbInterface {
	return bb.db
}

func (bb *BaseBuilder) parseStruct(s interface{}) *TableInfo {
	rv := reflect.Indirect(reflect.ValueOf(s))
	rt := reflect.TypeOf(rv.Interface())

	// fetch tablename
	tableName := strings.ToLower(rt.Name())

	// fetch tableinfo from global variable (cache)
	// tinfo, ok := parsedTable[tableName]
	// if !ok {

	// table info object
	tinfo := &TableInfo{}
	tinfo.TableName = tableName

	// fetch fields
	fs := make(map[string]*FieldInfo, 0)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fi := &FieldInfo{}
		// fetch field name
		fname := f.Name

		f.Type.Name()
		// fetch tags
		tagType := f.Tag.Get("type")
		tagLength := f.Tag.Get("length")

		fi.SetType(tagType)
		fi.SetLength(tagLength)
		fi.SetSType(f)

		// set value
		fi.SetValue(rv.Field(i))

		// check field is seq
		tagSeq := f.Tag.Get("isSeq")
		if tagType == "increments" || tagSeq == "true" {
			// is seq field
			tinfo.SetSeqField(fi)
		}

		fs[fname] = fi
	}

	tinfo.Fields = fs
	tinfo.FieldNum = len(fs)

	// add to global variable
	// parsedTable[tableName] = tinfo
	// }

	return tinfo
}
