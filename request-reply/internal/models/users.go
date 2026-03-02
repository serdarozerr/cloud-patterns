package models

import (
	"github.com/go-pg/pg/v10"
)
type User struct{
	Name string `pg:"name"`
	Email string `pg:"email,unique"`
	Age uint `pg:"age"`
}

func InsertUser(pg *pg.DB, user *User)error{
	_,err:=pg.Model(user).Insert()
	return err
}

