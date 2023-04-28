package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MazeSolverMock struct {
	mock.Mock
}

func (m *MazeSolverMock) FindSolutions(maze *Maze) ([]MazeSolution, error) {
	args := m.Called(maze)
	return args.Get(0).([]MazeSolution), args.Error(1)
}

type MazeRepositoryMock struct {
	mock.Mock
}

func (m *MazeRepositoryMock) Insert(userId string, entranceX uint16, entranceY uint16, gridWidth uint16, gridHeight uint16, walls []byte) (string, error) {
	args := m.Called(userId, entranceX, entranceY, gridWidth, gridHeight, walls)
	return args.String(0), args.Error(1)
}
func (m *MazeRepositoryMock) SelectById(id string) (*Maze, error) {
	args := m.Called(id)
	return args.Get(0).(*Maze), args.Error(1)
}
func (m *MazeRepositoryMock) CountAll() (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MazeRepositoryMock) SelectAllByUserId(userId string) ([]*Maze, error) {
	args := m.Called(userId)
	return args.Get(0).([]*Maze), args.Error(1)
}

func Test_mazeControllerImpl_CreateMaze(t *testing.T) {
	type fields struct {
		mazeRepository MazeRepository
		mazeSolver     MazeSolver
	}
	type args struct {
		userId string
		maze   *Maze
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Exit is not on the bottom",
			fields: fields{
				mazeRepository: nil,
				mazeSolver: func() MazeSolver {
					solver := &MazeSolverMock{}
					solver.On("FindSolutions", mock.Anything).Return([]MazeSolution{
						{
							Exit: "A1",
						},
					}, nil)
					return solver
				}(),
			},
			args: args{
				userId: "8Wa",
				maze: &Maze{
					EntranceX:  0,
					EntranceY:  0,
					GridWidth:  10,
					GridHeight: 10,
					Walls: []byte{
						1, 2, 3, 4, 5,
					},
				},
			},
			want:    "",
			wantErr: true,
		},

		{
			name: "Create",
			fields: fields{
				mazeRepository: func() MazeRepository {
					repo := &MazeRepositoryMock{}
					repo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("8Wa", nil)
					return repo
				}(),
				mazeSolver: func() MazeSolver {
					solver := &MazeSolverMock{}
					solver.On("FindSolutions", mock.Anything).Return([]MazeSolution{
						{
							Exit: "A8",
						},
					}, nil)
					return solver
				}(),
			},
			args: args{
				userId: "8Wa",
				maze: &Maze{
					EntranceX:  0,
					EntranceY:  0,
					GridWidth:  8,
					GridHeight: 8,
					Walls: []byte{
						68, 85, 20, 118, 18, 218, 74, 2, 0,
					},
				},
			},
			want:    "8Wa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeControllerImpl{
				mazeRepository: tt.fields.mazeRepository,
				mazeSolver:     tt.fields.mazeSolver,
			}
			got, err := m.CreateMaze(tt.args.userId, tt.args.maze)
			if (err != nil) != tt.wantErr {
				t.Errorf("mazeControllerImpl.CreateMaze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mazeControllerImpl.CreateMaze() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mazeControllerImpl_FindSolutionById(t *testing.T) {
	type fields struct {
		mazeRepository MazeRepository
		mazeSolver     MazeSolver
	}
	type args struct {
		mazeId string
		steps  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		verify  func(t *testing.T, fields *fields, got *MazeSolution)
	}{
		{
			name: "Find min solution",
			fields: fields{
				mazeRepository: func() MazeRepository {
					repo := &MazeRepositoryMock{}
					repo.On("SelectById", mock.Anything).Return(&Maze{
						EntranceX:  0,
						EntranceY:  0,
						GridWidth:  8,
						GridHeight: 8,
						Walls: []byte{
							68, 85, 20, 118, 18, 218, 74, 2, 0,
						},
					}, nil)
					return repo
				}(),
				mazeSolver: func() MazeSolver {
					solver := &MazeSolverMock{}
					solver.On("FindSolutions", mock.Anything).Return([]MazeSolution{
						{
							Length: 2,
						},
						{
							Length: 3,
						},
					}, nil)
					return solver
				}(),
			},
			args: args{
				mazeId: "8Wa",
				steps:  "min",
			},
			verify: func(t *testing.T, fields *fields, got *MazeSolution) {
				if got.Length != 2 {
					t.Errorf("mazeControllerImpl.FindSolutionById() = %v, want %v", got, 2)
				}
			},
		},
		{
			name: "Find max solution",
			fields: fields{
				mazeRepository: func() MazeRepository {
					repo := &MazeRepositoryMock{}
					repo.On("SelectById", mock.Anything).Return(&Maze{
						EntranceX:  0,
						EntranceY:  0,
						GridWidth:  8,
						GridHeight: 8,
						Walls: []byte{
							68, 85, 20, 118, 18, 218, 74, 2, 0,
						},
					}, nil)
					return repo
				}(),
				mazeSolver: func() MazeSolver {
					solver := &MazeSolverMock{}
					solver.On("FindSolutions", mock.Anything).Return([]MazeSolution{
						{
							Length: 2,
						},
						{
							Length: 3,
						},
					}, nil)
					return solver
				}(),
			},
			args: args{
				mazeId: "8Wa",
				steps:  "max",
			},
			verify: func(t *testing.T, fields *fields, got *MazeSolution) {
				if got.Length != 3 {
					t.Errorf("mazeControllerImpl.FindSolutionById() = %v, want %v", got, 2)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeControllerImpl{
				mazeRepository: tt.fields.mazeRepository,
				mazeSolver:     tt.fields.mazeSolver,
			}
			got, err := m.FindSolutionById(tt.args.mazeId, tt.args.steps)
			if (err != nil) != tt.wantErr {
				t.Errorf("mazeControllerImpl.FindSolutionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(t, &tt.fields, got)
		})
	}
}
