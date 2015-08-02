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

func GetAssignment(id int64) (Assignment, error) {
	a := Assignment{
		ID: id,
	}

	var timeString string

	row := db.QueryRow("SELECT name, description, explanation, due FROM assignments WHERE id = ?", a.ID)
	err := row.Scan(&a.Name, &a.Description, &a.Explanation, &timeString)
	if err != nil {
		return a, err
	}

	a.Due, err = time.Parse(timeFormat, timeString)
	if err != nil {
		return a, err
	}

	return a, nil
}
