package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/guregu/kami"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"
	"github.com/rlmcpherson/s3gof3r"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/handler"
	"github.com/keichi/scribble/model"
)

func InitDB() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "scribble.db")
	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Session{}, "sessions").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Note{}, "notes").SetKeys(true, "Id")
	dbMap.AddTableWithName(model.Image{}, "images").SetKeys(true, "Id")
	dbMap.CreateTablesIfNotExists()

	return dbMap
}

func InitS3() *s3gof3r.Bucket {
	keys, err := s3gof3r.EnvKeys()
	if err != nil {
		panic(err)
	}

	s3 := s3gof3r.New("", keys)
	bucket := s3.Bucket("scribble-image-store")

	return bucket
}

func main() {
	dbMap := InitDB()
	defer dbMap.Db.Close()
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "scribble: ", log.Lmicroseconds))

	bucket := InitS3()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "s3", bucket)
	kami.Context = ctx

	kami.PanicHandler = handler.Panic

	// Middlwares
	kami.Use("/api/", auth.AuthMiddleware)

	// Authentication APIs
	kami.Post("/api/register", handler.Register)
	kami.Post("/api/login", handler.Login)
	kami.Post("/api/logout", handler.Logout)

	// Note APIs
	kami.Get("/api/notes", handler.ListNotes)
	kami.Post("/api/notes", handler.AddNote)
	kami.Get("/api/notes/:noteId", handler.GetNote)
	kami.Put("/api/notes/:noteId", handler.UpdateNote)
	kami.Delete("/api/notes/:noteId", handler.DeleteNote)

	// Note APIs
	kami.Post("/api/images", handler.AddImage)
	kami.Get("/api/images/:imageId", handler.GetImage)
	kami.Delete("/api/images/:imageId", handler.DeleteImage)

	kami.Serve()
}
