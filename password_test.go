package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	pass, err := hashPassword("MyPass")
	assert.Nil(t, err)
	assert.Greater(t, len(pass), 16+32)
}

func TestGenerateRandomBytes(t *testing.T) {
	bytes, err := generateRandomBytes(16)
	assert.Nil(t, err)
	assert.Equal(t, len(bytes), 16)
}

func TestComparePasswordAndHash(t *testing.T) {
	hash, err := hashPassword("MyPass")
	assert.Nil(t, err)

	match, err := comparePasswordAndHash("MyPass", hash)
	assert.Nil(t, err)
	assert.True(t, match)
}

func TestComparePasswordAndHashNoMatch(t *testing.T) {
	hash, err := hashPassword("MyPass")
	assert.Nil(t, err)

	match, err := comparePasswordAndHash("MyPassWrong", hash)
	assert.Nil(t, err)
	assert.False(t, match)
}
