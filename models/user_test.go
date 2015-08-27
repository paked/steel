package models

import (
	"fmt"
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

func TestGetClasses(t *testing.T) {
	u, err := RegisterUser("get_classes_test", "x", "zz")
	if err != nil {
		t.Error("Could not register user")
	}

	for i := 0; i < 10; i++ {
		c, err := u.NewClass(fmt.Sprintf("class %v", i), "x")
		if err != nil {
			t.Error("Could not create class")
		}

		_, err = c.Invite(u)
		if err != nil {
			t.Error("Could not invite to class")
		}
	}

	u2, err := RegisterUser("get_classes_test_second", "x", "zz")
	if err != nil {
		t.Error("Could not register second user")
	}

	for i := 0; i < 10; i++ {
		c, err := u.NewClass(fmt.Sprintf("other class %v", i), "x")
		if err != nil {
			t.Error("Could not create class")
		}

		_, err = c.Invite(u)
		if err != nil {
			t.Error("could not invite to class ")
		}

		_, err = c.Invite(u2)
		if err != nil {
			t.Error("Could not invite to class")
		}
	}

	cs, err := u2.Classes()
	if err != nil {
		t.Error("Could not get classes")
	}

	if len(cs) != 10 {
		t.Error("Wrong amount of classes wanted 10, got", len(cs))
	}

	cs, err = u.Classes()
	if err != nil {
		t.Error("Could not get classes", err)
	}

	if len(cs) != 20 {
		t.Error("Wrong amount of classes wanted 20, got ", len(cs))
	}
}
