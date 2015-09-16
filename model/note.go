package model

type ShareState int8

const (
	SHARE_STATE_PUBLIC int = iota
	SHARE_STATE_PRIVATE
)

type Note struct {
	Id         int64      `db:"id" json"id,omitempty"`
	Title      string     `db:"title" json:"title"`
	Content    string     `db:"content" json:"content"`
	OwnerId    int64      `db:"owner_id" json:"ownerId"`
	ShareState ShareState `db:"share_state" json:"share_state"`
	CreatedAt  int64      `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt  int64      `db:"updated_at" json:"updatedAt,omitempty"`
}

func (note *Note) Authorize(user *User) {
}
