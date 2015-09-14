package model

type User struct {
	Id           int64  `db:"id" json:"id,omitempty"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Email        string `db:"email" json:"email"`
	CreatedAt    int64  `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt    int64  `db:"updated_at" json:"updatedAt,omitempty"`
}

type Session struct {
	Id        int64  `db:"id"`
	Token     string `db:"token"`
	UserId    int64  `db:"user_id"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

type Note struct {
	Id        int64  `db:"id" json"id,omitempty"`
	Title     string `db:"title" json:"title"`
	Content   string `db:"content" json:"content"`
	OwnerId   int64  `db:"owner_id" json:"ownerId"`
	CreatedAt int64  `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt int64  `db:"updated_at" json:"updatedAt,omitempty"`
}
