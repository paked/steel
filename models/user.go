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

func LoginUser(username, password string) (User, error) {
	u := User{}
	err := checkCredentials(username, password)
	if err != nil {
		return u, err
	}

	row := db.QueryRow("SELECT username, password_hash, id, email, permission_level FROM users WHERE username = ? AND password_hash = ?", username, password)
	err = row.Scan(&u.Username, &u.PasswordHash, &u.ID, &u.Email, &u.Permissions)

	if err != nil {
		return u, err
	}

	return u, nil
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

func checkCredentials(username, password string) error {
	if username == "" {
		return errors.New("Not a valid username")
	}

	if password == "" {
		return errors.New("Not a valid password")
	}

	return nil
}
