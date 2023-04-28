package main

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/speps/go-hashids"
	"github.com/stretchr/testify/assert"
)

func newMazeTestDb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(MAZE_REPO_CREATE_TABLE)
	if err != nil {
		panic(err)
	}
	return db
}

func Test_mazeRepositoryImpl_Insert(t *testing.T) {
	type fields struct {
		db     *sql.DB
		hashid *hashids.HashID
	}
	type args struct {
		userId     string
		entranceX  uint16
		entranceY  uint16
		gridWidth  uint16
		gridHeight uint16
		walls      []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		verify  func(t *testing.T, fields *fields)
	}{
		{
			name: "Insert",
			fields: fields{
				db:     newMazeTestDb(),
				hashid: newHashId(),
			},
			args: args{
				userId:     "8Wa",
				entranceX:  0,
				entranceY:  0,
				gridWidth:  10,
				gridHeight: 10,
				walls: []byte{
					1, 2, 3, 4, 5,
				},
			},
			wantErr: false,
			want:    "8Wa",
			verify: func(t *testing.T, fields *fields) {
				m := &mazeRepositoryImpl{
					db:     fields.db,
					hashid: fields.hashid,
				}
				count, err := m.CountAll()
				assert.Nil(t, err)
				assert.Equal(t, uint64(1), count)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeRepositoryImpl{
				db:     tt.fields.db,
				hashid: tt.fields.hashid,
			}
			got, err := m.Insert(tt.args.userId, tt.args.entranceX, tt.args.entranceY, tt.args.gridWidth, tt.args.gridHeight, tt.args.walls)
			if (err != nil) != tt.wantErr {
				t.Errorf("mazeRepositoryImpl.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mazeRepositoryImpl.Insert() = %v, want %v", got, tt.want)
			}
			if tt.verify != nil {
				tt.verify(t, &tt.fields)
			}
		})
	}
}

func Test_mazeRepositoryImpl_SelectById(t *testing.T) {
	type fields struct {
		db     *sql.DB
		hashid *hashids.HashID
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Maze
		wantErr bool
	}{
		{
			name: "SelectById",
			fields: fields{
				db: func() *sql.DB {
					db := newMazeTestDb()
					m := &mazeRepositoryImpl{
						db:     db,
						hashid: newHashId(),
					}
					m.Insert("8Wa", 3, 4, 10, 10, []byte{1, 2, 3, 4, 5})
					return db
				}(),
				hashid: newHashId(),
			},
			args: args{
				id: "8Wa",
			},
			want: &Maze{
				Id:         "8Wa",
				EntranceX:  3,
				EntranceY:  4,
				GridWidth:  10,
				GridHeight: 10,
				Walls: []byte{
					1, 2, 3, 4, 5,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeRepositoryImpl{
				db:     tt.fields.db,
				hashid: tt.fields.hashid,
			}
			got, err := m.SelectById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("mazeRepositoryImpl.SelectById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mazeRepositoryImpl.SelectById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mazeRepositoryImpl_SelectAllByUserId(t *testing.T) {
	userId := "username"
	type fields struct {
		db     *sql.DB
		hashid *hashids.HashID
	}
	type args struct {
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		verify  func(t *testing.T, fields *fields, got []*Maze)
	}{
		{
			name: "SelectAllByUserId",
			fields: fields{
				db: func() *sql.DB {
					db := newMazeTestDb()
					_, err := db.Exec(USER_REPO_CREATE_TABLE)
					if err != nil {
						panic(err)
					}
					u := &userRepositoryImpl{
						db:     db,
						hashid: newHashId(),
					}
					userId, err = u.Insert("username", "password")
					if err != nil {
						panic(err)
					}
					m := &mazeRepositoryImpl{
						db:     db,
						hashid: newHashId(),
					}
					_, err = m.Insert("username", 3, 4, 10, 10, []byte{1, 2, 3, 4, 5})
					if err != nil {
						panic(err)
					}
					_, err = m.Insert("username", 4, 5, 11, 11, []byte{1, 2, 3, 4, 5})
					if err != nil {
						panic(err)
					}
					return db
				}(),
				hashid: newHashId(),
			},
			args: args{
				userId: userId,
			},
			wantErr: false,
			verify: func(t *testing.T, fields *fields, got []*Maze) {
				assert.Equal(t, 2, len(got))
				assert.Equal(t, "8Wa", got[0].Id)
				assert.Equal(t, uint16(3), got[0].EntranceX)
				assert.Equal(t, uint16(4), got[0].EntranceY)
				assert.Equal(t, uint16(10), got[0].GridWidth)
				assert.Equal(t, uint16(10), got[0].GridHeight)
				assert.Equal(t, []byte{1, 2, 3, 4, 5}, got[0].Walls)
				assert.Equal(t, "E42", got[1].Id)
				assert.Equal(t, uint16(4), got[1].EntranceX)
				assert.Equal(t, uint16(5), got[1].EntranceY)
				assert.Equal(t, uint16(11), got[1].GridWidth)
				assert.Equal(t, uint16(11), got[1].GridHeight)
				assert.Equal(t, []byte{1, 2, 3, 4, 5}, got[1].Walls)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeRepositoryImpl{
				db:     tt.fields.db,
				hashid: tt.fields.hashid,
			}
			got, err := m.SelectAllByUserId(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("mazeRepositoryImpl.SelectAllByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.verify != nil {
				tt.verify(t, &tt.fields, got)
			}
		})
	}
}
