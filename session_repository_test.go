package main

import (
	"database/sql"
	"testing"
)

func newSessionTestDb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(SESSION_REPO_CREATE_TABLE)
	if err != nil {
		panic(err)
	}
	return db
}

func Test_sessionRepositoryImpl_Insert(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		sessionId string
		userId    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "session id and user id ok",
			fields: fields{
				db: newSessionTestDb(),
			},
			args: args{
				sessionId: "abc",
				userId:    "abc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sessionRepositoryImpl{
				db: tt.fields.db,
			}
			if err := s.Insert(tt.args.sessionId, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("sessionRepositoryImpl.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sessionRepositoryImpl_SelectBySessionId(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		sessionId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "session id ok",
			fields: fields{
				db: func() *sql.DB {
					db := newSessionTestDb()
					_, err := db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?, ?)", "sessionid", "userid")
					if err != nil {
						panic(err)
					}
					return db
				}(),
			},
			args: args{
				sessionId: "sessionid",
			},
			want:    "userid",
			wantErr: false,
		},
		{
			name: "session id not found",
			fields: fields{
				db: newSessionTestDb(),
			},
			args: args{
				sessionId: "sessionid",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sessionRepositoryImpl{
				db: tt.fields.db,
			}
			got, err := s.SelectBySessionId(tt.args.sessionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("sessionRepositoryImpl.SelectBySessionId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sessionRepositoryImpl.SelectBySessionId() = %v, want %v", got, tt.want)
			}
		})
	}
}
