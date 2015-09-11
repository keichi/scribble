package util

import (
	"net/http"
	"encoding/json"

	"golang.org/x/net/context"
)

func JsonResponseMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	w.Header().Set("Content-Type", "application/json")

	return ctx
}

func RenderJson(st int, v interface{}, w http.ResponseWriter) {
	w.WriteHeader(st)

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}
