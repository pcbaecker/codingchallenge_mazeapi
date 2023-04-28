package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SessionRepositoryMock struct {
	mock.Mock
}

func (m *SessionRepositoryMock) Insert(sessionId string, userId string) error {
	args := m.Called(sessionId, userId)
	return args.Error(0)
}

func (m *SessionRepositoryMock) SelectBySessionId(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return args.String(0), args.Error(1)
}

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Insert(username string, passwordHash string) (string, error) {
	args := m.Called(username, passwordHash)
	return args.String(0), args.Error(1)
}

func (m *UserRepositoryMock) SelectByUsername(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func Test_userControllerImpl_CreateUser(t *testing.T) {
	type fields struct {
		userRepository UserRepository
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		verify  func(t *testing.T, f *fields, got string)
	}{
		{
			name: "username and password too short",
			fields: fields{
				userRepository: &UserRepositoryMock{},
			},
			args: args{
				username: "a",
				password: "a",
			},
			wantErr: true,
			verify: func(t *testing.T, f *fields, got string) {
				m := f.userRepository.(*UserRepositoryMock)
				m.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
			},
		},
		{
			name: "username and password ok",
			fields: fields{
				userRepository: func() UserRepository {
					m := &UserRepositoryMock{}
					m.On("Insert", mock.Anything, mock.Anything).Return("userid", nil)
					return m
				}(),
			},
			args: args{
				username: "abc",
				password: "abc",
			},
			wantErr: false,
			verify: func(t *testing.T, f *fields, got string) {
				assert.Equal(t, "userid", got)
				m := f.userRepository.(*UserRepositoryMock)
				m.AssertCalled(t, "Insert", mock.Anything, mock.Anything)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userControllerImpl{
				userRepository: tt.fields.userRepository,
			}
			got, err := u.CreateUser(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userControllerImpl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.verify != nil {
				tt.verify(t, &fields{
					userRepository: tt.fields.userRepository,
				}, got)
			}
		})
	}
}

func Test_userControllerImpl_Login(t *testing.T) {
	password := "abc"
	passwordHash := "$argon2id$v=19$m=65536,t=1,p=4$lizu6Pb8PTek6GrGM84e1Q$s8RyFZLrr3tMYQYUV+resj0pLfdVmm3OHSgmH8PrZGI"
	type fields struct {
		userRepository    UserRepository
		sessionRepository SessionRepository
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		verify  func(t *testing.T, f *fields, got string)
	}{
		{
			name: "username and password ok",
			fields: fields{
				userRepository: func() UserRepository {
					m := &UserRepositoryMock{}
					m.On("SelectByUsername", mock.Anything).Return(passwordHash, nil)
					return m
				}(),
				sessionRepository: func() SessionRepository {
					m := &SessionRepositoryMock{}
					m.On("Insert", mock.Anything, mock.Anything).Return(nil)
					return m
				}(),
			},
			args: args{
				username: "abc",
				password: password,
			},
			wantErr: false,
			verify: func(t *testing.T, f *fields, got string) {
				assert.NotEmpty(t, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userControllerImpl{
				userRepository:    tt.fields.userRepository,
				sessionRepository: tt.fields.sessionRepository,
			}
			got, err := u.Login(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userControllerImpl.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.verify != nil {
				tt.verify(t, &tt.fields, got)
			}
		})
	}
}
