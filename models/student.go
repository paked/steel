package models

import (
	"errors"
	"time"
)

func GetStudentByID(id int64) (Student, error) {
	s := Student{}

	row := db.QueryRow("SELECT id, user, permission_level, class FROM students WHERE id = ?", id)
	err := row.Scan(&s.ID, &s.User, &s.Permissions, &s.Class)
	if err != nil {
		return s, err
	}

	return s, nil
}

// Student represents a student who is taking place in a class
type Student struct {
	ID          int64
	User        int64
	Permissions int
	Class       int64
}

func (s *Student) Delete() error {
	_, err := db.Exec("DELETE FROM students WHERE ID = ?", s.ID)

	return err
}

func (s *Student) IsAdmin() bool {
	if s.Permissions == AdminPermissions {
		return true
	}

	return false
}

func (s *Student) MakeAdmin() error {
	return s.updatePermissions(AdminPermissions)
}

func (s *Student) DemoteAdmin() error {
	return s.updatePermissions(DefaultPermissions)
}

func (s *Student) updatePermissions(level int) error {
	_, err := db.Exec("UPDATE students SET permission_level=? WHERE id = ?", level, s.ID)

	if err != nil {
		return err
	}

	s.Permissions = level

	return nil
}

func (s *Student) DueAssignments(before time.Time) ([]Assignment, error) {
	var as []Assignment
	rows, err := db.Query("SELECT id, name, description, explanation, due, class FROM assignments WHERE due < ? AND class = ?", before.UnixNano(), s.Class)
	if err != nil {
		return as, err
	}

	defer rows.Close()

	for rows.Next() {
		a := Assignment{}
		var unixTime int64

		err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Explanation, &unixTime, &a.Class)
		if err != nil {
			return as, err
		}

		a.Due = time.Unix(0, unixTime)

		as = append(as, a)
	}

	return as, nil
}

func (s *Student) CreateAssignment(name, description, explanation string) (Assignment, error) {
	a := Assignment{}

	if s.Permissions != AdminPermissions {
		return a, errors.New("Incorrect permissions")
	}

	a = Assignment{
		Name:        name,
		Description: description,
		Explanation: explanation,
		Due:         time.Now(),
	}

	result, err := db.Exec("INSERT INTO assignments (name, description, explanation, due, created_by, class) VALUES (?, ?, ?, ?, ?, ?)", a.Name, a.Description, a.Explanation, a.Due.UnixNano(), s.ID, s.Class)
	if err != nil {
		return a, err
	}

	a.ID, err = result.LastInsertId()

	return a, err
}

func (s *Student) StartAssignment(a Assignment) (Submission, error) {
	sub := Submission{
		TeamName:   "Assignment",
		Assignment: a.ID,
	}

	res, err := db.Exec("INSERT INTO submissions (team_name, assignment) VALUES (?, ?)", sub.TeamName, sub.Assignment)
	if err != nil {
		return sub, err
	}

	sub.ID, err = res.LastInsertId()
	if err != nil {
		return sub, err
	}

	err = sub.Invite(*s)

	return sub, err
}
