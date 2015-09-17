package model

// Image is an image attached to a note
type Image struct {
	ID          int64      `db:"id" json"id,omitempty"`
	ContentType string     `db:"content_type" json:"contentType"`
	UUID        string     `db:"uuid" json:"uuid"`
	OwnerID     int64      `db:"owner_id" json:"ownerId"`
	NoteID      int64      `db:"note_id" json:"noteId"`
	ShareState  ShareState `db:"share_state" json:"share_state"`
	CreatedAt   int64      `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   int64      `db:"updated_at" json:"updatedAt,omitempty"`
}

// Authorize determines what kind of action are allowed on this image
func (image *Image) Authorize(user *User, action AuthorizedAction) bool {
	switch action {
	case ActionCreate:
		return true
	case ActionRead:
		return image.ShareState == ShareStatePublic || image.OwnerID == user.ID
	case ActionUpdate:
		return image.OwnerID == user.ID
	case ActionDelete:
		return image.OwnerID == user.ID
	}

	return false
}
