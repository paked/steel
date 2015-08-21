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

func (c *Class) Invite(u User) (Student, error) {
	s := Student{}

	row := db.QueryRow("SELECT id FROM students WHERE user = ? AND class = ?", u.ID, c.ID)
	err := row.Scan() // why does this work...
	if err != sql.ErrNoRows {
		return s, errors.New("That user is already in this class")
	}

	s = Student{
		User:        u.ID,
		Class:       c.ID,
		Permissions: DefaultPermissions,
	}

	_, err = db.Exec("INSERT INTO students (user, class, permission_level) VALUES (?, ?, ?)", s.User, s.Class, s.Permissions)

	return s, err
}

func (c *Class) Students() ([]Student, error) {
	var st []Student
	rows, err := db.Query("SELECT id, user, permission_level FROM students WHERE class = ?", c.ID)
	if err != nil {
		return st, err
	}

	defer rows.Close()

	for rows.Next() {
		s := Student{
			Class: c.ID,
		}

		err = rows.Scan(&s.ID, &s.User, &s.Permissions)
		if err != nil {
			return st, err
		}

		st = append(st, s)
	}

	return st, nil
}
