package models

import (
	"database/sql"
	"time"
)

type User struct {
	id        int
	name      string
	email     string
	password  string
	isAdmin   bool
	createdAt string
	updatedAt string
}

func New(name string, email string, password string, isAdmin bool) *User {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return &User{
		name:      name,
		email:     email,
		password:  password,
		isAdmin:   isAdmin,
		createdAt: currentTime,
		updatedAt: currentTime,
	}
}

func (u User) Save(db *sql.DB) error {
	_, err := db.Exec(
		`INSERT INTO users (name, email, password, is_admin, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		u.name, u.email, u.password, u.isAdmin, u.createdAt, u.updatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u User) GetByEmail(db *sql.DB, email string) (User, error) {
	var user User
	err := db.
		QueryRow(`SELECT * FROM users WHERE email = ?`, email).
		Scan(&user.id, &user.name, &user.email, &user.password, &user.isAdmin, &user.createdAt, &user.updatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
