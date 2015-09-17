package model

type Image struct {
	Id          int64      `db:"id" json"id,omitempty"`
	ContentType string     `db:"content_type" json:"contentType"`
	Uuid        string     `db:"uuid" json:"uuid"`
	OwnerId     int64      `db:"owner_id" json:"ownerId"`
	NoteId      int64      `db:"note_id" json:"noteId"`
	ShareState  ShareState `db:"share_state" json:"share_state"`
	CreatedAt   int64      `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   int64      `db:"updated_at" json:"updatedAt,omitempty"`
}

func (image *Image) Authorize(user *User, action AuthorizedAction) bool {
	switch action {
	case ACTION_CREATE:
		return true
	case ACTION_READ:
		return image.ShareState == SHARE_STATE_PUBLIC || image.OwnerId == user.Id
	case ACTION_UPDATE:
		return image.OwnerId == user.Id
	case ACTION_DELETE:
		return image.OwnerId == user.Id
	}

	return false
}
