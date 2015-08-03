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

	if _, err := RegisterUser("paked", "pw", "hat"); err != nil {
		t.Error("Hello")
	}

	if _, err := RegisterUser("newbie", "pw", "mat"); err != nil {
		t.Error("Could not register user", err)
	}

	u, err := GetUser("username", "paked")
	if err != nil {
		t.Error("Could not get user", err)
	}

	if u.Username != "paked" || u.PasswordHash != "pw" || u.Email != "hat" {
		t.Error("invalid data")
	}
}

func TestLogin(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", "../db.db")

	if err != nil {
		t.Log(err)
		t.Error("Could not connect to DB")
	}

	if _, err := LoginUser("paked", "pw"); err != nil {
		t.Error("Could not log in user:", err)
	}

	if _, err := LoginUser("paked", "notpw"); err == nil {
		t.Error("Could log in user with wrong password")
	}

	if _, err := LoginUser("thisuserdoesnotexist", "pass"); err == nil {
		t.Error("Could login user")
	}

}

func TestAdmin(t *testing.T) {
	u, _ := LoginUser("paked", "pw")

	if u.Permissions != DefaultPermissions {
		t.Error("Wrong default permissions")
	}

	if u2, _ := LoginUser("paked", "pw"); u2.Permissions != DefaultPermissions {
		t.Error("Wrong permissions level in DB")
	}

	err := u.MakeAdmin()

	if err != nil {
		t.Error("Admin creation error: ", err)
	}

	if u.Permissions != AdminPermissions {
		t.Error("Wrong permission level")
	}

	if u2, _ := LoginUser("paked", "pw"); u2.Permissions != AdminPermissions {
		t.Error("Wrong permissions level in DB")
	}
}

func TestDelete(t *testing.T) {
	paked, _ := LoginUser("paked", "pw")
	newbie, _ := LoginUser("newbie", "pw")

	if err := paked.Delete(); err != nil {
		t.Error("Could not delete paked")
	}

	if err := newbie.Delete(); err != nil {
		t.Error("Could not delete newbie")
	}
}
