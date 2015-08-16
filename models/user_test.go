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

func TestAdmin(t *testing.T) {
	u, _ := GetUser("username", "paked")
	if u.Permissions != DefaultPermissions {
		t.Error("Initial permissions are not default!")
	}

	err := u.MakeAdmin()
	if err != nil {
		t.Error("Error creating admin: ", err)
	}

	if u.Permissions != AdminPermissions {
		t.Error("Local permissions have not been changed")
	}

	// pull user from database
	u, _ = GetUser("username", "paked")
	if u.Permissions != AdminPermissions {
		t.Error("Wrong permissions in DB")
	}

	err = u.DemoteAdmin()
	if err != nil {
		t.Error("Could not demote admin...", err)
	}

	if u.Permissions != DefaultPermissions {
		t.Error("Local changes not made")
	}

	// pull user from database
	u, _ = GetUser("username", "paked")
	if u.Permissions != DefaultPermissions {
		t.Error("Changes not in DB")
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
