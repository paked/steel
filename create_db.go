package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paked/configure"
)

var (
	conf = configure.New()

	db = conf.String("db", "database.db", "the file where you want you database")
)

func main() {
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewJSONFromFile("config.json"))

	conf.Parse()

	database, err := sql.Open("sqlite3", *db)
	if err != nil {
		log.Println("Unable to open database file")
		log.Printf("\t%v", err)
	}

	create, err := ioutil.ReadFile("models/create_db.sql")
	if err != nil {
		log.Println("unable to open create_db.sql in models/create_db.sql")
		log.Printf("\t%v", err)
	}

	_, err = database.Exec(string(create))
	if err != nil {
		log.Println("unable to open create_db.sql in models/create_db.sql")
		log.Printf("\t%v", err)
	}

	log.Println("Finished creating your database!")
}
