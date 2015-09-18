package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/goamz/goamz/s3/s3test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/model"
)

func initDB() *gorp.DbMap {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(model.Session{}, "sessions").SetKeys(true, "ID")
	dbMap.AddTableWithName(model.Note{}, "notes").SetKeys(true, "ID")
	dbMap.AddTableWithName(model.Image{}, "images").SetKeys(true, "ID")
	dbMap.CreateTables()

	return dbMap
}

func initS3() *s3.Bucket {
	srv, err := s3test.NewServer(&s3test.Config{})
	if err != nil {
		panic(err)
	}
	region := aws.Region{
		Name:                 "dummy-region-1",
		S3Endpoint:           srv.URL(),
		S3LocationConstraint: true,
	}

	s3 := s3.New(aws.Auth{}, region)
	bucket := s3.Bucket("scribble-image-store")

	return bucket
}

func request(t *testing.T, url string, st int,
	req interface{}) interface{} {
	return requestWithHeader(t, url, st, req, http.Header{})
}

func requestWithHeader(t *testing.T, url string, st int,
	req interface{}, hdr http.Header) interface{} {
	assert := assert.New(t)

	bts, err := json.Marshal(req)
	assert.Nil(err, "Failed to encode request to JSON: %v", err)

	if req == nil {
		bts = make([]byte, 0)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bts))
	assert.Nil(err, "Failed to create request: %v", err)

	request.Header = hdr
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", strconv.Itoa(len(bts)))

	response, err := http.DefaultClient.Do(request)
	assert.Nil(err, "Failed to execute request: %v", err)

	body, err := ioutil.ReadAll(response.Body)
	assert.Nil(err, "Error while reading resp body: %v", err)
	defer response.Body.Close()

	assert.Equal(st, response.StatusCode, "Wrong status code")

	var respJSON interface{}

	if len(body) == 0 {
		return respJSON
	}

	err = json.Unmarshal(body, &respJSON)
	assert.Nil(err, "Error while parsing response to JSON: %v", err)

	return respJSON
}
