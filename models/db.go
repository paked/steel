package models

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func createTestDB() {
	var err error

	db, err = sql.Open("sqlite3", "tst.db")
	if err != nil {
		panic(err)
	}

	create, err := ioutil.ReadFile("create_db.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(create))
	if err != nil {
		panic(err)
	}
}

func deleteTestDB() {
	os.Remove("tst.db")
}

func InitDB(file string) {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}
}
