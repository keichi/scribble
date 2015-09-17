package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/model"
)

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Session{}, "sessions").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Note{}, "notes").SetKeys(true, "Id")
	dbMap.CreateTables()

	return dbMap
}

func request(t *testing.T, url string, st int,
	req interface{}) interface{} {
	return requestWithHeader(t, url, st, req, http.Header{})
}

func requestWithHeader(t *testing.T, url string, st int,
	req interface{}, hdr http.Header) interface{} {
	assert := assert.New(t)

	bts, err := json.Marshal(req)
	assert.Nil(err, "Failed to encode request to json")

	if req == nil {
		bts = make([]byte, 0)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bts))
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

	var respJSON interface{}

	if len(body) == 0 {
		return respJSON
	}

	err = json.Unmarshal(body, &respJSON)
	assert.Nil(err, "Error while parsing response to JSON")

	return respJSON
}
