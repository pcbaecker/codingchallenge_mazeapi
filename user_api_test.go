package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type UserControllerMock struct {
	mock.Mock
}

func (m *UserControllerMock) CreateUser(username string, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *UserControllerMock) Login(username string, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *UserControllerMock) GetUserForSession(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return args.String(0), args.Error(1)
}

func Test_userApiImpl_CreateUser(t *testing.T) {
	body, err := json.Marshal(&User{
		Username: "abc",
		Password: "abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
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
		verify func(t *testing.T, f *fields)
	}{
		{
			name: "Create user",
			fields: fields{
				userController: func() UserController {
					m := &UserControllerMock{}
					m.On("CreateUser", mock.Anything, mock.Anything).Return("userId", nil)
					return m
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/user", bytes.NewReader(body)),
			},
			verify: func(t *testing.T, f *fields) {
				m := f.userController.(*UserControllerMock)
				m.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userApiImpl{
				userController: tt.fields.userController,
			}
			u.CreateUser(tt.args.w, tt.args.r)
			if tt.verify != nil {
				tt.verify(t, &tt.fields)
			}
		})
	}
}

func Test_userApiImpl_Login(t *testing.T) {
	type fields struct {
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
		verify func(t *testing.T, f *fields)
	}{
		{
			name: "Login",
			fields: fields{
				userController: func() UserController {
					m := &UserControllerMock{}
					m.On("Login", mock.Anything, mock.Anything).Return("sessionId", nil)
					return m
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/login", nil),
			},
			verify: func(t *testing.T, f *fields) {
				m := f.userController.(*UserControllerMock)
				m.AssertCalled(t, "Login", mock.Anything, mock.Anything)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userApiImpl{
				userController: tt.fields.userController,
			}
			u.Login(tt.args.w, tt.args.r)
		})
	}
}
