package models

import (
	"database/sql"
	"errors"
)

func GetClassByID(id int64) (Class, error) {
	c := Class{}

	row := db.QueryRow("SELECT id, name, description, image_url FROM classes WHERE id = ?", id)

	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Image)
	if err != nil {
		return c, err
	}

	return c, nil
}

type Class struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image_url"`
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

	res, err := db.Exec("INSERT INTO students (user, class, permission_level) VALUES (?, ?, ?)", s.User, s.Class, s.Permissions)

	s.ID, err = res.LastInsertId()

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

func (c *Class) Student(u User) (Student, error) {
	s := Student{
		Class: c.ID,
		User:  u.ID,
	}

	row := db.QueryRow("SELECT id, user, permission_level FROM students WHERE class = ? AND user = ?", c.ID, u.ID)

	err := row.Scan(&s.ID, &s.User, &s.Permissions)

	return s, err
}
