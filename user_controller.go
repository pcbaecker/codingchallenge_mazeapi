package main

import (
	"errors"
	"log"
	"math/rand"
)

type UserController interface {
	CreateUser(username string, password string) (string, error)
	Login(username string, password string) (string, error)
	GetUserForSession(sessionId string) (string, error)
}

func NewUserController(userRepository UserRepository, sessionRepository SessionRepository) UserController {
	return &userControllerImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

type userControllerImpl struct {
	userRepository    UserRepository
	sessionRepository SessionRepository
}

func (u *userControllerImpl) CreateUser(username string, password string) (string, error) {
	if (len(username) < 3) || (len(password) < 3) {
		return "", errors.New("username and password must be at least 3 characters long")
	}

	hash, err := hashPassword(password)
	log.Println(hash)
	if err != nil {
		return "", err
	}
	userId, err := u.userRepository.Insert(username, hash)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (u *userControllerImpl) Login(username string, password string) (string, error) {
	// Find user in db
	hash, err := u.userRepository.SelectByUsername(username)
	if err != nil {
		return "", err
	}
	if hash == "" {
		return "", errors.New("invalid username or password")
	}

	// Check password
	passwordOk, err := comparePasswordAndHash(password, hash)
	if err != nil {
		return "", err
	}
	if !passwordOk {
		return "", errors.New("invalid username or password")
	}

	// Create session
	sessionId := randStringBytesRmndr(32)
	err = u.sessionRepository.Insert(sessionId, username)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

const letters = "abcdefghijklmnopqrstuvwxyz1234567890"

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func (u *userControllerImpl) GetUserForSession(sessionId string) (string, error) {
	return u.sessionRepository.SelectBySessionId(sessionId)
}
