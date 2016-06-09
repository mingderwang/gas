package model

import (
	"fmt"
	"reflect"
	"strings"
	// "encoding/json"
)

var (
	FieldType = []string{
		"BIT",
		"TINYINT",
		"SMALLINT",
		"MEDIUMINT",
		"INT",
		"BIGINT",
		"DOUBLE",
		"FLOAT",
		"DECIMAL",
		"NUMBER",
		"DATE",
		"TIME",
		"TIMESTAMP",
		"DATETIME",
		"VARCHAR",
		"CHAR",
		"BLOB",
		"TEXT",
		"ENUM",
	}

	CustomType = []string{
		"STRING",
		"INCREMENTS",
	}

	FieldDefinition = []string{
		"NOT NULL",
		"NULL",
		"DEFAULT",
		"AUTO_INCREMENT",
		"UNIQUE",
		"PRIMARY",
		"COMMENT",
	}

	DefaultLength = map[string]string{
		"STRING":     "255",
		"VARCHAR":    "255",
		"INCREMENTS": "10",
		"INT":        "10",
	}
)

type FieldInfo struct {
	ftype  string              // field type in DB
	stype  reflect.StructField // field type in struct
	length string              // length

	refValue reflect.Value
	strValue string
}

func (fi *FieldInfo) GetType() string {
	return fi.ftype
}

func (fi *FieldInfo) SetType(t string) {
	fi.ftype = t
}

func (fi *FieldInfo) GetSType() reflect.StructField {
	return fi.stype
}

func (fi *FieldInfo) SetSType(t reflect.StructField) {
	fi.stype = t
}

func (fi *FieldInfo) String() string {
	return fi.strValue
}

func (fi *FieldInfo) GetReflectValue() reflect.Value {
	return fi.refValue
}

func (fi *FieldInfo) SetValue(rv reflect.Value) {
	fi.refValue = rv
	if rv.CanInterface() {
		fi.strValue = fmt.Sprintf("%v", rv.Interface())
	} else {
		fi.strValue = ""
	}

}

func (fi *FieldInfo) GetLength() string {
	if fi.length == "" {
		return fi.returnDefaultLength()
	}

	return fi.length
}

func (fi *FieldInfo) SetLength(l string) {
	fi.length = l
}

func (fi *FieldInfo) returnDefaultLength() string {
	v, ok := DefaultLength[strings.ToUpper(fi.GetType())]

	if ok {
		return v
	} else {
		return ""
	}
}
