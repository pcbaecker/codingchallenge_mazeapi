package main

import (
	"database/sql"
	"log"

	"github.com/speps/go-hashids"
)

const MAZE_REPO_CREATE_TABLE = "CREATE TABLE IF NOT EXISTS mazes (id INTEGER PRIMARY KEY, user_id VARCHAR(255) NOT NULL, entrance_x INTEGER NOT NULL, entrance_y INTEGER NOT NULL, grid_width INTEGER NOT NULL, grid_height INTEGER NOT NULL, walls BLOB NOT NULL)"

type MazeRepository interface {
	Insert(user_id string, entranceX uint16, entranceY uint16, gridWidth uint16, gridHeight uint16, walls []byte) (string, error)
	SelectById(id string) (*Maze, error)
	SelectAllByUserId(userId string) ([]*Maze, error)
	CountAll() (uint64, error)
}

func NewMazeRepository(db *sql.DB, hashid *hashids.HashID) MazeRepository {
	_, err := db.Exec(MAZE_REPO_CREATE_TABLE)
	if err != nil {
		panic(err)
	}
	return &mazeRepositoryImpl{
		db:     db,
		hashid: hashid,
	}
}

type mazeRepositoryImpl struct {
	db     *sql.DB
	hashid *hashids.HashID
}

func (m *mazeRepositoryImpl) Insert(user_id string, entranceX uint16, entranceY uint16, gridWidth uint16, gridHeight uint16, walls []byte) (string, error) {
	result, err := m.db.Exec("INSERT INTO mazes (user_id,entrance_x, entrance_y, grid_width, grid_height, walls) VALUES (?, ?, ?, ?, ?, ?)", user_id, entranceX, entranceY, gridWidth, gridHeight, walls)
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return m.hashid.EncodeInt64([]int64{int64(id)})
}

func (m *mazeRepositoryImpl) SelectById(id string) (*Maze, error) {
	decoded, err := m.hashid.DecodeInt64WithError(id)
	if err != nil {
		return nil, err
	}

	rows, err := m.db.Query("SELECT id, entrance_x, entrance_y, grid_width, grid_height, walls FROM mazes WHERE id = ?", decoded[0])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	maze := &Maze{}
	var walls []byte
	var resultId uint64
	err = rows.Scan(&resultId, &maze.EntranceX, &maze.EntranceY, &maze.GridWidth, &maze.GridHeight, &walls)
	if err != nil {
		return nil, err
	}
	maze.Walls = walls
	log.Println(resultId)
	maze.Id, err = m.hashid.EncodeInt64([]int64{int64(resultId)})
	return maze, err
}

func (m *mazeRepositoryImpl) CountAll() (uint64, error) {
	rows, err := m.db.Query("SELECT COUNT(*) FROM mazes")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}

	var count uint64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *mazeRepositoryImpl) SelectAllByUserId(userId string) ([]*Maze, error) {
	rows, err := m.db.Query("SELECT id, entrance_x, entrance_y, grid_width, grid_height, walls FROM mazes WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mazes := []*Maze{}
	for rows.Next() {
		maze := &Maze{}
		var walls []byte
		var resultId uint64
		err = rows.Scan(&resultId, &maze.EntranceX, &maze.EntranceY, &maze.GridWidth, &maze.GridHeight, &walls)
		if err != nil {
			return nil, err
		}
		maze.Walls = walls
		maze.Id, err = m.hashid.EncodeInt64([]int64{int64(resultId)})
		if err != nil {
			return nil, err
		}
		mazes = append(mazes, maze)
	}

	return mazes, nil
}
