package handler


import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guregu/kami"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestPanic(t *testing.T) {
	assert := assert.New(t)

	kami.Reset()
	kami.Context = context.Background()
	kami.PanicHandler = Panic
	kami.Post("/panic", func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("This is a test panic"))
	})
	server := httptest.NewServer(kami.Handler())
	defer server.Close()

	resp := request(t, server.URL + "/panic", http.StatusInternalServerError, nil)
	assert.NotNil(resp)
	assert.Equal("Handler panic: This is a test panic", resp.(map[string]interface{})["message"])
	assert.NotEmpty(resp.(map[string]interface{})["context"])
}
