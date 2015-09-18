package model

// Image is an image attached to a note
type Image struct {
	ID          int64      `db:"id" json:"id,omitempty"`
	ContentType string     `db:"content_type" json:"contentType"`
	UUID        string     `db:"uuid" json:"uuid"`
	NoteID      int64      `db:"note_id" json:"noteId,omitempty"`
	Note        *Note      `db:"-" json:"-"`
	CreatedAt   int64      `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   int64      `db:"updated_at" json:"updatedAt,omitempty"`
}

// Authorize determines what kind of action are allowed on this image
func (image *Image) Authorize(user *User, action AuthorizedAction) bool {
	switch action {
	case ActionCreate:
		return image.Note.OwnerID == user.ID
	case ActionRead:
		return image.Note.ShareState == ShareStatePublic || image.Note.OwnerID == user.ID
	case ActionUpdate:
		return image.Note.OwnerID == user.ID
	case ActionDelete:
		return image.Note.OwnerID == user.ID
	}

	return false
}
