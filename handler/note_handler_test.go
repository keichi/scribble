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

func TestListNotes(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)

	dbMap.Insert(&model.Note{
		Id:        0,
		Title:     "Test Title 1",
		Content:   "lorem ipsum dolor sit amet consetetur.",
		OwnerId:   0,
		CreatedAt: 1442284669000,
		UpdatedAt: 1442284669000,
	}, &model.Note{
		Id:        0,
		Title:     "Test Title 2",
		Content:   "lorem ipsum dolor sit amet consetetur.",
		OwnerId:   0,
		CreatedAt: 1442284669000,
		UpdatedAt: 1442284669000,
	}, &model.Note{
		Id:        0,
		Title:     "Test Title 3",
		Content:   "lorem ipsum dolor sit amet consetetur.",
		OwnerId:   0,
		CreatedAt: 1442284669000,
		UpdatedAt: 1442284669000,
	})

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ListNotes(ctx, w, r)
		},
	))
	defer server.Close()

	resp := request(t, server, http.StatusOK, nil)
	assert.NotNil(resp)

	notes := resp.([]interface{})
	assert.Equal(len(notes), 3)
	assert.Equal("Test Title 1", notes[0].(map[string]interface{})["title"])
	assert.Equal("Test Title 2", notes[1].(map[string]interface{})["title"])
	assert.Equal("Test Title 3", notes[2].(map[string]interface{})["title"])
}

func TestAddNote(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{})

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			AddNote(ctx, w, r)
		},
	))
	defer server.Close()

	resp := request(t, server, http.StatusOK, map[string]interface{}{
		"title":    "Test Title",
		"content":  "lorem ipsum dolor sit amet consetetur.",
		"owner_id": 0,
	})
	assert.NotNil(resp)

	note := resp.(map[string]interface{})
	assert.Equal("Test Title", note["title"])
	assert.Equal("lorem ipsum dolor sit amet consetetur.", note["content"])
	assert.EqualValues(0, note["ownerId"])
	assert.NotZero(note["createdAt"])
	assert.NotZero(note["updatedAt"])

	count, err := dbMap.SelectInt("SELECT COUNT(id) FROM notes")
	assert.Nil(err)
	assert.EqualValues(1, count)
}
