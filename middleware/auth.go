package middleware

import (
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

// AuthMiddleWare acquires current login state and user info using the
// session token stored in request header
func Auth(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	dbMap := ctx.Value("db").(*gorp.DbMap)
	token := r.Header.Get("X-Scribble-Session")

	if token == "" {
		authCtx := &auth.Context{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	var session model.Session
	if err := dbMap.SelectOne(&session, "select * from sessions where token = ?", token); err != nil {
		authCtx := &auth.Context{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	var user model.User
	if err := dbMap.SelectOne(&user, "select * from users where id = ?", session.UserID); err != nil {
		authCtx := &auth.Context{false, &model.User{}, &model.Session{}}
		return context.WithValue(ctx, "auth", authCtx)
	}

	authCtx := &auth.Context{true, &user, &session}
	return context.WithValue(ctx, "auth", authCtx)
}
