package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"time"
)

type PathItem struct {
	X    uint16
	Y    uint16
	Prev *PathItem
}

type MazeController interface {
	GetUserMazes(userId string) ([]*Maze, error)
	Generate(width uint16, height uint16) (*Maze, error)
	CreateMaze(userId string, maze *Maze) (string, error)
	DrawMaze(maze *Maze, pathItem *PathItem, w io.Writer) error
	FindSolutionById(mazeId string, steps string) (*MazeSolution, error)
}

func NewMazeController(mazeRepository MazeRepository, mazeSolver MazeSolver) MazeController {
	return &mazeControllerImpl{
		mazeRepository: mazeRepository,
		mazeSolver:     mazeSolver,
	}
}

type mazeControllerImpl struct {
	mazeRepository MazeRepository
	mazeSolver     MazeSolver
}

func (m *mazeControllerImpl) FindSolutionById(mazeId string, steps string) (*MazeSolution, error) {
	maze, err := m.mazeRepository.SelectById(mazeId)
	if err != nil {
		return nil, err
	}

	solutions, err := m.mazeSolver.FindSolutions(maze)
	if err != nil {
		return nil, err
	}
	if len(solutions) == 0 {
		return nil, errors.New("No solution found")
	}

	// For steps == min return the shortest solution
	if steps == "min" {
		min := solutions[0]
		for _, solution := range solutions {
			if solution.Length < min.Length {
				min = solution
			}
		}
		return &min, nil
	}

	// Return the longest solution
	max := solutions[0]
	for _, solution := range solutions {
		if solution.Length > max.Length {
			max = solution
		}
	}
	return &max, nil
}

func (m *mazeControllerImpl) GetUserMazes(userId string) ([]*Maze, error) {
	return m.mazeRepository.SelectAllByUserId(userId)
}

func (m *mazeControllerImpl) Generate(width uint16, height uint16) (*Maze, error) {
	maze := &Maze{}
	maze.InitWalls(width, height)
	for i := 0; i < len(maze.Walls); i++ {
		maze.Walls[i] = 255
	}

	// Set entrace
	maze.EntranceX = 1 + uint16(rand.Int31n(int32(maze.GridWidth-2)))
	maze.EntranceY = 0
	maze.SetWall(maze.EntranceX, maze.EntranceY, false)

	// Fill the maze
	generateMaze_nextCell(maze.EntranceX, maze.EntranceY+1, maze)

	// Set exit
	longestPath := findLongestPathFromEntrace(maze)
	path := longestPath
	for path != nil {
		if path.X == 1 {
			maze.SetWall(0, path.Y, false)
			break
		}
		if path.Y == 1 {
			maze.SetWall(path.X, 0, false)
			break
		}
		if path.X == maze.GridWidth-2 {
			maze.SetWall(maze.GridWidth-1, path.Y, false)
			break
		}
		if path.Y == maze.GridHeight-2 {
			maze.SetWall(path.X, maze.GridHeight-1, false)
			break
		}
		path = path.Prev
	}

	return maze, nil
}

func (m *mazeControllerImpl) CreateMaze(userId string, maze *Maze) (string, error) {
	if maze == nil || maze.GridWidth == 0 || maze.GridHeight == 0 {
		return "", errors.New("Invalid maze size")
	}
	if maze.EntranceX != 0 && maze.EntranceY != 0 && maze.EntranceX != maze.GridWidth-1 && maze.EntranceY != maze.GridHeight-1 {
		return "", errors.New("Invalid entrance")
	}

	// Check that at least one solution can be found
	solutions, err := m.mazeSolver.FindSolutions(maze)
	if err != nil {
		return "", err
	}
	if len(solutions) == 0 {
		return "", errors.New("No solution found")
	}

	// Make sure there is only one exit
	exit := solutions[0].Exit
	for _, solution := range solutions {
		if solution.Exit != exit {
			return "", errors.New("Multiple exits found")
		}
	}

	// Make sure that the exit is on the bottom edge
	_, exitY, err := readPosition(exit)
	if err != nil {
		return "", err
	}
	if exitY != maze.GridHeight-1 {
		return "", errors.New("Exit is not on the bottom edge")
	}

	return m.mazeRepository.Insert(userId, maze.EntranceX, maze.EntranceY, maze.GridWidth, maze.GridHeight, maze.Walls)
}

func (m *mazeControllerImpl) DrawMaze(maze *Maze, path *PathItem, w io.Writer) error {
	var img = image.NewRGBA(image.Rect(0, 0, int(maze.GridWidth), int(maze.GridHeight)))

	// Draw maze
	for x := uint16(0); x < maze.GridWidth; x++ {
		for y := uint16(0); y < maze.GridHeight; y++ {
			if maze.IsWall(x, y) {
				img.Set(int(x), int(y), color.RGBA{0, 0, 0, 255})
			} else {
				img.Set(int(x), int(y), color.RGBA{255, 255, 255, 255})
			}
		}
	}

	// Draw path
	for path != nil {
		img.Set(int(path.X), int(path.Y), color.RGBA{255, 0, 0, 255})
		path = path.Prev
	}

	png.Encode(w, img)
	return nil
}

func generateMaze_nextCell(x uint16, y uint16, maze *Maze) {
	// We are at the right border
	if x == maze.GridWidth-1 {
		return
	}

	// We are at the left border
	if x == 0 {
		return
	}

	// We are at the bottom border
	if y == maze.GridHeight-1 {
		return
	}

	// We are at the top border
	if y == 0 {
		return
	}

	// Check if this cell is already set
	if !maze.IsWall(x, y) {
		return
	}

	// Check if neighbors are already set
	if countNeighborsWithoutWall(x, y, maze) > 1 {
		return
	}

	maze.SetWall(x, y, false)

	// Make all possible branches
	var branches []func()
	branches = append(branches, func() {
		generateMaze_nextCell(x+1, y, maze)
	})
	branches = append(branches, func() {
		generateMaze_nextCell(x-1, y, maze)
	})
	branches = append(branches, func() {
		generateMaze_nextCell(x, y+1, maze)
	})
	branches = append(branches, func() {
		generateMaze_nextCell(x, y-1, maze)
	})

	// Shuffle branches and call them
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(branches), func(i, j int) { branches[i], branches[j] = branches[j], branches[i] })
	for _, branch := range branches {
		branch()
	}
}

func countNeighborsWithoutWall(x uint16, y uint16, maze *Maze) uint8 {
	count := uint8(0)
	if !maze.IsWall(x+1, y) {
		count++
	}
	if !maze.IsWall(x-1, y) {
		count++
	}
	if !maze.IsWall(x, y+1) {
		count++
	}
	if !maze.IsWall(x, y-1) {
		count++
	}
	return count
}
