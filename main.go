package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/goamz/goamz/s3/s3test"
	"github.com/guregu/kami"
	"github.com/mattes/migrate/migrate"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/handler"
	"github.com/keichi/scribble/middleware"
	"github.com/keichi/scribble/model"
)

func initDB() *gorp.DbMap {
	errors, ok := migrate.UpSync("sqlite3://scribble.db", "./migrations")
	if !ok {
		panic(errors)
	}

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

	return dbMap
}

func initS3() *s3.Bucket {
	region := aws.Region{}
	auth, err := aws.EnvAuth()

	if err != nil {
		srv, err := s3test.NewServer(&s3test.Config{})
		if err != nil {
			panic(err)
		}
		region = aws.Region{
			Name:                 "dummy-region-1",
			S3Endpoint:           srv.URL(),
			S3LocationConstraint: true,
		}
	} else {
		region = aws.APNortheast
	}

	bucketName := os.Getenv("S3_BUCKET")
	if bucketName == "" {
		bucketName = "scribble-image-store"
	}

	s3 := s3.New(auth, region)
	bucket := s3.Bucket(bucketName)

	return bucket
}

func main() {
	dbMap := initDB()
	defer dbMap.Db.Close()

	bucket := initS3()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dbMap)
	ctx = context.WithValue(ctx, "s3", bucket)
	kami.Context = ctx

	kami.PanicHandler = handler.Panic

	// Middlwares
	kami.Use("/api/", middleware.Auth)
	kami.Use("/api/notes/:noteId", middleware.CheckIfNoteExists)
	kami.Use("/api/notes/:noteId/images/:imageId", middleware.CheckIfImageExists)

	// Authentication APIs
	kami.Get("/api/auth", handler.GetMyUser)
	kami.Post("/api/auth/register", handler.Register)
	kami.Post("/api/auth/login", handler.Login)
	kami.Post("/api/auth/logout", handler.Logout)

	// Note APIs
	kami.Get("/api/notes", handler.ListNotes)
	kami.Post("/api/notes", handler.AddNote)
	kami.Get("/api/notes/:noteId", handler.GetNote)
	kami.Put("/api/notes/:noteId", handler.UpdateNote)
	kami.Delete("/api/notes/:noteId", handler.DeleteNote)

	// Image APIs
	kami.Post("/api/notes/:noteId/images", handler.AddImage)
	kami.Get("/api/notes/:noteId/images/:imageId", handler.GetImage)
	kami.Delete("/api/notes/:noteId/images/:imageId", handler.DeleteImage)

	// Personal APIs
	kami.Use("/api/my/", middleware.CheckIfLoggedIn)
	kami.Get("/api/my/notes", handler.ListMyNotes)
	kami.Post("/api/my/notes", handler.AddNote)
	kami.Get("/api/my/notes/:noteId", handler.GetNote)
	kami.Put("/api/my/notes/:noteId", handler.UpdateNote)
	kami.Delete("/api/my/notes/:noteId", handler.DeleteNote)

	// Ping API
	kami.Get("/api/ping", handler.Ping)

	fileServer := http.FileServer(http.Dir("static"))
	kami.Get("/", fileServer)
	kami.Get("/404.html", fileServer)
	kami.Get("/robots.txt", fileServer)
	kami.Get("/favicon.ico", fileServer)
	kami.Get("/fonts/*path", fileServer)
	kami.Get("/scripts/*path", fileServer)
	kami.Get("/styles/*path", fileServer)

	kami.Serve()
}
