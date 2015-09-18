package handler

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"reflect"

	"github.com/guregu/kami"
)

// JSONHandler is a generic type for JSON API handlers
type JSONHandler func(ctx context.Context, req interface{}) (interface{}, *ErrorResponse)

// ErrorResponse represents errors returned from APIs
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

// Render writes the content of ErrorResponse as a JSON to a ResponseWriter
func (err *ErrorResponse) Render(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	encoder.Encode(map[string]string{
		"message": err.Error(),
	})
}

type emptyRequest struct{}

// wrapJSONHandler wraps JsonHandler as a kami.HandlerFunc
func wrapJSONHandler(v interface{}, h JSONHandler) kami.HandlerFunc {
	t := reflect.TypeOf(v)

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var input interface{}
		if t != reflect.TypeOf(emptyRequest{}) {
			input = reflect.New(t).Interface()
			if err := decoder.Decode(input); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				encoder.Encode(map[string]string{
					"message": fmt.Sprintf("JSON decode fail: %v", err),
				})
				return
			}
		}

		queryParams := r.URL.Query()
		ctx = context.WithValue(ctx, "query", queryParams)

		output, resp := h(ctx, input)
		if resp != nil {
			resp.Render(w)
			return
		}
		if output == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := encoder.Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(map[string]string{
				"message": fmt.Sprintf("JSON encode fail: %v", err),
			})
			return
		}
	}
}
