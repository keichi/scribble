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

const strechCount int = 1000
const fixedSalt string = "kRAKG5PRXZryyrPnMAXwCfGYHFfxi"

// HashPassword generates SHA256 hash of the password with salt & stretching
func HashPassword(username string, password string) string {
	pwd := []byte(password)
	salt := []byte(username + fixedSalt)
	hash := [32]byte{}

	for i := 0; i < strechCount; i++ {
		var next []byte
		next = append(next, hash[:]...)
		next = append(next, pwd...)
		next = append(next, salt...)
		hash = sha256.Sum256(next)
	}

	return fmt.Sprintf("%x", hash)
}

// NewToken creates 32byte token using CSPRNG algorithm
func NewToken() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)

	return fmt.Sprintf("%x", randBytes)
}

// AuthMiddleWare acquires current login state and user info using the
// session token stored in request header
func Middleware(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	dbMap := ctx.Value("db").(*gorp.DbMap)
	token := r.Header.Get("X-Scribble-Session")

	if token == "" {
		authCtx := &Context{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	var session model.Session
	if err := dbMap.SelectOne(&session, "select * from sessions where token = ?", token); err != nil {
		authCtx := &Context{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	var user model.User
	if err := dbMap.SelectOne(&user, "select * from users where id = ?", session.UserID); err != nil {
		authCtx := &Context{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	authCtx := &Context{true, &user, &session}
	return context.WithValue(ctx, "auth", authCtx)
}
