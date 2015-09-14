package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/guregu/kami"
	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

func listNotes(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)

	var notes []model.Note
	if _, err := db.Select(&notes, "select * from notes"); err != nil {
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

	return note, nil
}

var GetNote = WrapJsonHandler(nil, getNote)

func updateNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
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

	if _, err := db.Delete(note); err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Delete failed: %v", err),
		}
	}

	return nil, nil
}

var DeleteNote = WrapJsonHandler(nil, deleteNote)
