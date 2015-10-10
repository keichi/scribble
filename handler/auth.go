package handler

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

const (
	// Invalidate session after this period (milli seconds)
	sessionPeriod = 7 * 24 * 60 * 60 * 1000
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *registerRequest) Validate() error {
	if req.Email == "" {
		return fmt.Errorf("email is empty")
	}

	if req.Password == "" {
		return fmt.Errorf("password is empty")
	}

	return nil
}

func register(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	input := req.(*registerRequest)
	dbMap := ctx.Value("db").(*gorp.DbMap)

	count, err := dbMap.SelectInt("select count(id) from users where email = ?", input.Email)
	if err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}
	if count > 0 {
		return nil, &ErrorResponse{http.StatusBadRequest, "user exists"}
	}

	salt := auth.NewToken()

	user := model.User{
		Email:        input.Email,
		PasswordSalt: salt,
		PasswordHash: auth.HashPassword(input.Email + salt, input.Password),
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

func (req *loginRequest) Validate() error {
	if req.Email == "" {
		return fmt.Errorf("email is empty")
	}

	if req.Password == "" {
		return fmt.Errorf("password is empty")
	}

	return nil
}

func login(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	input := req.(*loginRequest)
	dbMap := ctx.Value("db").(*gorp.DbMap)
	authCtx := ctx.Value("auth").(*auth.Context)

	if authCtx.IsLoggedIn {
		return nil, &ErrorResponse{http.StatusBadRequest, "already logged in"}
	}

	var user model.User


	err := dbMap.SelectOne(&user, "select * from users where email = ?", input.Email)
	if err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, "user does not exist"}
	}

	passwordHash := auth.HashPassword(input.Email + user.PasswordSalt, input.Password)
	if passwordHash != user.PasswordHash {
		return nil, &ErrorResponse{http.StatusBadRequest, "password is wrong"}
	}

	session := model.Session{
		Token:     auth.NewToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(sessionPeriod*time.Millisecond).UnixNano() / int64(time.Millisecond),
	}

	if err := dbMap.Insert(&session); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return map[string]interface{}{
		"token":         session.Token,
		"sessionPeriod": sessionPeriod,
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

func getMyUser(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	authCtx := ctx.Value("auth").(*auth.Context)

	if !authCtx.IsLoggedIn {
		return nil, &ErrorResponse{http.StatusBadRequest, "not logged in"}
	}

	return authCtx.User, nil
}

// GetMyUser handles get current users requests
var GetMyUser = wrapJSONHandler(emptyRequest{}, getMyUser)
