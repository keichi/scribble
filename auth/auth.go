package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/model"
)

const STRETCH_COUNT int = 1000
const FIXED_SALT string = "kRAKG5PRXZryyrPnMAXwCfGYHFfxi"

type AuthContext struct {
	IsLoggedIn bool
	User       *model.User
	Session    *model.Session
}

func HashPassword(username string, password string) string {
	pwd := []byte(password)
	salt := []byte(username + FIXED_SALT)
	hash := [32]byte{}

	for i := 0; i < STRETCH_COUNT; i++ {
		next := make([]byte, 0)
		next = append(next, hash[:]...)
		next = append(next, pwd...)
		next = append(next, salt...)
		hash = sha256.Sum256(next)
	}

	return fmt.Sprintf("%x", hash)
}

func NewToken() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)

	return fmt.Sprintf("%x", randBytes)
}

func AuthMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	dbMap := ctx.Value("db").(*gorp.DbMap)
	token := r.Header.Get("X-Session-Token")

	if token == "" {
		authCtx := &AuthContext{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	var session model.Session
	if err := dbMap.SelectOne(&session, "select * from sessions where token = ?", token); err != nil {
		authCtx := &AuthContext{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	var user model.User
	if err := dbMap.SelectOne(&user, "select * from users where id = ?", session.UserId); err != nil {
		authCtx := &AuthContext{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	authCtx := &AuthContext{true, &user, &session}
	return context.WithValue(ctx, "auth", authCtx)
}
