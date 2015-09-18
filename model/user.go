package model

import (
	"time"

	"gopkg.in/gorp.v1"
)

// User holds various data on an user
type User struct {
	ID           int64  `db:"id" json:"id,omitempty"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Email        string `db:"email" json:"email"`
	CreatedAt    int64  `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt    int64  `db:"updated_at" json:"updatedAt,omitempty"`
}

// PreDelete is fired before user is deleted and deletes associated notes
func (user *User) PreDelete(s gorp.SqlExecutor) error {
	_, err := s.Exec("delete from notes where owner_id = ?", user.ID)
	if err != nil {
		return err
	}

	return nil
}

// PreInsert is fired before the entity is being inserted
func (user *User) PreInsert(s gorp.SqlExecutor) error {
	user.CreatedAt = time.Now().UnixNano()
	user.UpdatedAt = user.CreatedAt
	return nil
}

// PreUpdate is fired before the entity is being updated
func (user *User) PreUpdate(s gorp.SqlExecutor) error {
	user.UpdatedAt = time.Now().UnixNano()
	return nil
}
