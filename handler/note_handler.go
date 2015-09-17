package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/guregu/kami"
	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

func listNotes(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	queryParams := ctx.Value("query").(url.Values)

	limit := 10
	offset := 0

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		i, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, &ErrorResponse{
				http.StatusInternalServerError,
				fmt.Sprintf("Invalid limit parameter format: %v", err),
			}
		}
		limit = i
	}

	offsetStr := queryParams.Get("offset")
	if offsetStr != "" {
		i, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, &ErrorResponse{
				http.StatusInternalServerError,
				fmt.Sprintf("Invalid offset parameter format: %v", err),
			}
		}
		offset = i
	}

	var notes []model.Note
	var err error
	if auth.IsLoggedIn {
		_, err = db.Select(&notes, `select * from notes where share_state = ?
				or share_state = ? and owner_id = ? limit ? offset ?`,
			model.SHARE_STATE_PUBLIC, model.SHARE_STATE_PRIVATE,
			auth.User.Id, limit, offset)
	} else {
		_, err = db.Select(&notes, `select * from notes where share_state = ?
				limit ? offset ?`,
			model.SHARE_STATE_PUBLIC, limit, offset)
	}
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	return notes, nil
}

var ListNotes = WrapJsonHandler(nil, listNotes)

func addNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	note := req.(*model.Note)

	if !note.Authorize(auth.User, model.ACTION_CREATE) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	if auth.IsLoggedIn {
		note.OwnerId = auth.User.Id
	} else {
		note.OwnerId = 0
	}
	note.CreatedAt = time.Now().UnixNano()
	note.UpdatedAt = time.Now().UnixNano()

	if err := db.Insert(note); err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Insert failed: %v", err),
		}
	}

	return note, nil
}

var AddNote = WrapJsonHandler(model.Note{}, addNote)

func getNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	noteId, err := strconv.Atoi(kami.Param(ctx, "noteId"))

	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid note id format: %v", err),
		}
	}

	note := new(model.Note)
	err = db.SelectOne(note, "select * from notes where id = ?", noteId)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !note.Authorize(auth.User, model.ACTION_READ) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	return note, nil
}

var GetNote = WrapJsonHandler(nil, getNote)

func updateNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	newNote := req.(*model.Note)
	noteId, err := strconv.Atoi(kami.Param(ctx, "noteId"))

	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid note id format: %v", err),
		}
	}

	note := new(model.Note)
	err = db.SelectOne(note, "select * from notes where id = ?", noteId)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !note.Authorize(auth.User, model.ACTION_UPDATE) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	note.Title = newNote.Title
	note.Content = newNote.Content
	note.OwnerId = newNote.OwnerId
	note.UpdatedAt = time.Now().UnixNano()

	if _, err := db.Update(note); err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Update failed: %v", err),
		}
	}

	return note, nil
}

var UpdateNote = WrapJsonHandler(model.Note{}, updateNote)

func deleteNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	noteId, err := strconv.Atoi(kami.Param(ctx, "noteId"))

	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid note id format: %v", err),
		}
	}

	note := new(model.Note)
	err = db.SelectOne(note, "select * from notes where id = ?", noteId)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !note.Authorize(auth.User, model.ACTION_DELETE) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	if _, err := db.Delete(note); err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Delete failed: %v", err),
		}
	}

	return nil, nil
}

var DeleteNote = WrapJsonHandler(nil, deleteNote)
