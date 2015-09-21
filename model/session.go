package model

import (
	"time"

	"gopkg.in/gorp.v1"
)

// Session holds the status of a session initiated by the user's login
type Session struct {
	ID        int64  `db:"id"`
	Token     string `db:"token"`
	UserID    int64  `db:"user_id"`
	ExpiresAt int64  `db:"expires_at"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

// PreInsert is fired before the entity is being inserted
func (session *Session) PreInsert(s gorp.SqlExecutor) error {
	session.CreatedAt = time.Now().UnixNano() / int64(time.Millisecond)
	session.UpdatedAt = session.CreatedAt
	return nil
}

// PreUpdate is fired before the entity is being updated
func (session *Session) PreUpdate(s gorp.SqlExecutor) error {
	session.UpdatedAt = time.Now().UnixNano() / int64(time.Millisecond)
	return nil
}
