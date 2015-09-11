package models

import (
	"testing"
)

// This is a mess (tm).
func TestClass(t *testing.T) {
	u, err := RegisterUser("class_tester", "pw", "matt")
	c, err := u.NewClass("My Class", "A class for learning JavaScript")
	if err != nil {
		t.Error("Could not create class:", err)
	}

	if c.Name != "My Class" || c.Description != "A class for learning JavaScript" {
		t.Error("Incorrect information in class")
	}

	err = c.SetImage("imgs/house.png")
	if err != nil {
		t.Error("Could not set image", err)
	}

	if c.Image != "imgs/house.png" {
		t.Error("Incorrect image url")
	}

	cT, err := GetClassByID(c.ID)
	if err != nil {
		t.Error("Could not get class by ID", err)
	}

	if cT.ID != c.ID {
		t.Error("Not the same ID")
	}

	if cT.Name != c.Name {
		t.Error("Not the same name")
	}

	if cT.Description != c.Description {
		t.Error("Not the same description")
	}

	if cT.Image != c.Image {
		t.Error("not the same image")
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error("Could not get user")
	}

	if s.User != u.ID {
		t.Error("Wrong user ID")
	}

	if s.Class != c.ID {
		t.Error("Wrong class ID")
	}

	if s.Permissions != AdminPermissions {
		t.Error("Wrong permissions level")
	}

	s, err = c.Invite(u)
	if err == nil {
		t.Error("Was able to add the same student twice to a class")
	}

	st, err := c.Students()
	if err != nil {
		t.Error("Could not get students in this class", err)
	}

	if len(st) != 1 {
		t.Error("Wrong amount of students")
	}
}
