package models

import (
	"database/sql"
	"errors"
)

type Class struct {
	ID          int64
	Name        string
	Description string
	Image       string
}

func NewClass(name, description string) (Class, error) {
	c := Class{
		Name:        name,
		Description: description,
	}

	result, err := db.Exec("INSERT INTO classes (name, description) VALUES (?, ?)", c.Name, c.Description)
	if err != nil {
		return c, err
	}

	c.ID, err = result.LastInsertId()
	if err != nil {
		return c, err
	}

	return c, nil
}

func (c *Class) AddUser(u User) error {
	rows := db.QueryRow("SELECT id FROM students WHERE user = ? AND class = ?", u.ID, c.ID)
	err := rows.Scan() // why does this work...
	if err != sql.ErrNoRows {
		return errors.New("That user is already in this class")
	}

	_, err = db.Exec("INSERT INTO students (user, class) VALUES (?, ?)", u.ID, c.ID)

	return err
}
