package models

import "time"

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
