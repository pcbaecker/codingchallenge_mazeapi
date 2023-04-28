package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserApi(userController UserController) ApiEndpoint {
	return &userApiImpl{
		userController: userController,
	}
}

type userApiImpl struct {
	userController UserController
}

func (u *userApiImpl) Init(router *mux.Router) {
	router.HandleFunc("/user", u.CreateUser).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/login", u.Login).Methods("POST").Headers("Content-Type", "application/json")
}

func (u *userApiImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process request
	userId, err := u.userController.CreateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"userId": %s}`, userId)))
}

func (u *userApiImpl) Login(w http.ResponseWriter, r *http.Request) {
	// Read request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process request
	sessionId, err := u.userController.Login(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"sessionId": %s}`, sessionId)))
}
