package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type MazeControllerMock struct {
	mock.Mock
}

func (m *MazeControllerMock) GetUserMazes(userId string) ([]*Maze, error) {
	args := m.Called(userId)
	return args.Get(0).([]*Maze), args.Error(1)
}

func (m *MazeControllerMock) Generate(width uint16, height uint16) (*Maze, error) {
	args := m.Called(width, height)
	return args.Get(0).(*Maze), args.Error(1)
}
func (m *MazeControllerMock) CreateMaze(userId string, maze *Maze) (string, error) {
	args := m.Called(maze)
	return args.String(0), args.Error(1)
}
func (m *MazeControllerMock) DrawMaze(maze *Maze, pathItem *PathItem, w io.Writer) error {
	args := m.Called(maze, pathItem, w)
	return args.Error(0)
}

func (m *MazeControllerMock) FindSolutionById(mazeId string, steps string) (*MazeSolution, error) {
	args := m.Called(mazeId, steps)
	return args.Get(0).(*MazeSolution), args.Error(1)
}

func Test_mazeApiImpl_CreateMaze(t *testing.T) {
	type fields struct {
		mazeController MazeController
		userController UserController
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		verify func(t *testing.T, f *fields, w *httptest.ResponseRecorder)
	}{
		{
			name: "Create maze",
			fields: fields{
				mazeController: func() MazeController {
					m := &MazeControllerMock{}
					m.On("CreateMaze", mock.Anything, mock.Anything).Return("aaa", nil)
					return m
				}(),
				userController: func() UserController {
					m := &UserControllerMock{}
					m.On("GetUserForSession", mock.Anything).Return("aaa", nil)
					return m
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest("POST", "/maze", serializeMazeApiDao(&MazeApiDao{
						Entrace:  "A1",
						GridSize: "10x10",
						Walls:    []string{},
					}))
					req.Header.Add("Authorization", "Bearer bbaaaaab")
					return req
				}(),
			},
			verify: func(t *testing.T, f *fields, w *httptest.ResponseRecorder) {
				if w.Code != http.StatusCreated {
					t.Errorf("CreateMaze() status code = %v, want %v", w.Code, http.StatusOK)
				}
				if w.Body.String() != `{"mazeId": "aaa"}` {
					t.Errorf("CreateMaze() body = %v, want %v", w.Body.String(), `{"mazeId": "aaa"}`)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeApiImpl{
				mazeController: tt.fields.mazeController,
				userController: tt.fields.userController,
			}
			m.CreateMaze(tt.args.w, tt.args.r)
			tt.verify(t, &tt.fields, tt.args.w.(*httptest.ResponseRecorder))
		})
	}
}

func serializeMazeApiDao(input *MazeApiDao) io.Reader {
	result := new(bytes.Buffer)
	err := json.NewEncoder(result).Encode(input)
	if err != nil {
		panic(err)
	}
	return result
}

func Test_getUserId(t *testing.T) {
	type args struct {
		r              *http.Request
		userController UserController
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Find userid",
			args: args{
				r: func() *http.Request {
					request := httptest.NewRequest("GET", "/maze", nil)
					request.Header.Add("Authorization", "Bearer asdfghjkl123456789")
					return request
				}(),
				userController: func() UserController {
					m := &UserControllerMock{}
					m.On("GetUserForSession", mock.Anything).Return("aaa", nil)
					return m
				}(),
			},
			want:    "aaa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserId(tt.args.r, tt.args.userController)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mazeApiImpl_FindSolutionById(t *testing.T) {
	type fields struct {
		mazeController MazeController
		userController UserController
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		verify func(t *testing.T, f *fields, w *httptest.ResponseRecorder)
	}{
		{
			name: "Find solution",
			fields: fields{
				mazeController: func() MazeController {
					m := &MazeControllerMock{}
					m.On("FindSolutionById", mock.Anything, mock.Anything).Return(&MazeSolution{
						Length: 10,
						Path:   []string{"A1", "A2"},
						Exit:   "A2",
					}, nil)
					return m
				}(),
				userController: func() UserController {
					m := &UserControllerMock{}
					m.On("GetUserForSession", mock.Anything).Return("aaa", nil)
					return m
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest("GET", "/maze/10/solution?steps=min", nil)
					req.Header.Add("Authorization", "Bearer bbaaaaab")
					req = mux.SetURLVars(req, map[string]string{
						"mazeId": "abcd",
					})
					return req
				}(),
			},
			verify: func(t *testing.T, f *fields, w *httptest.ResponseRecorder) {
				if w.Code != http.StatusOK {
					t.Errorf("FindSolutionById() status code = %v, want %v", w.Code, http.StatusOK)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeApiImpl{
				mazeController: tt.fields.mazeController,
				userController: tt.fields.userController,
			}
			m.FindSolutionById(tt.args.w, tt.args.r)
			tt.verify(t, &tt.fields, tt.args.w.(*httptest.ResponseRecorder))
		})
	}
}
