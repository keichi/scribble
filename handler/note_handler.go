package handler

import (
	"golang.org/x/net/context"
	"net/http"
	"strconv"

	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/model"
	"github.com/keichi/scribble/util"
)

func ListNotes(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := ctx.Value("db").(*gorp.DbMap)

	var notes []model.Note
	if _, err := db.Select(&notes, "select * from notes"); err != nil {
		panic(err)
	}

	util.RenderJson(http.StatusOK, notes, w)
}

func AddNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func GetNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := ctx.Value("db").(*gorp.DbMap)
	noteId, err := strconv.Atoi(ctx.Value("noteId").(string))

	if err != nil {
		util.RenderJson(http.StatusBadRequest, map[string]string {
			"message": "Failed to get note id",
		}, w)
		return
	}

	var note model.Note
	err = db.SelectOne(&note, "select * from notes where id = ?", noteId)
	if err != nil {
		util.RenderJson(http.StatusBadRequest, map[string]string {
			"message": "Note with specified id does not exist",
		}, w)
		return
	}

	util.RenderJson(http.StatusOK, note, w)
}

func UpdateNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func DeleteNote(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}
