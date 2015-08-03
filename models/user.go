package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// These should probably go into the DB
const (
	DefaultPermissions = iota
	AdminPermissions
)

var (
	db *sql.DB
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Salt         string
	Email        string
	Permissions  int
}

func GetUser(key, value string) (User, error) {
	u := User{}

	row := db.QueryRow("SELECT id, username, password_hash, email, permission_level FROM users WHERE "+key+" = ?", value)
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Permissions)
	if err != nil {
		return u, err
	}

	return u, nil
}

func RegisterUser(username, password, email string) (User, error) {
	u := User{}

	err := checkCredentials(username, password)
	if err != nil {
		return u, err
	}

	if email == "" {
		return u, errors.New("Not a valid email")
	}

	row := db.QueryRow("SELECT username FROM users WHERE username = ? AND password_hash = ?", username, password)
	row.Scan(&u.Username)

	if u.Username != "" {
		return u, errors.New("That username is already taken")
	}

	u = User{
		Username:     username,
		PasswordHash: password,
		Email:        email,
	}

	result, err := db.Exec("INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)", u.Username, u.PasswordHash, u.Email)
	if err != nil {
		return u, err
	}

	u.ID, err = result.LastInsertId()

	return u, err
}

func (u *User) Login(password string) (bool, error) {
	if u.PasswordHash == password {
		return true, nil
	}

	return false, nil
}

func (u *User) Delete() error {
	_, err := db.Exec("DELETE FROM users WHERE username = ?", u.Username)

	return err
}

func (u *User) IsAdmin() bool {
	if u.Permissions == AdminPermissions {
		return true
	}

	return false
}

func (u *User) MakeAdmin() error {
	return u.updatePermissions(AdminPermissions)
}

func (u *User) DemoteAdmin() error {
	return u.updatePermissions(DefaultPermissions)
}

func (u *User) updatePermissions(level int) error {
	_, err := db.Exec("UPDATE users SET permission_level=? WHERE id = ?", level, u.ID)

	if err != nil {
		return err
	}

	u.Permissions = level

	return nil
}

func (u *User) CreateAssignment(name, description, explanation string) (Assignment, error) {
	a := Assignment{
		Name:        name,
		Description: description,
		Explanation: explanation,
		Due:         time.Now(),
	}

	result, err := db.Exec("INSERT INTO assignments (name, description, explanation, due, created_by) VALUES (?, ?, ?, ?, ?)", a.Name, a.Description, a.Explanation, a.Due, u.ID)
	if err != nil {
		return a, err
	}

	a.ID, err = result.LastInsertId()

	return a, err
}

func (u *User) StartAssignment(id int64) (Submission, error) {
	s := Submission{
		TeamName: u.Username + "'s Assignment",
	}

	res, err := db.Exec("INSERT INTO submissions (team_name) VALUES (?)", s.TeamName)
	if err != nil {
		return s, err
	}

	s.ID, err = res.LastInsertId()

	return s, err
}

func checkCredentials(username, password string) error {
	if username == "" {
		return errors.New("Not a valid username")
	}

	if password == "" {
		return errors.New("Not a valid password")
	}

	return nil
}