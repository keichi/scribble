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
	auth := ctx.Value("auth").(*auth.Context)
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
			model.ShareStatePublic, model.ShareStatePrivate,
			auth.User.ID, limit, offset)
	} else {
		_, err = db.Select(&notes, `select * from notes where share_state = ?
				limit ? offset ?`,
			model.ShareStatePublic, limit, offset)
	}
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	return notes, nil
}

// ListNotes handles list notes requests
var ListNotes = WrapJSONHandler(nil, listNotes)

func addNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.Context)
	note := req.(*model.Note)

	if !note.Authorize(auth.User, model.ActionCreate) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	if auth.IsLoggedIn {
		note.OwnerID = auth.User.ID
	} else {
		note.OwnerID = 0
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

// AddNote handles add note requests
var AddNote = WrapJSONHandler(model.Note{}, addNote)

func getNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.Context)
	noteID, _ := strconv.ParseInt(kami.Param(ctx, "noteId"), 10, 64)

	note := new(model.Note)
	err := db.SelectOne(note, "select * from notes where id = ?", noteID)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !note.Authorize(auth.User, model.ActionRead) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	return note, nil
}

// GetNote handles get note requests
var GetNote = WrapJSONHandler(nil, getNote)

func updateNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.Context)
	newNote := req.(*model.Note)
	noteID, _ := strconv.ParseInt(kami.Param(ctx, "noteId"), 10, 64)

	note := new(model.Note)
	err := db.SelectOne(note, "select * from notes where id = ?", noteID)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !note.Authorize(auth.User, model.ActionUpdate) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	note.Title = newNote.Title
	note.Content = newNote.Content
	note.OwnerID = newNote.OwnerID
	note.UpdatedAt = time.Now().UnixNano()

	if _, err := db.Update(note); err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Update failed: %v", err),
		}
	}

	return note, nil
}

// UpdateNote handles update note requests
var UpdateNote = WrapJSONHandler(model.Note{}, updateNote)

func deleteNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.Context)
	noteID, _ := strconv.ParseInt(kami.Param(ctx, "noteId"), 10, 64)

	note := new(model.Note)
	err := db.SelectOne(note, "select * from notes where id = ?", noteID)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !note.Authorize(auth.User, model.ActionDelete) {
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

// DeleteNote handles delete note requests
var DeleteNote = WrapJSONHandler(nil, deleteNote)
