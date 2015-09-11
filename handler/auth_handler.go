package handler

import (
	"golang.org/x/net/context"
	"net/http"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
	"github.com/keichi/scribble/util"
)

func Register(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")

	dbMap := ctx.Value("db").(*gorp.DbMap)

	if username == "" {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "username is empty"}, w)
		return
	}

	if password == "" {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "password is empty", }, w)
		return
	}

	count, err := dbMap.SelectInt("select count(id) from users where username = ?", username)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "user already exists"}, w)
	}

	user := model.User{
		0,
		username,
		auth.HashPassword(username, password),
		email,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	}

	if err := dbMap.Insert(&user); err != nil {
		panic(err)
	}

	util.RenderJson(http.StatusOK, map[string]string{"message": "user created"}, w)
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "Invalid request format"}, w)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.AuthContext)

	if authCtx.IsLoggedIn {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "already logged in"}, w)
		return
	}

	if username == "" {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "username is empty"}, w)
		return
	}

	if password == "" {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "password is empty", }, w)
		return
	}

	var user model.User

	passwordHash := auth.HashPassword(username, password)

	err := dbMap.SelectOne(&user, "select * from users where username = ? and password_hash = ?", username, passwordHash)
	if err != nil {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "username or password is wrong"}, w)
		return
	}

	session := model.Session{
		0,
		auth.NewToken(),
		user.Id,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	}

	if err := dbMap.Insert(&session); err != nil {
		panic(err)
	}

	util.RenderJson(http.StatusOK, map[string]string{"token":  session.Token}, w)
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.AuthContext)

	if !authCtx.IsLoggedIn {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "not logged in"}, w)
		return
	}

	count, err := dbMap.Delete(authCtx.Session)
	if count <= 0 || err != nil {
		util.RenderJson(http.StatusBadRequest, map[string]string{"message": "log out failed"}, w)
		return
	}

	util.RenderJson(http.StatusOK, map[string]string{"message": "logged out"}, w)
}
