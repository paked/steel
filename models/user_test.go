package models

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	createTestDB()
	c := m.Run()
	deleteTestDB()

	os.Exit(c)
}
func TestRegister(t *testing.T) {
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
	paked, _ := GetUser("username", "paked")
	newbie, _ := GetUser("username", "newbie")

	if ok, err := paked.Login("pw"); err != nil || !ok {
		t.Error("Could not log in user:", err)
	}

	if ok, err := paked.Login("notpw"); err != nil || ok {
		t.Error("Could log in user with wrong password")
	}

	if ok, err := newbie.Login("pw"); err != nil || !ok {
		t.Error("Could not log in user:", err)
	}

	if ok, err := newbie.Login("notpw"); err != nil || ok {
		t.Error("Could log in user with wrong password")
	}
}

func TestDelete(t *testing.T) {
	paked, _ := GetUser("username", "paked")
	newbie, _ := GetUser("username", "newbie")

	if err := paked.Delete(); err != nil {
		t.Error("Could not delete paked")
	}

	if err := newbie.Delete(); err != nil {
		t.Error("Could not delete newbie")
	}
}
