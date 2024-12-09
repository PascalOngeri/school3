package handlers

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
}
