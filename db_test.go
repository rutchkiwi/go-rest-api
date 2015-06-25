package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	_, user := dbWriteNewUser("user1")
	assert.Equal(t, "user1", user.Username)
}

func TestDbWriteAndRead(t *testing.T) {
	id, _ := dbWriteNewUser("user1")
	user, _ := dbGetUser(id)
	assert.Equal(t, "user1", user.Username)
}
