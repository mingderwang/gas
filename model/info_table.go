package model

import (
	"errors"
	"reflect"
)

type TableInfo struct {
	TableName string
	Fields    map[string]*FieldInfo // map[columnName]FieldInfo
	SeqField  *FieldInfo
	FieldNum  int
}

func (ti *TableInfo) InsertValid() bool {

	return true
}

func (ti *TableInfo) SetSeqField(f *FieldInfo) {
	ti.SeqField = f
}

func (ti *TableInfo) GetSeqField() *FieldInfo {
	return ti.SeqField
}

func (ti *TableInfo) SetValue(col string, val interface{}) error {
	f, ok := ti.Fields[col]

	if ok {
		f.SetValue(reflect.ValueOf(val))

		return nil
	}

	return errors.New(col + " column not exist in " + ti.TableName + " table")
}

func (ti *TableInfo) GetValue(col string) string {
	f, ok := ti.Fields[col]

	if ok {
		return f.String()
	}

	return ""
}
