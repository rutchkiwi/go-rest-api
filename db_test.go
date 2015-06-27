package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	db := newInMemoryDb()

	user := db.dbWriteNewUser("user1", "passwd1")
	assert.Equal(t, "user1", user.Username)
}

func TestDbWriteAndRead(t *testing.T) {
	db := newInMemoryDb()

	createdUser := db.dbWriteNewUser("user", "passwd")
	gotUser, _ := db.dbGetUser(createdUser.id)
	assert.Equal(t, "user", gotUser.Username)
}

func TestDbGetUserAndPassword(t *testing.T) {
	db := newInMemoryDb()

	db.dbWriteNewUser("user", "passwd")
	user, pass, _ := db.dbGetUserAndPasswordForUsername("user")
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, "passwd", *pass)
}

func TestDbGetNonExistingUserAndPassword(t *testing.T) {
	db := newInMemoryDb()

	_, pass, err := db.dbGetUserAndPasswordForUsername("user")
	assert.Error(t, err)
	assert.Nil(t, pass)
}

//TODO: needs test for username wrong (so it doesnt fall back to "")
