package models

import (
	"testing"
)

func TestClass(t *testing.T) {
	u, err := RegisterUser("class_tester", "pw", "matt")
	c, err := NewClass("My Class", "A class for learning JavaScript")
	if err != nil {
		t.Error("Could not create class:", err)
	}

	if c.Name != "My Class" || c.Description != "A class for learning JavaScript" {
		t.Error("Incorrect information in class")
	}

	err = c.AddUser(u)
	if err != nil {
		t.Error("Could not add user")
	}

	err = c.AddUser(u)
	if err == nil {
		t.Error("Was able to add the same student twice to a class")
	}

	st, err := c.Students()
	if err != nil {
		t.Error("Could not get students in this class")
	}

	if len(st) != 1 {
		t.Error("Wrong amount of students")
	}
}
