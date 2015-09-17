package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/guregu/kami"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/handler"
	"github.com/keichi/scribble/model"
)

func initDB() *gorp.DbMap {
	// TODO Use MySQL at production environment
	db, err := sql.Open("sqlite3", "scribble.db")
	if err != nil {
		panic(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(model.Session{}, "sessions").SetKeys(true, "ID")
	dbMap.AddTableWithName(model.Note{}, "notes").SetKeys(true, "ID")
	dbMap.AddTableWithName(model.Image{}, "images").SetKeys(true, "ID")
	dbMap.CreateTablesIfNotExists()

	return dbMap
}

func initS3() *s3.Bucket {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}

	s3 := s3.New(auth, aws.APNortheast)
	bucket := s3.Bucket("scribble-image-store")

	return bucket
}

func main() {
	dbMap := initDB()
	defer dbMap.Db.Close()

	// TODO Trace SQL only during development
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "scribble: ", log.Lmicroseconds))

	bucket := initS3()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "s3", bucket)
	kami.Context = ctx

	kami.PanicHandler = handler.Panic

	// Middlwares
	kami.Use("/api/", auth.Middleware)

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

	// Image APIs
	kami.Post("/api/images", handler.AddImage)
	kami.Get("/api/images/:imageId", handler.GetImage)
	kami.Delete("/api/images/:imageId", handler.DeleteImage)

	kami.Serve()
}
