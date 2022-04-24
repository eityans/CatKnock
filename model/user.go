package model

type User struct {
	Id int `db:"id"`
	Name string `db:"name"`
	Age int `db:"age"`
}
