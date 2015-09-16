package model

type Session struct {
	Id        int64  `db:"id"`
	Token     string `db:"token"`
	UserId    int64  `db:"user_id"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}
