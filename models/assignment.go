package models

import (
	"database/sql"
	"errors"
	"time"
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
