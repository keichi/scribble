package model

// Session holds the status of a session initiated by the user's login
type Session struct {
	ID        int64  `db:"id"`
	Token     string `db:"token"`
	UserID    int64  `db:"user_id"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}
