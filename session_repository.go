package main

import "database/sql"

const SESSION_REPO_CREATE_TABLE = "CREATE TABLE IF NOT EXISTS sessions (session_id CHAR(32) NOT NULL PRIMARY KEY, user_id VARCHAR(256))"

type SessionRepository interface {
	Insert(sessionId string, userId string) error
	SelectBySessionId(sessionId string) (string, error)
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	_, err := db.Exec(SESSION_REPO_CREATE_TABLE)
	if err != nil {
		panic(err)
	}
	return &sessionRepositoryImpl{
		db: db,
	}
}

type sessionRepositoryImpl struct {
	db *sql.DB
}

func (s *sessionRepositoryImpl) Insert(sessionId string, userId string) error {
	_, err := s.db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?, ?)", sessionId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRepositoryImpl) SelectBySessionId(sessionId string) (string, error) {
	var userId string
	err := s.db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionId).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return userId, nil
}
