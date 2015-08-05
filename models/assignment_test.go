package models

import (
	"database/sql"
	"testing"
	"time"
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

	a, err := u.CreateAssignment("Test", "descritpion", "explanation")
	if err != nil {
		t.Error("Could not create Assignment: ", err)
	}

	if a.Name != "Test" || a.Description != "descritpion" || a.Explanation != "explanation" {
		t.Error("Assignment values were not set properly")
	}

	s, err := u.StartAssignment(a.ID)
	if err != nil {
		t.Error("Could not start assignment: ", err)
	}

	if s.TeamName != u.Username+"'s Assignment" {
		t.Error("Submission values were not correct")
	}

	a.Delete()
	u.Delete()
}

func TestDueAssignments(t *testing.T) {
	u, err := RegisterUser("due_assignments_test", "go", "golang.com")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		a, _ := u.CreateAssignment("Test", "testing", "terster")
		defer a.Delete()
	}

	tm := time.Now()

	as, err := u.DueAssignments(tm)

	if len(as) != 10 {
		t.Errorf("Expecting 10 assignments after %v got %v", tm.UnixNano(), len(as))
	}

	u.Delete()
}
