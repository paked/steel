package models

import (
	"database/sql"
	"errors"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

type Assignment struct {
	ID          int64
	Name        string
	Description string
	Explanation string
	Due         time.Time
}

func (a *Assignment) Delete() error {
	_, err := db.Exec("DELETE FROM assignments WHERE id = ?", a.ID)

	return err
}

type Submission struct {
	ID       int64
	Thoughts string
	TeamName string
}

type SubmissionMember struct {
	ID           int64
	SubmissionID int64
	UserID       int64
}

func GetAssignment(id int64) (Assignment, error) {
	a := Assignment{
		ID: id,
	}

	var unixTime int64

	row := db.QueryRow("SELECT name, description, explanation, due FROM assignments WHERE id = ?", a.ID)
	err := row.Scan(&a.Name, &a.Description, &a.Explanation, &unixTime)
	if err != nil {
		return a, err
	}

	a.Due = time.Unix(0, unixTime)
	if err != nil {
		return a, err
	}

	return a, nil
}

func GetSubmission(id int64) (Submission, error) {
	s := Submission{
		ID: id,
	}

	row := db.QueryRow("SELECT team_name, thoughts FROM submissions WHERE id = ?", s.ID)
	err := row.Scan(&s.TeamName, &s.Thoughts)

	return s, err
}

func (s *Submission) AddMember(id int64) error {
	row := db.QueryRow("SELECT user FROM team_members WHERE submission = ? AND user = ?", s.ID, id)

	var uid *int64
	err := row.Scan(uid)

	if uid != nil || (err != nil && err != sql.ErrNoRows) {
		return errors.New("This user is alreayd a member")
	}

	_, err = db.Exec("INSERT INTO team_members (submission, user) VALUES (?, ?)", s.ID, id)

	return err
}

func (s *Submission) Members() ([]User, error) {
	var us []User
	rows, err := db.Query("SELECT user FROM team_members WHERE submission = ?", s.ID)
	if err != nil {
		return us, err
	}

	defer rows.Close()

	for rows.Next() {
		var uid int64
		rows.Scan(&uid)
		u, err := GetUserByID("id", uid)
		if err != nil {
			return us, err
		}

		us = append(us, u)
	}

	return us, nil
}
