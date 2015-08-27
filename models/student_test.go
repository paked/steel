package models

import "testing"

func TestStudentCreation(t *testing.T) {
	u, err := RegisterUser("student_creation_user", "x", "y")
	if err != nil {
		t.Error(err)
	}

	c, err := u.NewClass("Student", "Creation")
	if err != nil {
		t.Error(err)
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error("Could not retrieve user", err)
	}

	if s.User != u.ID {
		t.Error("User ID not correct")
	}

	if s.Class != c.ID {
		t.Error("ClassID not correct")
	}

	if s.Permissions != DefaultPermissions {
		t.Error("Normal permissions are broken")
	}

	sR, err := GetStudentByID(s.ID)
	if err != nil {
		t.Error("Could not get student")
	}

	if sR.User != s.User || sR.Class != s.Class || sR.Permissions != s.Permissions {
		t.Error("Values not sunc in DB")
	}

}

func TestStudentPermissions(t *testing.T) {
	u, err := RegisterUser("student_permission", "x", "y")
	if err != nil {
		t.Error(err)
	}

	c, err := u.NewClass("Permission", "Creation")
	if err != nil {
		t.Error(err)
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error("Could not retrieve user", err)
	}

	if s.Permissions != DefaultPermissions {
		t.Error("Incorrect permissions")
	}

	err = s.MakeAdmin()
	if err != nil {
		t.Error(err)
	}

	if s.Permissions != AdminPermissions {
		t.Error("pemrissinos ")
	}

	sR, err := GetStudentByID(s.ID)
	if err != nil {
		t.Error("Could not get student")
	}

	if s.Permissions != sR.Permissions {
		t.Error("Incorrect remote permissions")
	}

}
