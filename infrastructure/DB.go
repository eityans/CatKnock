package infrastructure

import (
	"database/sql"
	"log"
	"os"

	"catknock/model"

	"github.com/go-gorp/gorp"
)

func GetDb() *gorp.DbMap {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbMap.AddTableWithName(model.User{}, "users")

	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatal(err)
	}
	return dbMap
}
