package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	assert := assert.New(t)

	token1 := NewToken()
	assert.NotEmpty(token1)

	raw, _ := base64.StdEncoding.DecodeString(token1)
	assert.EqualValues(16, len(raw))

	token2 := NewToken()

	assert.NotEqual(token1, token2)
}
