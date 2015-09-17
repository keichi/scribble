package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guregu/kami"
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
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{User: &model.User{}})

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

	resp := request(t, server.URL, http.StatusOK, nil)
	assert.NotNil(resp)

	notes := resp.([]interface{})
	assert.EqualValues(3, len(notes))
	assert.Equal("Test Title 1", notes[0].(map[string]interface{})["title"])
	assert.Equal("Test Title 2", notes[1].(map[string]interface{})["title"])
	assert.Equal("Test Title 3", notes[2].(map[string]interface{})["title"])
}

func TestListNotesPagination(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{User: &model.User{}})

	for i := 1; i <= 100; i++ {
		dbMap.Insert(&model.Note{
			Id:        0,
			Title:     fmt.Sprintf("Test Title %d", i),
			Content:   "lorem ipsum dolor sit amet consetetur.",
			OwnerId:   0,
			CreatedAt: 1442284669000,
			UpdatedAt: 1442284669000,
		})
	}

	kami.Reset()
	kami.Context = ctx
	kami.Post("/api/notes", ListNotes)
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	resp := request(t, server.URL+"/api/notes", http.StatusOK, nil)
	assert.NotNil(resp)
	assert.EqualValues(10, len(resp.([]interface{})))

	resp = request(t, server.URL+"/api/notes?limit=25", http.StatusOK, nil)
	assert.NotNil(resp)
	assert.EqualValues(25, len(resp.([]interface{})))

	resp = request(t, server.URL+"/api/notes?limit=50", http.StatusOK, nil)
	assert.NotNil(resp)
	assert.EqualValues(50, len(resp.([]interface{})))
	note := resp.([]interface{})[0].(map[string]interface{})
	assert.Equal("Test Title 1", note["title"])

	resp = request(t, server.URL+"/api/notes?offset=25&limit=50", http.StatusOK, nil)
	assert.NotNil(resp)
	assert.EqualValues(50, len(resp.([]interface{})))
	note = resp.([]interface{})[0].(map[string]interface{})
	assert.Equal("Test Title 26", note["title"])
}

func TestAddNote(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{User: &model.User{}})

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			AddNote(ctx, w, r)
		},
	))
	defer server.Close()

	resp := request(t, server.URL, http.StatusOK, map[string]interface{}{
		"title":   "Test Title",
		"content": "lorem ipsum dolor sit amet consetetur.",
		"ownerId": 0,
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

func TestGetNote(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{User: &model.User{}})

	dbMap.Insert(&model.Note{
		Id:        0,
		Title:     "Test Title 1",
		Content:   "lorem ipsum dolor sit amet consetetur.",
		OwnerId:   0,
		CreatedAt: 1442284669000,
		UpdatedAt: 1442292926000,
	})

	kami.Reset()
	kami.Context = ctx
	kami.Post("/api/notes/:noteId", GetNote)
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	resp := request(t, server.URL+"/api/notes/1", http.StatusOK, nil)
	assert.NotNil(resp)

	note := resp.(map[string]interface{})
	assert.Equal("Test Title 1", note["title"])
	assert.Equal("lorem ipsum dolor sit amet consetetur.", note["content"])
	assert.EqualValues(0, note["ownerId"])
	assert.EqualValues(1442284669000, note["createdAt"])
	assert.EqualValues(1442292926000, note["updatedAt"])
}

func TestUpdateNote(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{User: &model.User{}})

	dbMap.Insert(&model.Note{
		Id:        0,
		Title:     "Test Title 1",
		Content:   "lorem ipsum dolor sit amet consetetur.",
		OwnerId:   0,
		CreatedAt: 1442284669000,
		UpdatedAt: 1442292926000,
	})

	kami.Reset()
	kami.Context = ctx
	kami.Post("/api/notes/:noteId", UpdateNote)
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	resp := request(t, server.URL+"/api/notes/1", http.StatusOK, map[string]interface{}{
		"title":   "Test Title 2",
		"content": "hoge piyo hoge piyo.",
		"ownerId": 1,
	})
	assert.NotNil(resp)

	note := resp.(map[string]interface{})
	assert.Equal("Test Title 2", note["title"])
	assert.Equal("hoge piyo hoge piyo.", note["content"])
	assert.EqualValues(1, note["ownerId"])
	assert.EqualValues(1442284669000, note["createdAt"])

	count, err := dbMap.SelectInt("SELECT COUNT(id) FROM notes")
	assert.Nil(err)
	assert.EqualValues(1, count)

	n := new(model.Note)
	err = dbMap.SelectOne(n, "SELECT * FROM notes")
	assert.Nil(err)
	assert.Equal("Test Title 2", n.Title)
	assert.Equal("hoge piyo hoge piyo.", n.Content)
	assert.EqualValues(1, n.OwnerId)
	assert.EqualValues(1442284669000, n.CreatedAt)
}

func TestDeleteNote(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDb()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.AuthContext{User: &model.User{}})

	dbMap.Insert(&model.Note{
		Id:        0,
		Title:     "Test Title 1",
		Content:   "lorem ipsum dolor sit amet consetetur.",
		OwnerId:   0,
		CreatedAt: 1442284669000,
		UpdatedAt: 1442292926000,
	})

	count, err := dbMap.SelectInt("SELECT COUNT(id) FROM notes")
	assert.Nil(err)
	assert.EqualValues(1, count)

	kami.Reset()
	kami.Context = ctx
	kami.Post("/api/notes/:noteId", DeleteNote)
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	request(t, server.URL+"/api/notes/1", http.StatusOK, nil)

	count, err = dbMap.SelectInt("SELECT COUNT(id) FROM notes")
	assert.Nil(err)
	assert.EqualValues(0, count)
}
