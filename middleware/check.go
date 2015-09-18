package middleware

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strconv"

	"github.com/guregu/kami"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/handler"
)

// CheckIfNoteExists middleware checks if a note with the  specified id exists
// It also ensures the note id is an integer
func CheckIfNoteExists(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	db := ctx.Value("db").(*gorp.DbMap)
	noteID, err := strconv.Atoi(kami.Param(ctx, "noteId"))

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
