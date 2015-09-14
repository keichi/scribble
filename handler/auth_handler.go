package handler

import (
	"golang.org/x/net/context"
	"net/http"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func register(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	input := req.(*registerRequest)
	dbMap := ctx.Value("db").(*gorp.DbMap)

	if input.Username == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "username is empty"}
	}

	if input.Password == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "password is empty"}
	}

	count, err := dbMap.SelectInt("select count(id) from users where username = ?", input.Username)
	if err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}
	if count > 0 {
		return nil, &ErrorResponse{http.StatusBadRequest, "user exists"}
	}

	user := model.User{
		0,
		input.Username,
		auth.HashPassword(input.Username, input.Password),
		input.Email,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	}

	if err := dbMap.Insert(&user); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return map[string]string{"message": "user created"}, nil
}

var Register = WrapJsonHandler(registerRequest{}, register)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	input := req.(*loginRequest)
	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.AuthContext)

	if authCtx.IsLoggedIn {
		return nil, &ErrorResponse{http.StatusBadRequest, "already logged in"}
	}

	if input.Username == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "username is empty"}
	}

	if input.Password == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "password is empty"}
	}

	var user model.User

	passwordHash := auth.HashPassword(input.Username, input.Password)

	err := dbMap.SelectOne(&user, "select * from users where username = ? and password_hash = ?", input.Username, passwordHash)
	if err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, "username or password is wrong"}
	}

	session := model.Session{
		0,
		auth.NewToken(),
		user.Id,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	}

	if err := dbMap.Insert(&session); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return map[string]string{"token": session.Token}, nil
}

var Login = WrapJsonHandler(loginRequest{}, login)

func logout(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.AuthContext)

	if !authCtx.IsLoggedIn {
		return nil, &ErrorResponse{http.StatusBadRequest, "not logged in"}
	}

	count, err := dbMap.Delete(authCtx.Session)
	if count <= 0 || err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, "log out failed"}
	}

	return map[string]string{"message": "logged out"}, nil
}

var Logout = WrapJsonHandler(nil, logout)
