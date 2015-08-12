package models

import (
	"database/sql"
	"errors"
)

type Submission struct {
	ID         int64  `json:"id"`
	Assignment int64  `json:"assignment"`
	Thoughts   string `json:"thoughts"`
	TeamName   string `json:"team_name"`
}

type SubmissionMember struct {
	ID           int64 `json:"id"`
	SubmissionID int64 `json:"submission"`
	Assignment   int64 `json:"assignment"`
	UserID       int64 `json:"user"`
}

func GetSubmission(id int64) (Submission, error) {
	s := Submission{
		ID: id,
	}

	row := db.QueryRow("SELECT team_name, thoughts, assignment FROM submissions WHERE id = ?", s.ID)
	err := row.Scan(&s.TeamName, &s.Thoughts, &s.Assignment)

	return s, err
}

func (s *Submission) AddMember(u User) error {
	row := db.QueryRow("SELECT user FROM team_members WHERE assignment = ? AND user = ?", s.Assignment, u.ID)

	var uid *int64
	err := row.Scan(uid)

	if uid != nil || (err != nil && err != sql.ErrNoRows) {
		return errors.New("This user is alreayd a member of another team")
	}

	_, err = db.Exec("INSERT INTO team_members (submission, user, assignment) VALUES (?, ?, ?)", s.ID, u.ID, s.Assignment)

	return err
}

func (s *Submission) Members() ([]User, error) {
	var us []User
	rows, err := db.Query("SELECT user FROM team_members WHERE submission = ? AND assignment = ?", s.ID, s.Assignment)
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

func (s *Submission) Rename(name string) error {
	_, err := db.Exec("UPDATE submissions SET team_name=? WHERE id = ?", name, s.ID)

	if err != nil {
		return err
	}

	s.TeamName = name

	return nil
}
