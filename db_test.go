package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	writeNewUser("user1")
}

func TestDbWriteAndRead(t *testing.T) {
	id := writeNewUser("user1")
	user := getUser(id)
	assert.Equal(t, "user1", user.username)
}
