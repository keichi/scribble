package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)

	handlerFunc := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Register(ctx, w, r)
		},
	)

	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	resp := request(t, server.URL, http.StatusOK,
		map[string]string{
			"email":    "test@example.com",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "user created"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "email is empty"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "test@example.com",
			"password": "",
		},
	)
	assert.Equal(map[string]interface{}{"message": "password is empty"}, resp)

	count, err := dbMap.SelectInt("select count(id) from users")
	assert.Nil(err, "Error while querying user count")
	assert.EqualValues(1, count, "Wrong user count")
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()

	authCtx := &auth.Context{false, &model.User{}, &model.Session{}}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", authCtx)

	handlerFunc := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Login(ctx, w, r)
		},
	)

	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	resp := request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "email is empty"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "test@example.com",
			"password": "",
		},
	)
	assert.Equal(map[string]interface{}{"message": "password is empty"}, resp)

	authCtx.IsLoggedIn = true
	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "test@example.com",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "already logged in"}, resp)
	authCtx.IsLoggedIn = false

	salt := auth.NewToken()
	user := model.User{
		Email:        "test@example.com",
		PasswordSalt: salt,
		PasswordHash: auth.HashPassword("test@example.com" + salt, "testpassword"),
	}

	err := dbMap.Insert(&user)
	assert.Nil(err, "Failed to insert test user")

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "unknown-user@example.com",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "user does not exist"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "password is wrong"}, resp)

	resp = request(t, server.URL, http.StatusOK,
		map[string]string{
			"email":    "test@example.com",
			"password": "testpassword",
		},
	)
	token := resp.(map[string]interface{})["token"]
	assert.True(token != "", "Session token must not be empty")
}

func TestLogout(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()

	authCtx := &auth.Context{false, &model.User{}, &model.Session{}}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", authCtx)

	handlerFunc := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Logout(ctx, w, r)
		},
	)

	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	resp := request(t, server.URL, http.StatusBadRequest,
		map[string][]string{},
	)
	assert.Equal(map[string]interface{}{"message": "not logged in"}, resp)

	authCtx.IsLoggedIn = true
	const testSessionToken string = "3a11779677f844f581448ba6337225499dae0850c26665a83ae344609157774"

	session := model.Session{
		Token:  testSessionToken,
		UserID: 0,
	}

	err := dbMap.Insert(&session)
	assert.Nil(err, "Failed to insert test session")

	header := http.Header{"X-Session-Token": []string{testSessionToken}}
	requestWithHeader(t, server.URL, http.StatusBadRequest, map[string]string{}, header)
}
