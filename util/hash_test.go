package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLHash(t *testing.T) {
	assert := assert.New(t)
	h, _ := BaseHashEncode([]byte("foobarbiz"), 2)
	assert.Equal(len(h), 2)

	h1, _ := BaseHashEncode([]byte("foobarbiz"), 18)
	h2, _ := BaseHashEncode([]byte("xxx"), 18)
	h3, _ := BaseHashEncode([]byte("foobarbiz"), 18)
	assert.Equal(len(h1), 18)
	assert.Equal(len(h2), 18)
	assert.Equal(len(h3), 18)
	assert.NotEqual(h1, h2)
	assert.NotEqual(h1, h3)
}

func TestGUID(t *testing.T) {
	assert := assert.New(t)
	g1 := GetUID()
	g2 := GetUID()
	g3 := GetUID()
	assert.GreaterOrEqual(len(g1), 12)
	assert.GreaterOrEqual(len(g2), 12)
	assert.GreaterOrEqual(len(g3), 12)
	assert.NotEqual(g1, g2)
	assert.NotEqual(g1, g3)
	assert.NotEqual(g2, g3)
}
