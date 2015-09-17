package model

// User holds various data on an user
type User struct {
	ID           int64  `db:"id" json:"id,omitempty"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Email        string `db:"email" json:"email"`
	CreatedAt    int64  `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt    int64  `db:"updated_at" json:"updatedAt,omitempty"`
}
