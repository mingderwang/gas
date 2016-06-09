package models

import (
// "github.com/gowebtw/goslim/model"
)

type TestUser struct {
	ID   int    `type:"Interger",prop:"autoincrement"`
	Name string `type:"Varchar",length:"10"`
}
