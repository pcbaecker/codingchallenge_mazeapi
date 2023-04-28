package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Maze struct {
	Id         string
	EntranceX  uint16
	EntranceY  uint16
	GridWidth  uint16
	GridHeight uint16
	Walls      []byte
}

func (maze *Maze) InitWalls(width uint16, height uint16) {
	maze.GridWidth = width
	maze.GridHeight = height
	numberOfBytes := ((maze.GridWidth * maze.GridHeight) / 8) + 1
	maze.Walls = make([]byte, numberOfBytes)
	for i := 0; i < len(maze.Walls); i++ {
		maze.Walls[i] = 0
	}
}

func (m *Maze) GetByteAddress(x uint16, y uint16) uint32 {
	return uint32((y*m.GridWidth + x) / 8)
}

func (m *Maze) GetBitAddress(x uint16, y uint16) uint8 {
	return uint8((y*m.GridWidth + x) % 8)
}

func (m *Maze) SetWall(x uint16, y uint16, isWall bool) {
	if isWall {
		m.Walls[m.GetByteAddress(x, y)] |= 1 << m.GetBitAddress(x, y)
	} else {
		m.Walls[m.GetByteAddress(x, y)] &= ^(1 << m.GetBitAddress(x, y))
	}
}

func (m *Maze) IsWall(x uint16, y uint16) bool {
	w := m.Walls[m.GetByteAddress(x, y)]
	return w&(1<<m.GetBitAddress(x, y)) != 0
}

func (m *Maze) WallsToStrings() []string {
	var result []string
	for x := uint16(0); x < m.GridWidth; x++ {
		for y := uint16(0); y < m.GridHeight; y++ {
			if m.IsWall(x, y) {
				result = append(result, toApiAddress(x, y))
			}
		}
	}
	return result
}

func toApiAddress(x uint16, y uint16) string {
	column := string(rune('A' + x))
	return column + strconv.Itoa(int(y+1))
}

var POSITION_VALIDATION_PATTERN = regexp.MustCompile("^[A-Z]{1,2}[0-9]{1,2}$")
var POSITION_COLUMN_PATTERN = regexp.MustCompile("^[A-Z]{1,2}")
var POSITION_ROW_PATTERN = regexp.MustCompile("[0-9]{1,2}$")

func readPosition(position string) (uint16, uint16, error) {
	if !POSITION_VALIDATION_PATTERN.MatchString(position) {
		return 0, 0, errors.New("Invalid position")
	}

	// Column
	column := uint16(0)
	columnChars := POSITION_COLUMN_PATTERN.FindString(position)
	for i := 0; i < len(columnChars); i++ {
		v := uint16(columnChars[i] - byte('A'))
		column += v + uint16(i)*uint16(26)
	}

	// Row
	row, err := strconv.Atoi(POSITION_ROW_PATTERN.FindString(position))
	if err != nil {
		return 0, 0, err
	}

	return column, uint16(row - 1), nil
}

var GRIDSIZE_VALIDATION_PATTERN = regexp.MustCompile("^[0-9]{1,}x[0-9]{1,}$")

func readGridSize(gridsize string) (uint16, uint16, error) {
	if !GRIDSIZE_VALIDATION_PATTERN.MatchString(gridsize) {
		return 0, 0, errors.New("Invalid grid size")
	}
	sizes := strings.Split(gridsize, "x")
	width, err := strconv.ParseUint(sizes[0], 10, 16)
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.ParseUint(sizes[1], 10, 16)
	if err != nil {
		return 0, 0, err
	}

	return uint16(width), uint16(height), nil
}

func applyStringWalls(maze *Maze, walls []string) error {
	for _, wall := range walls {
		x, y, err := readPosition(wall)
		if err != nil {
			return err
		}
		if x >= maze.GridWidth || y >= maze.GridHeight {
			return errors.New("wall position out of range x=" + strconv.Itoa(int(x)) + " y=" + strconv.Itoa(int(y)))
		}
		maze.SetWall(x, y, true)
	}
	return nil
}

func stringsToPath(path []string) (*PathItem, error) {
	var pathItem *PathItem
	for _, p := range path {
		x, y, err := readPosition(p)
		if err != nil {
			return nil, err
		}

		pathItem = &PathItem{X: x, Y: y, Prev: pathItem}
	}

	return pathItem, nil
}
