package handler

import (
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

func TestListNotesAuth(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.Context{
		true,
		&model.User{
			1,
			"testuser",
			"000000",
			"test@emaple.com",
			1442284669000,
			1442284669000,
		},
		&model.Session{},
	})

	dbMap.Insert(&model.Note{
		ID:         0,
		Title:      "Test Title 1",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    1,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 2",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 3",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePublic,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
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
	assert.EqualValues(2, len(notes))
	assert.Equal("Test Title 1", notes[0].(map[string]interface{})["title"])
	assert.Equal("Test Title 3", notes[1].(map[string]interface{})["title"])
}

func TestGetNoteAuth(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.Context{
		true,
		&model.User{
			1,
			"testuser",
			"000000",
			"test@emaple.com",
			1442284669000,
			1442284669000,
		},
		&model.Session{},
	})

	dbMap.Insert(&model.Note{
		ID:         0,
		Title:      "Test Title 1",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    1,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 2",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 3",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePublic,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	})

	kami.Reset()
	kami.Context = ctx
	kami.Post("/api/notes/:noteId", GetNote)
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	resp := request(t, server.URL+"/api/notes/1", http.StatusOK, nil)
	assert.NotNil(resp)

	request(t, server.URL+"/api/notes/2", http.StatusUnauthorized, nil)

	request(t, server.URL+"/api/notes/3", http.StatusOK, nil)
}

func TestUpdateNoteAuth(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.Context{
		true,
		&model.User{
			1,
			"testuser",
			"000000",
			"test@emaple.com",
			1442284669000,
			1442284669000,
		},
		&model.Session{},
	})

	dbMap.Insert(&model.Note{
		ID:         0,
		Title:      "Test Title 1",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    1,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 2",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 3",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePublic,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
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

	request(t, server.URL+"/api/notes/2", http.StatusUnauthorized, map[string]interface{}{
		"title":   "Test Title 2",
		"content": "hoge piyo hoge piyo.",
		"ownerId": 1,
	})
	request(t, server.URL+"/api/notes/3", http.StatusUnauthorized, map[string]interface{}{
		"title":   "Test Title 2",
		"content": "hoge piyo hoge piyo.",
		"ownerId": 1,
	})
}

func TestDeleteNoteAuth(t *testing.T) {
	assert := assert.New(t)

	dbMap := initDB()
	defer dbMap.Db.Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "auth", &auth.Context{
		true,
		&model.User{
			1,
			"testuser",
			"000000",
			"test@emaple.com",
			1442284669000,
			1442284669000,
		},
		&model.Session{},
	})

	dbMap.Insert(&model.Note{
		ID:         0,
		Title:      "Test Title 1",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    1,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 2",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePrivate,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	}, &model.Note{
		ID:         0,
		Title:      "Test Title 3",
		Content:    "lorem ipsum dolor sit amet consetetur.",
		OwnerID:    2,
		ShareState: model.ShareStatePublic,
		CreatedAt:  1442284669000,
		UpdatedAt:  1442284669000,
	})

	kami.Reset()
	kami.Context = ctx
	kami.Post("/api/notes/:noteId", DeleteNote)
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	request(t, server.URL+"/api/notes/1", http.StatusOK, nil)
	request(t, server.URL+"/api/notes/2", http.StatusUnauthorized, nil)
	request(t, server.URL+"/api/notes/3", http.StatusUnauthorized, nil)

	count, err := dbMap.SelectInt("SELECT COUNT(id) FROM notes")
	assert.Nil(err)
	assert.EqualValues(2, count)
}
