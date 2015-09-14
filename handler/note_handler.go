package handler

import (
	"golang.org/x/net/context"
	"net/http"
	"strconv"

	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/model"
)

func listNotes(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)

	var notes []model.Note
	if _, err := db.Select(&notes, "select * from notes"); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return notes, nil
}

var ListNotes = WrapJsonHandler(nil, listNotes)

func addNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	note := req.(*model.Note)

	if err := db.Insert(note); err != nil {
		return nil, &ErrorResponse{http.StatusInternalServerError, err.Error()}
	}

	return note, nil
}

var AddNote = WrapJsonHandler(model.Note{}, addNote)

func getNote(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	noteId, err := strconv.Atoi(ctx.Value("noteId").(string))

	if err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, err.Error()}
	}

	var note model.Note
	err = db.SelectOne(&note, "select * from notes where id = ?", noteId)
	if err != nil {
		return nil, &ErrorResponse{http.StatusBadRequest, err.Error()}
	}

	return note, nil
}

var GetNote = WrapJsonHandler(nil, getNote)

func UpdateNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func DeleteNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}
