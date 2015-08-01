package models

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
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

	row := db.QueryRow("SELECT username, password_hash, id, email FROM users WHERE username = ? AND password_hash = ?", username, password)
	err = row.Scan(&u.Username, &u.PasswordHash, &u.ID, &u.Email)
	if err != nil {
		return u, err
	}

	return u, nil
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
