package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/gorp.v1"

	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)


func listMyNotes(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
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
	_, err := db.Select(&notes, `select * from notes where
									owner_id = ? limit ? offset ?`,
						auth.User.ID, limit, offset)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	return notes, nil
}

// ListNotes handles list notes requests
var ListMyNotes = wrapJSONHandler(emptyRequest{}, listMyNotes)
