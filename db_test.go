package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	dbWriteNewUser("user1")
}

func TestDbWriteAndRead(t *testing.T) {
	id := dbWriteNewUser("user1")
	user, _ := dbGetUser(id)
	assert.Equal(t, "user1", user.username)
}
