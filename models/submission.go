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

func GetSubmission(id int64) (Submission, error) {
	s := Submission{
		ID: id,
	}

	row := db.QueryRow("SELECT team_name, thoughts, assignment FROM submissions WHERE id = ?", s.ID)
	err := row.Scan(&s.TeamName, &s.Thoughts, &s.Assignment)

	return s, err
}

func (sub *Submission) Invite(s Student) error {
	row := db.QueryRow("SELECT student FROM team_members WHERE student = ? AND submission = ?", s.ID, sub.ID)

	var uid *int64

	err := row.Scan(uid)
	if uid != nil || (err != nil && err != sql.ErrNoRows) {
		return errors.New("This user is alreayd a member of another team")
	}

	_, err = db.Exec("INSERT INTO team_members (student, submission, assignment) VALUES (?, ?, ?)", s.ID, sub.ID, sub.Assignment)

	return err
}

func (sub *Submission) Members() ([]Student, error) {
	var us []Student
	rows, err := db.Query("SELECT student FROM team_members WHERE submission = ? ", sub.ID)
	if err != nil {
		return us, err
	}

	defer rows.Close()

	for rows.Next() {
		var uid int64
		err = rows.Scan(&uid)
		if err != nil {
			return us, err
		}

		s, err := GetStudentByID(uid)
		if err != nil {
			return us, err
		}

		us = append(us, s)
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
