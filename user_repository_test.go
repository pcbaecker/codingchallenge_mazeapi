package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/speps/go-hashids"
)

func newUserTestDb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(USER_REPO_CREATE_TABLE)
	if err != nil {
		panic(err)
	}
	return db
}

func newHashId() *hashids.HashID {
	hashIdData := hashids.NewData()
	hashIdData.Salt = "f948b5c5723ac19643226d73d52662d6"
	hashIdData.MinLength = 3
	hashId, err := hashids.NewWithData(hashIdData)
	if err != nil {
		panic(err)
	}
	return hashId
}

func Test_userRepositoryImpl_Insert(t *testing.T) {
	type fields struct {
		db     *sql.DB
		hashId *hashids.HashID
	}
	type args struct {
		username     string
		passwordHash string
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
				db:     newUserTestDb(),
				hashId: newHashId(),
			},
			args: args{
				username:     "abc",
				passwordHash: "abc",
			},
			wantErr: false,
			verify: func(t *testing.T, f *fields, got string) {
				if got == "" {
					t.Errorf("userRepositoryImpl.Insert() = %v, wantErr %v", got, false)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userRepositoryImpl{
				db:     tt.fields.db,
				hashid: tt.fields.hashId,
			}
			got, err := u.Insert(tt.args.username, tt.args.passwordHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepositoryImpl.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.verify != nil {
				tt.verify(t, &fields{
					db: tt.fields.db,
				}, got)
			}
		})
	}
}

func Test_userRepositoryImpl_SelectByUsername(t *testing.T) {
	type fields struct {
		db     *sql.DB
		hashid *hashids.HashID
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "username and password ok",
			fields: fields{
				db: func() *sql.DB {
					db := newUserTestDb()
					_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", "username", "password")
					if err != nil {
						panic(err)
					}
					return db
				}(),
				hashid: newHashId(),
			},
			args: args{
				username: "username",
			},
			want:    "password",
			wantErr: false,
		},
		{
			name: "username not found",
			fields: fields{
				db:     newUserTestDb(),
				hashid: newHashId(),
			},
			args: args{
				username: "username",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userRepositoryImpl{
				db:     tt.fields.db,
				hashid: tt.fields.hashid,
			}
			got, err := u.SelectByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepositoryImpl.SelectByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userRepositoryImpl.SelectByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
