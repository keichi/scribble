package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/goamz/goamz/s3/s3test"
	"github.com/guregu/kami"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/handler"
	"github.com/keichi/scribble/middleware"
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
	//	auth, err := aws.EnvAuth()
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	s3 := s3.New(auth, aws.APNortheast)
	//	bucket := s3.Bucket("scribble-image-store")
	//
	//	return bucket

	// TODO Fake S3 server -- only use it during development
	srv, err := s3test.NewServer(&s3test.Config{})
	if err != nil {
		panic(err)
	}
	region := aws.Region{
		Name:                 "dummy-region-1",
		S3Endpoint:           srv.URL(),
		S3LocationConstraint: true,
	}

	bucket := s3.New(aws.Auth{}, region).Bucket("scribble-image-store")
	bucket.PutBucket(s3.PublicRead)

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

	// TODO Allow any CORS request -- only during development!
	kami.Use("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		return ctx
	})
	kami.Options("/*path", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	kami.Serve()
}
