package model

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
