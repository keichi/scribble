package model

import (
	"time"

	"gopkg.in/gorp.v1"
)

// ShareState represents the scope of disclosure of an entity
type ShareState int8

const (
	// ShareStatePublic means the entity is public
	ShareStatePublic ShareState = iota
	// ShareStatePrivate means the entity is private
	ShareStatePrivate
)

// Note is a single markdown note
type Note struct {
	ID         int64      `db:"id" json:"id,omitempty"`
	Title      string     `db:"title" json:"title"`
	Content    string     `db:"content" json:"content"`
	OwnerID    int64      `db:"owner_id" json:"ownerId"`
	ShareState ShareState `db:"share_state" json:"share_state"`
	Images     []Image    `db:"-" json:"images"`
	CreatedAt  int64      `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt  int64      `db:"updated_at" json:"updatedAt,omitempty"`
}

// Authorize determines what kind of action are allowed on this note
func (note *Note) Authorize(user *User, action AuthorizedAction) bool {
	switch action {
	case ActionCreate:
		return true
	case ActionRead:
		return note.ShareState == ShareStatePublic || note.OwnerID == user.ID
	case ActionUpdate:
		return note.OwnerID == user.ID
	case ActionDelete:
		return note.OwnerID == user.ID
	}

	return false
}

// PostGet is fired after the entity is acquired from the database
func (note *Note) PostGet(s gorp.SqlExecutor) error {
	var images []Image
	_, err := s.Select(&images, "select * from images where note_id = ?", note.ID)
	if err != nil {
		return err
	}

	note.Images = images
	return nil
}

// PreInsert is fired before the entity is being inserted
func (note *Note) PreInsert(s gorp.SqlExecutor) error {
	note.CreatedAt = time.Now().UnixNano()
	note.UpdatedAt = note.CreatedAt
	return nil
}

// PreUpdate is fired before the entity is being updated
func (note *Note) PreUpdate(s gorp.SqlExecutor) error {
	note.UpdatedAt = time.Now().UnixNano()
	return nil
}
