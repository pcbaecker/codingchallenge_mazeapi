package main

import (
	"sort"
)

type MazeSolution struct {
	Length uint64   `json:"length"`
	Path   []string `json:"path"`
	Exit   string   `json:"exit"`
}

type MazeSolver interface {
	FindSolutions(maze *Maze) ([]MazeSolution, error)
}

func NewMazeSolver() MazeSolver {
	return &mazeSolverImpl{}
}

type mazeSolverImpl struct {
}

func (m *mazeSolverImpl) FindSolutions(maze *Maze) ([]MazeSolution, error) {
	firstPathItem := &PathItem{
		X: maze.EntranceX,
		Y: maze.EntranceY,
	}
	var completePaths []*PathItem
	findAllPathsToExit_cell(maze, firstPathItem, &completePaths)

	var result []MazeSolution
	for _, path := range completePaths {
		pathLength := uint64(0)
		var pathStrings []string
		exit := toApiAddress(path.X, path.Y)
		for path != nil {
			pathStrings = append(pathStrings, toApiAddress(path.X, path.Y))
			path = path.Prev
			pathLength++
		}
		reverseSlice(pathStrings)
		result = append(result, MazeSolution{
			Length: pathLength,
			Path:   pathStrings,
			Exit:   exit,
		})
	}

	return result, nil
}

func reverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func findLongestPathFromEntrace(maze *Maze) *PathItem {
	firstPathItem := &PathItem{
		X: maze.EntranceX,
		Y: maze.EntranceY,
	}
	var completePaths []*PathItem
	findLongestPath_cell(maze, firstPathItem, &completePaths)

	var longestPath *PathItem
	longestPathLength := uint16(0)
	for _, path := range completePaths {
		tmpLength := countPathLength(path)
		if tmpLength > longestPathLength {
			longestPathLength = tmpLength
			longestPath = path
		}
	}

	return longestPath
}

func findLongestPath_cell(maze *Maze, lastPathItem *PathItem, completePaths *[]*PathItem) {
	foundAtLeastOne := false

	// Look right
	if lastPathItem.X < maze.GridWidth-1 &&
		!isInPath(lastPathItem.X+1, lastPathItem.Y, lastPathItem) &&
		!maze.IsWall(lastPathItem.X+1, lastPathItem.Y) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X + 1,
			Y:    lastPathItem.Y,
			Prev: lastPathItem,
		}
		findLongestPath_cell(maze, nextPathItem, completePaths)
		foundAtLeastOne = true
	}

	// Look down
	if lastPathItem.Y < maze.GridHeight-1 &&
		!isInPath(lastPathItem.X, lastPathItem.Y+1, lastPathItem) &&
		!maze.IsWall(lastPathItem.X, lastPathItem.Y+1) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X,
			Y:    lastPathItem.Y + 1,
			Prev: lastPathItem,
		}
		findLongestPath_cell(maze, nextPathItem, completePaths)
		foundAtLeastOne = true
	}

	// Look left
	if lastPathItem.X > 0 &&
		!isInPath(lastPathItem.X-1, lastPathItem.Y, lastPathItem) &&
		!maze.IsWall(lastPathItem.X-1, lastPathItem.Y) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X - 1,
			Y:    lastPathItem.Y,
			Prev: lastPathItem,
		}
		findLongestPath_cell(maze, nextPathItem, completePaths)
		foundAtLeastOne = true
	}

	// Look up
	if lastPathItem.Y > 0 &&
		!isInPath(lastPathItem.X, lastPathItem.Y-1, lastPathItem) &&
		!maze.IsWall(lastPathItem.X, lastPathItem.Y-1) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X,
			Y:    lastPathItem.Y - 1,
			Prev: lastPathItem,
		}
		findLongestPath_cell(maze, nextPathItem, completePaths)
		foundAtLeastOne = true
	}

	if !foundAtLeastOne {
		// We are at the end of a path
		*completePaths = append(*completePaths, lastPathItem)
	}
}

func isInPath(x uint16, y uint16, path *PathItem) bool {
	for path != nil {
		if path.X == x && path.Y == y {
			return true
		}
		path = path.Prev
	}
	return false
}

func countPathLength(path *PathItem) uint16 {
	count := uint16(0)
	for path != nil {
		count++
		path = path.Prev
	}
	return count
}

func findAllPathsToExit_cell(maze *Maze, lastPathItem *PathItem, completePaths *[]*PathItem) {
	// Look for an exit
	if lastPathItem.Y == maze.GridHeight-1 {
		// Every path that reaches the bottom is a complete path
		*completePaths = append(*completePaths, lastPathItem)
		return
	}

	// Look right
	if lastPathItem.X < maze.GridWidth-1 &&
		!isInPath(lastPathItem.X+1, lastPathItem.Y, lastPathItem) &&
		!maze.IsWall(lastPathItem.X+1, lastPathItem.Y) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X + 1,
			Y:    lastPathItem.Y,
			Prev: lastPathItem,
		}
		findAllPathsToExit_cell(maze, nextPathItem, completePaths)
	}

	// Look down
	if lastPathItem.Y < maze.GridHeight-1 &&
		!isInPath(lastPathItem.X, lastPathItem.Y+1, lastPathItem) &&
		!maze.IsWall(lastPathItem.X, lastPathItem.Y+1) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X,
			Y:    lastPathItem.Y + 1,
			Prev: lastPathItem,
		}
		findAllPathsToExit_cell(maze, nextPathItem, completePaths)
	}

	// Look left
	if lastPathItem.X > 0 &&
		!isInPath(lastPathItem.X-1, lastPathItem.Y, lastPathItem) &&
		!maze.IsWall(lastPathItem.X-1, lastPathItem.Y) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X - 1,
			Y:    lastPathItem.Y,
			Prev: lastPathItem,
		}
		findAllPathsToExit_cell(maze, nextPathItem, completePaths)
	}

	// Look up
	if lastPathItem.Y > 0 &&
		!isInPath(lastPathItem.X, lastPathItem.Y-1, lastPathItem) &&
		!maze.IsWall(lastPathItem.X, lastPathItem.Y-1) {
		nextPathItem := &PathItem{
			X:    lastPathItem.X,
			Y:    lastPathItem.Y - 1,
			Prev: lastPathItem,
		}
		findAllPathsToExit_cell(maze, nextPathItem, completePaths)
	}
}
