package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type MazeApiDao struct {
	Entrace  string   `json:"entrance"`
	GridSize string   `json:"gridSize"`
	Walls    []string `json:"walls"`
}

type MazeWithIdDao struct {
	MazeApiDao
	Id string `json:"id"`
}

type DrawMazeApiDao struct {
	MazeApiDao
	Path []string `json:"path"`
}

type MyMazesDao struct {
	Mazes []MazeWithIdDao `json:"mazes"`
}

type MazeSolutionResponse struct {
	Path []string `json:"path"`
}

func NewMazeApi(mazeController MazeController, userController UserController) ApiEndpoint {
	return &mazeApiImpl{
		mazeController: mazeController,
		userController: userController,
	}
}

type mazeApiImpl struct {
	mazeController MazeController
	userController UserController
}

func (m *mazeApiImpl) Init(router *mux.Router) {
	router.HandleFunc("/maze", m.GetMyMazes).Methods("GET")
	router.HandleFunc("/maze", m.CreateMaze).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/maze/draw", m.DrawMaze).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/maze/generate", m.Generate).Methods("GET")
	router.HandleFunc("/maze/{mazeId}/solution", m.FindSolutionById).Methods("GET")
}

func (m *mazeApiImpl) FindSolutionById(w http.ResponseWriter, r *http.Request) {
	_, err := getUserId(r, m.userController)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mazeId := mux.Vars(r)["mazeId"]
	if mazeId == "" {
		http.Error(w, "the mazeId must be provided", http.StatusBadRequest)
		return
	}

	stepsParam := orDefault(r.URL.Query().Get("steps"), "min")
	if stepsParam != "min" && stepsParam != "max" {
		http.Error(w, "the steps parameter must be either 'min' or 'max'", http.StatusBadRequest)
		return
	}

	solution, err := m.mazeController.FindSolutionById(mazeId, stepsParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(solution)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *mazeApiImpl) GetMyMazes(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r, m.userController)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mazes, err := m.mazeController.GetUserMazes(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	var response MyMazesDao
	for _, maze := range mazes {
		response.Mazes = append(response.Mazes, MazeWithIdDao{
			Id: maze.Id,
			MazeApiDao: MazeApiDao{
				Entrace:  toApiAddress(maze.EntranceX, maze.EntranceY),
				GridSize: fmt.Sprintf("%dx%d", maze.GridWidth, maze.GridHeight),
				Walls:    maze.WallsToStrings(),
			},
		})
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *mazeApiImpl) Generate(w http.ResponseWriter, r *http.Request) {
	_, err := getUserId(r, m.userController)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	widthStr := orDefault(r.URL.Query().Get("width"), "16")
	heightStr := orDefault(r.URL.Query().Get("height"), "16")

	// Validate
	width, err := strconv.ParseUint(widthStr, 10, 32)
	if err != nil {
		http.Error(w, "the width must be numeric", http.StatusBadRequest)
		return
	}
	height, err := strconv.ParseUint(heightStr, 10, 32)
	if err != nil {
		http.Error(w, "the height must be numeric", http.StatusBadRequest)
		return
	}

	// Process
	maze, err := m.mazeController.Generate(uint16(width), uint16(height))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toMazeApiDao(maze))
}

func (m *mazeApiImpl) CreateMaze(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r, m.userController)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Read request
	var mazeDao MazeApiDao
	err = json.NewDecoder(r.Body).Decode(&mazeDao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entranceX, entranceY, err := readPosition(mazeDao.Entrace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	width, height, err := readGridSize(mazeDao.GridSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	maze := &Maze{
		EntranceX: entranceX,
		EntranceY: entranceY,
	}
	maze.InitWalls(width, height)
	err = applyStringWalls(maze, mazeDao.Walls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process request
	mazeId, err := m.mazeController.CreateMaze(userId, maze)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"mazeId": "%s"}`, mazeId)))
}

func (m *mazeApiImpl) DrawMaze(w http.ResponseWriter, r *http.Request) {
	_, err := getUserId(r, m.userController)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Read request
	var mazeDao DrawMazeApiDao
	err = json.NewDecoder(r.Body).Decode(&mazeDao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entranceX, entranceY, err := readPosition(mazeDao.Entrace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	width, height, err := readGridSize(mazeDao.GridSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	maze := &Maze{
		EntranceX: entranceX,
		EntranceY: entranceY,
	}
	maze.InitWalls(width, height)
	err = applyStringWalls(maze, mazeDao.Walls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, b := range maze.Walls {
		fmt.Println(b)
	}
	pathItem, err := stringsToPath(mazeDao.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process request
	err = m.mazeController.DrawMaze(maze, pathItem, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
}

func orDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func toMazeApiDao(maze *Maze) *MazeApiDao {
	return &MazeApiDao{
		Entrace:  toApiAddress(maze.EntranceX, maze.EntranceY),
		GridSize: fmt.Sprintf("%dx%d", maze.GridWidth, maze.GridHeight),
		Walls:    maze.WallsToStrings(),
	}
}

var AUTHHEADER_VALID_PATTERN = regexp.MustCompile(`^Bearer [a-zA-Z0-9]{1,}$`)

func getUserId(r *http.Request, userController UserController) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 {
		return "", errors.New("invalid authorization header")
	}

	if !AUTHHEADER_VALID_PATTERN.MatchString(authHeader) {
		return "", errors.New("invalid authorization header pattern")
	}

	// Remove the 'Bearer '
	sessionId := authHeader[7:]

	// Find user for sessionid
	return userController.GetUserForSession(sessionId)
}
