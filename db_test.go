package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	db := newInMemoryDb()

	_, user := db.dbWriteNewUser("user1", "passwd1")
	assert.Equal(t, "user1", user.Username)
}

func TestDbWriteAndRead(t *testing.T) {
	db := newInMemoryDb()

	id, _ := db.dbWriteNewUser("user", "passwd")
	user, _ := db.dbGetUser(id)
	assert.Equal(t, "user", user.Username)
}

func TestDbGetPassoword(t *testing.T) {
	db := newInMemoryDb()

	db.dbWriteNewUser("user", "passwd")
	pass, _ := db.dbGetPasswordForUsername("user")
	assert.Equal(t, "passwd", pass)
}
