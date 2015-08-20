package models

func GetStudentByID(id int64) (Student, error) {
	s := Student{}

	row := db.QueryRow("SELECT id, user, permission_level, class FROM users WHERE id = ?", id)
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
