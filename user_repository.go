package main

import (
	"database/sql"

	"github.com/speps/go-hashids"
)

const USER_REPO_CREATE_TABLE = "CREATE TABLE IF NOT EXISTS users (username VARCHAR(256) NOT NULL PRIMARY KEY, password_hash VARCHAR(256))"

type UserRepository interface {
	Insert(username string, passwordHash string) (string, error)
	SelectByUsername(username string) (string, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	_, err := db.Exec(USER_REPO_CREATE_TABLE)
	if err != nil {
		panic(err)
	}
	return &userRepositoryImpl{
		db: db,
	}
}

type userRepositoryImpl struct {
	db     *sql.DB
	hashid *hashids.HashID
}

func (u *userRepositoryImpl) Insert(username string, passwordHash string) (string, error) {
	_, err := u.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (u *userRepositoryImpl) SelectByUsername(username string) (string, error) {
	var passwordHash string
	err := u.db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return passwordHash, nil
}
