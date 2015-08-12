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

	s, err := u.StartAssignment(a)
	if err != nil {
		t.Error("Could not start assignment: ", err)
	}

	if s.TeamName != u.Username+"'s Assignment" {
		t.Error("Submission values were not correct")
	}

	err = s.Rename("My Submission")
	if err != nil {
		t.Error("Could not rename submission")
	}

	if s.TeamName != "My Submission" {
		t.Error("Wrong submission name")
	}

	if sm, err := s.Members(); err != nil || len(sm) != 1 {
		t.Error("Failed wrong amount of members (0)", len(sm))
	}

	err = s.AddMember(u.ID)
	if err == nil {
		t.Error("Should have failed adding a user again")
	}

	u2, err := RegisterUser("member_add_test", "pw", "mat")
	if err != nil {
		panic(err)
	}

	err = s.AddMember(u2.ID)
	if err != nil {
		t.Error(err)
	}

	if sm, err := s.Members(); err != nil || len(sm) != 2 {
		t.Error("Failed wrong amount of members (1)")
	}

	a.Delete()
	u.Delete()
	u2.Delete()
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

func TestAllSubmissions(t *testing.T) {
	u, err := RegisterUser("all_submissions_test", "go", "golang.com")
	if err != nil {
		t.Error("User could not be registered", err)
		t.Fail()
	}

	a, err := u.CreateAssignment("A", "desc", "expl")
	if err != nil {
		t.Error("Could not make assignment")
	}

	sm := []Submission{}

	for i := 0; i < 10; i++ {
		s, err := u.StartAssignment(a)
		if err != nil {
			t.Error("Could not start assignment")
		}

		sm = append(sm, s)
	}

	if len(sm) != 10 {
		t.Error("Wrong amount of submissions")
	}

	a.Delete()
	u.Delete()
}
