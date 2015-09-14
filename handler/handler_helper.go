package handler

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	"reflect"

	"fmt"
	"github.com/guregu/kami"
)

type JsonHandler func(ctx context.Context, req interface{}) (interface{}, *ErrorResponse)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

type emptyRequest struct{}

func WrapJsonHandler(v interface{}, h JsonHandler) kami.HandlerFunc {
	t := reflect.TypeOf(v)

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		decoder := json.NewDecoder(r.Body)

		var input interface{}
		if t != reflect.TypeOf(nil) {
			input = reflect.New(t).Interface()
			if err := decoder.Decode(input); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				encoder.Encode(map[string]string{
					"message": fmt.Sprintf("json decode fail: %v", err),
				})
				return
			}
		}

		output, resp := h(ctx, input)
		if resp != nil {
			w.WriteHeader(resp.StatusCode)
			encoder.Encode(map[string]string{
				"message": resp.Error(),
			})
			return
		}
		if output == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := encoder.Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(map[string]string{
				"message": fmt.Sprintf("json encode fail: %v", err),
			})
			return
		}
	}
}
