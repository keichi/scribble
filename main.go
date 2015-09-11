package main

import (
	"database/sql"

	"github.com/guregu/kami"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/handler"
	"github.com/keichi/scribble/model"
	"github.com/keichi/scribble/util"
)

func InitDB() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "scribble.db")
	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Session{}, "sessions").SetKeys(true, "Id")
	dbMap.CreateTablesIfNotExists()

	return dbMap
}

func main() {
	dbMap := InitDB()
	defer dbMap.Db.Close()

	ctx := context.Background()
	kami.Context = context.WithValue(ctx, "db", dbMap)

	// Middlwares
	kami.Use("/api/", util.JsonResponseMiddleware)
	kami.Use("/api/", auth.AuthMiddleware)

	// Authentication APIs
	kami.Post("/api/register", handler.Register)
	kami.Post("/api/login", handler.Login)
	kami.Post("/api/logout", handler.Logout)

	// Note APIs
	kami.Get("/api/notes", handler.ListNotes)
	kami.Post("/api/notes", handler.AddNote)
	kami.Get("/api/notes/:noteId", handler.GetNote)
	kami.Post("/api/notes/:noteId", handler.UpdateNote)
	kami.Delete("/api/notes/:noteId", handler.DeleteNote)

	kami.Serve()
}
