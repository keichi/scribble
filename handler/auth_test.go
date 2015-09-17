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
			"username": "testuser",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "user created"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"username": "",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "username is empty"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"username": "testuser",
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
			"username": "",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "username is empty"}, resp)

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"username": "testuser",
			"password": "",
		},
	)
	assert.Equal(map[string]interface{}{"message": "password is empty"}, resp)

	authCtx.IsLoggedIn = true
	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"username": "testuser",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "already logged in"}, resp)
	authCtx.IsLoggedIn = false

	user := model.User{
		0,
		"testuser",
		auth.HashPassword("testuser", "testpassword"),
		"test@example.com",
		1441872075622000,
		1441872075622000,
	}

	err := dbMap.Insert(&user)
	assert.Nil(err, "Failed to insert test user")

	resp = request(t, server.URL, http.StatusBadRequest,
		map[string]string{
			"username": "test",
			"password": "test",
		},
	)
	assert.Equal(map[string]interface{}{"message": "username or password is wrong"}, resp)

	resp = request(t, server.URL, http.StatusOK,
		map[string]string{
			"username": "testuser",
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
		0,
		testSessionToken,
		0,
		1441872075622000,
		1441872075622000,
	}

	err := dbMap.Insert(&session)
	assert.Nil(err, "Failed to insert test session")

	header := http.Header{"X-Session-Token": []string{testSessionToken}}
	requestWithHeader(t, server.URL, http.StatusBadRequest, map[string]string{}, header)
}
