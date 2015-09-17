package handler

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"

	"github.com/guregu/kami"
)

// Panic handles panics happened in kami handler functions
func Panic(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)

	encoder := json.NewEncoder(w)
	e := kami.Exception(ctx)
	resp := map[string]string{
		"message": fmt.Sprintf("Handler panic: %v", e),
		"context": fmt.Sprint(ctx),
	}

	encoder.Encode(resp)

	log.Println("Panicked in kami handler")
	log.Printf("Panic detail: %v", e)
	log.Printf("Current context: %v", ctx)
}
