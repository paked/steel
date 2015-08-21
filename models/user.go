package models

import (
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

// These should probably go into the DB
const (
	DefaultPermissions = iota
	AdminPermissions
)

// A user represents one account, but not the entity which is in a *class*
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Salt         string `json:"-"`
	Email        string `json:"email"`
	Permissions  int    `json:"permissions"`
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

func GetUserByID(id int64) (User, error) {
	u := User{}

	row := db.QueryRow("SELECT id, username, password_hash, email, permission_level FROM users WHERE id = ?", id)
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

func (u *User) NewClass(name, description string) (Class, error) {
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

func checkCredentials(username, password string) error {
	if username == "" {
		return errors.New("Not a valid username")
	}

	if password == "" {
		return errors.New("Not a valid password")
	}

	return nil
}
