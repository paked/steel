package models

import (
	"database/sql"
	"testing"
)

func TestAssignments(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", "../db.db")
	if err != nil {
		panic(err)
	}

	u, err := RegisterUser("assignment_test", "pw", "mat")
	if err != nil {
		panic(err)
	}

	_, err = u.CreateAssignment("Test", "descritpion", "explanation")
	if err != nil {
		t.Error("Could not create Assignment: ", err)
	}

	u.Delete()
}
