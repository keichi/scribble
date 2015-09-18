package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	assert := assert.New(t)

	token1 := NewToken()
	assert.NotEmpty(token1)
	assert.EqualValues(32, len(token1))

	token2 := NewToken()

	assert.NotEqual(token1, token2)
}
