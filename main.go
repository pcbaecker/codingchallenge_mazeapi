package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/speps/go-hashids"
)

type ApiEndpoint interface {
	Init(router *mux.Router)
}

type Response struct {
	Persons []Person `json:"persons"`
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	log.Println("Preparing server ...")

	// Setup https://hashids.org/
	hashIdData := hashids.NewData()
	hashIdData.Salt = "f948b5c5723ac19643226d73d52662d6"
	hashIdData.MinLength = 3
	hashId, err := hashids.NewWithData(hashIdData)
	if err != nil {
		panic(err)
	}

	// Setup database
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}

	// Setup router
	router := mux.NewRouter()

	// Setup services
	mazeSolver := NewMazeSolver()
	sessionRepo := NewSessionRepository(db)
	userRepo := NewUserRepository(db)
	mazeRepo := NewMazeRepository(db, hashId)
	userController := NewUserController(userRepo, sessionRepo)
	mazeController := NewMazeController(mazeRepo, mazeSolver)
	userApi := NewUserApi(userController)
	mazeApi := NewMazeApi(mazeController, userController)
	userApi.Init(router)
	mazeApi.Init(router)

	// Start server
	log.Println("Starting server ...")
	http.Handle("/", router)
	http.ListenAndServe(":8080", router)

}
