package models

import (
	"fmt"
	"testing"
	"time"
)

func TestAssignments(t *testing.T) {
	u, err := RegisterUser("assignment_test", "pw", "mat")
	if err != nil {
		t.Error(err)
	}

	c, err := u.NewClass("assignments_test", "testing assignments")
	if err != nil {
		t.Error(err)
	}

	_, err = c.Invite(u)
	if err == nil {
		t.Error("Already in class")
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error("That student doesnt exist:", err)
	}

	err = s.MakeAdmin()
	if err != nil {
		t.Error("COuld not make admin", err)
	}

	a, err := s.CreateWorkshop("Test", "descritpion", "explanation")
	if err != nil {
		t.Error("Could not create Assignment: ", err)
	}

	if a.Name != "Test" || a.Description != "descritpion" || a.Explanation != "explanation" {
		t.Error("Assignment values were not set properly")
	}

	sub, err := s.StartAssignment(a)
	if err != nil {
		t.Error("Could not start assignment: ", err)
	}

	if sub.TeamName != "Assignment" {
		t.Error("Submission values were not correct")
	}

	err = sub.Rename("My Submission")
	if err != nil {
		t.Error("Could not rename submission")
	}

	if sub.TeamName != "My Submission" {
		t.Error("Wrong submission name")
	}

	if sm, err := sub.Members(); err != nil || len(sm) != 1 {
		t.Error("Failed wrong amount of members (0)", len(sm), err)
	}

	err = sub.Invite(s)
	if err == nil {
		t.Error("Should have failed adding a user again")
	}

	u2, err := RegisterUser("member_add_test", "pw", "mat")
	if err != nil {
		t.Error(err)
	}

	s2, err := c.Invite(u2)
	if err != nil {
		t.Error(err)
	}

	err = sub.Invite(s2)

	if sm, err := sub.Members(); err != nil || len(sm) != 2 {
		t.Error("Failed wrong amount of members (1)", len(sm), err)
	}

	a.Delete()
	u.Delete()
	u2.Delete()
}

func TestDueAssignments(t *testing.T) {
	u, err := RegisterUser("due_assignments_test", "go", "golang.com")
	if err != nil {
		t.Error(err)
	}

	c, err := u.NewClass("due assignments class", "things")
	if err != nil {
		t.Error(err)
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error(err)
	}

	s.MakeAdmin()

	for i := 0; i < 10; i++ {
		a, err := s.CreateWorkshop("Test", "testing", "terster")
		if err != nil {
			t.Error(err)
		}

		defer a.Delete()
	}

	tm := time.Now()

	as, err := s.Workshops(tm)

	if len(as) != 10 {
		t.Errorf("Expecting 10 assignments after %v got %v err: %v", tm.UnixNano(), len(as), err)
	}

	u.Delete()
}

func TestAllSubmissions(t *testing.T) {
	u, err := RegisterUser("all_submissions_test_master", "go", "golang.com")
	if err != nil {
		t.Error("User could not be registered", err)
		t.Fail()
	}

	c, err := u.NewClass("xyz", "xxx")
	if err != nil {
		t.Error("Class could not be registered", err)
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error("User could not be retrieved", err)
	}

	s.MakeAdmin()

	a, err := s.CreateWorkshop("A", "desc", "expl")
	if err != nil {
		t.Error("Could not make assignment")
	}

	sm := []Submission{}

	for i := 0; i < 10; i++ {
		u, err := RegisterUser(fmt.Sprintf("all_submissions_test_%v", i), "go", "golang.com")
		if err != nil {
			t.Error("User could not be registered", err)
			t.Fail()
		}

		stu, err := c.Invite(u)
		if err != nil {
			t.Error("User could not be invited", err)
			t.Fail()
		}

		s, err := stu.StartAssignment(a)
		if err != nil {
			t.Error("Could not start assignment")
		}

		sm = append(sm, s)

		defer u.Delete()
	}

	if len(sm) != 10 {
		t.Error("Wrong amount of submissions")
	}

	a.Delete()
	u.Delete()
}

func TestWorkshopPages(t *testing.T) {
	u, err := RegisterUser("test_workshop_pages", "xD", "someeon")
	if err != nil {
		t.Error("Could not register user", err)
	}

	c, err := u.NewClass("workshop pages", "something pages")
	if err != nil {
		t.Error("COuld not make new classs", err)
	}

	s, err := c.Student(u)
	if err != nil {
		t.Error("Could not get student", err)
	}

	w, err := s.CreateWorkshop("the page test workshop", "something description", "something description")
	if err != nil {
		t.Error("COuld not create workshop")
	}

	ps, err := w.Pages()
	if err != nil {
		t.Error("Could not get pages", err)
	}

	if len(ps) != 0 {
		t.Error("Wrong amount of pages! Expected 0, got ", len(ps))
	}

	p, err := w.CreatePage("Welcome to pages", "This is the beginning.")
	if err != nil {
		t.Error("Could not create page: ", err)
	}

	err = p.Edit("This is the end.")
	if err != nil {
		t.Error("Could not edit page")
	}

	ps, err = w.Pages()
	if err != nil {
		t.Error("Could not get pages", err)
	}

	if len(ps) != 1 {
		t.Error("Wrong amount of pages! Expected 1, got ", len(ps))
	}

	correct := []WorkshopPage{p}

	for i := 0; i < 10; i++ {
		p, err = w.CreatePage(fmt.Sprintf("Page #%v", i), "something description")
		if err != nil {
			t.Error("Could not add page #", i)
		}

		correct = append(correct, p)
	}

	ps, err = w.Pages()
	if err != nil {
		t.Error("Coudl not get pages")
	}

	if len(ps) != len(correct) {
		t.Error("Pages are not the same length")
	}

	for i := 0; i < len(correct); i++ {
		if ps[i] != correct[i] {
			t.Error("Pages are not identical,", ps[i], "is not", correct[i])
		}
	}
}
