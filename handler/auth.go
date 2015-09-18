package handler

import (
	"golang.org/x/net/context"
	"net/http"

	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
	"time"
)

const (
	// Invalidate session after this period (milli seconds)
	sessionPeriod = 7 * 24 * 60 * 60 * 1000
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func register(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	input := req.(*registerRequest)
	dbMap := ctx.Value("db").(*gorp.DbMap)

	if input.Email == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "email is empty"}
	}

	if input.Password == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "password is empty"}
	}

	count, err := dbMap.SelectInt("select count(id) from users where email = ?", input.Email)
	if err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}
	if count > 0 {
		return nil, &ErrorResponse{http.StatusBadRequest, "user exists"}
	}

	user := model.User{
		Email:        input.Email,
		PasswordHash: auth.HashPassword(input.Email, input.Password),
	}

	if err := dbMap.Insert(&user); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return map[string]string{"message": "user created"}, nil
}

// Register handles user registration requests
var Register = wrapJSONHandler(registerRequest{}, register)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	input := req.(*loginRequest)
	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.Context)

	if authCtx.IsLoggedIn {
		return nil, &ErrorResponse{http.StatusBadRequest, "already logged in"}
	}

	if input.Email == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "email is empty"}
	}

	if input.Password == "" {
		return nil, &ErrorResponse{http.StatusBadRequest, "password is empty"}
	}

	var user model.User

	passwordHash := auth.HashPassword(input.Email, input.Password)

	err := dbMap.SelectOne(&user, "select * from users where email = ? and password_hash = ?", input.Email, passwordHash)
	if err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, "email or password is wrong"}
	}

	session := model.Session{
		Token:     auth.NewToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(sessionPeriod * time.Millisecond).UnixNano(),
	}

	if err := dbMap.Insert(&session); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return map[string]string{
		"token":          session.Token,
		"session_period": string(sessionPeriod),
	}, nil
}

// Login handles user login requests
var Login = wrapJSONHandler(loginRequest{}, login)

func logout(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.Context)

	if !authCtx.IsLoggedIn {
		return nil, &ErrorResponse{http.StatusBadRequest, "not logged in"}
	}

	count, err := dbMap.Delete(authCtx.Session)
	if count <= 0 || err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, "log out failed"}
	}

	return map[string]string{"message": "logged out"}, nil
}

// Logout handles user logout requests
var Logout = wrapJSONHandler(emptyRequest{}, logout)
