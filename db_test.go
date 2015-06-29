package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	db := newInMemoryDb()

	user, _ := db.dbWriteNewUser("user1", "passwd1")
	assert.Equal(t, "user1", user.Username)
}

func TestDbWriteAndRead(t *testing.T) {
	db := newInMemoryDb()

	createdUser, _ := db.dbWriteNewUser("user", "passwd")
	gotUser, _ := db.dbGetUser(createdUser.Id)
	assert.Equal(t, "user", gotUser.Username)
}

func TestDbGetUserAndPassword(t *testing.T) {
	db := newInMemoryDb()

	db.dbWriteNewUser("user", "passwd")
	user, _ := db.dbGetUserAndPasswordForUsername("user")
	assert.Equal(t, "user", user.user.Username)
	assert.Equal(t, "passwd", *(user.password))
}

func TestDbGetNonExistingUserAndPassword(t *testing.T) {
	db := newInMemoryDb()

	userWithPassword, err := db.dbGetUserAndPasswordForUsername("user")
	assert.Error(t, err)
	assert.Nil(t, userWithPassword.password)
}

//TODO: needs test for username wrong (so it doesnt fall back to "")
