package middleware

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strconv"

	"github.com/guregu/kami"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/handler"
)

// CheckIfNoteExists middleware checks if a note with the  specified id exists
// It also ensures the note id is a valid integer
func CheckIfNoteExists(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	db := ctx.Value("db").(*gorp.DbMap)
	noteID, err := strconv.ParseInt(kami.Param(ctx, "noteId"), 10, 64)

	if err != nil {
		resp := &handler.ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid note id format: %v", err),
		}
		resp.Render(w)
		return nil
	}

	count, err := db.SelectInt("select count(id) from notes where id = ?", noteID)
	switch {
	case count <= 0:
		resp := &handler.ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Requested note does not exist"),
		}
		resp.Render(w)
		return nil
	case err != nil:
		resp := &handler.ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
		resp.Render(w)
		return nil
	}

	return ctx
}

// CheckIfImageExists middleware checks if an image with the  specified id exists
// It also ensures the image id is a valid integer
func CheckIfImageExists(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if CheckIfNoteExists(ctx, w, r) == nil {
		return nil
	}

	db := ctx.Value("db").(*gorp.DbMap)
	noteID, _ := strconv.ParseInt(kami.Param(ctx, "noteId"), 10, 64)
	imageID, err := strconv.ParseInt(kami.Param(ctx, "imageId"), 10, 64)

	if err != nil {
		resp := &handler.ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid image id format: %v", err),
		}
		resp.Render(w)
		return nil
	}

	count, err := db.SelectInt(`select count(id) from images where id = ?
								and note_id = ?`, imageID, noteID)
	switch {
	case count <= 0:
		resp := &handler.ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Requested image does not exist"),
		}
		resp.Render(w)
		return nil
	case err != nil:
		resp := &handler.ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
		resp.Render(w)
		return nil
	}

	return ctx
}

// CheckIfNoteExists middleware checks if user is logged in
func CheckIfLoggedIn(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	auth := ctx.Value("auth").(*auth.Context)

	// TODO development only
	if r.Method == "OPTIONS" {
		return ctx
	}

	if !auth.IsLoggedIn {
		resp := &handler.ErrorResponse{
			http.StatusUnauthorized,
			fmt.Sprintf("You have to logged in to execute this request"),
		}
		resp.Render(w)
		return nil
	}

	return ctx
}
