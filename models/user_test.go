package models

import (
	"database/sql"
	"testing"
)

func TestRegister(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", "../db.db")

	if err != nil {
		t.Log(err)
		t.Error("Could not connect to db")
	}

	if _, err := RegisterUser("paked", "pw", "hat"); err == nil {
		t.Error("Hello")
	}

	if _, err := RegisterUser("newbie", "pw", "mat"); err != nil {
		t.Error("Could not register user", err)
	}
}
