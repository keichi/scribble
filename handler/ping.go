package handler
import (
	"net/http"
	"golang.org/x/net/context"
)

func Ping(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
