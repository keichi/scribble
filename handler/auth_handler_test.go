package handler

import (
	"flag"
	"os"
	"net/http"
	"net/http/httptest"
	"database/sql"
	"testing"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"
	"github.com/stretchr/testify/assert"

	"github.com/keichi/scribble/model"
	"github.com/keichi/scribble/auth"
)

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Session{}, "sessions").SetKeys(true, "Id")
	dbMap.CreateTables()

	return dbMap
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func request(t *testing.T, sv *httptest.Server, st int,
					req interface{}) map[string]interface{} {
	return requestWithHeader(t, sv, st, req, http.Header{})
}

func requestWithHeader(t *testing.T, sv *httptest.Server, st int,
					req  interface{}, hdr http.Header) map[string]interface{} {
	assert := assert.New(t)

	bts, err := json.Marshal(req)
	assert.Nil(err, "Failed to encode request to json")

	request, err := http.NewRequest("POST", sv.URL, bytes.NewBuffer(bts))
	assert.Nil(err, "Failed to create request")

	request.Header = hdr
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", strconv.Itoa(len(bts)))

	response, err := http.DefaultClient.Do(request)
	assert.Nil(err, "Failed to do request")

	body, err := ioutil.ReadAll(response.Body)
	assert.Nil(err, "Error while reading resp body")
	defer response.Body.Close()

	assert.Equal(st, response.StatusCode, "Wrong status code")

	respJson := make(map[string]interface{})
	err = json.Unmarshal(body, &respJson)
	assert.Nil(err, "Error while parsing response to JSON")

	return respJson
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)

	handlerFunc := http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			Register(ctx, w, r)
		},
	)

	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	resp := request(t, server, http.StatusOK,
		map[string]string {
			"username": "testuser",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "user created"}, resp)

	resp = request(t, server, http.StatusBadRequest,
		map[string]string {
			"username": "",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "username is empty"}, resp)

	resp = request(t, server, http.StatusBadRequest,
		map[string]string {
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

	dbMap := initDb()
	defer dbMap.Db.Close()

	authCtx := &auth.AuthContext{false, &model.User{}, &model.Session{}}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", authCtx)

	handlerFunc := http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			Login(ctx, w, r)
		},
	)

	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	resp := request(t, server, http.StatusBadRequest,
		map[string]string {
			"username": "",
			"password": "testpassword",
		},
	)
	assert.Equal(map[string]interface{}{"message": "username is empty"}, resp)

	resp = request(t, server, http.StatusBadRequest,
		map[string]string {
			"username": "testuser",
			"password": "",
		},
	)
	assert.Equal(map[string]interface{}{"message": "password is empty"}, resp)

	authCtx.IsLoggedIn = true
	resp = request(t, server, http.StatusBadRequest,
		map[string]string {
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

	resp = request(t, server, http.StatusBadRequest,
		map[string]string {
			"username": "test",
			"password": "test",
		},
	)
	assert.Equal(map[string]interface{}{"message": "username or password is wrong"}, resp)

	resp = request(t, server, http.StatusOK,
		map[string]string {
			"username": "testuser",
			"password": "testpassword",
		},
	)
	assert.True(resp["token"] != "", "Session token must not be empty")
}

func TestLogout(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()

	authCtx := &auth.AuthContext{false, &model.User{}, &model.Session{}}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", authCtx)

	handlerFunc := http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			Logout(ctx, w, r)
		},
	)

	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	resp := request(t, server, http.StatusBadRequest,
		map[string][]string {},
	)
	assert.Equal(map[string]interface{}{"message": "not logged in"}, resp)

	authCtx.IsLoggedIn = true
	const testSessionToken string =
	"3a11779677f844f581448ba6337225499dae0850c26665a83ae344609157774"

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
	requestWithHeader(t, server, http.StatusBadRequest, map[string]string{}, header)
}
