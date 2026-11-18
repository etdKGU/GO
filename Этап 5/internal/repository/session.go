package repository

import (
	"database/sql"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(sessionID string, userID int, expires time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expires.UTC().Format(time.RFC3339),
	)
	return err
}

// Возвращает id пользователя, если сессия существует и срок ещё не истёк
func (r *SessionRepository) GetUserID(sessionID string) (int, error) {
	var userID int
	err := r.db.QueryRow(
		"SELECT user_id FROM sessions WHERE session_id = ? AND expires_at > ?",
		sessionID, time.Now().UTC().Format(time.RFC3339),
	).Scan(&userID)
	return userID, err
}

func (r *SessionRepository) Delete(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}
