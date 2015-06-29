package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbWrite(t *testing.T) {
	db := newInMemoryDb()

	user, _ := db.writeNewUser("user1", "passwd1")
	assert.Equal(t, "user1", user.Username)
}

func TestDbWriteAndRead(t *testing.T) {
	db := newInMemoryDb()

	createdUser, _ := db.writeNewUser("user", "passwd")
	gotUser, _ := db.getUser(createdUser.Id)
	assert.Equal(t, "user", gotUser.Username)
}

func TestgetUserAndPassword(t *testing.T) {
	db := newInMemoryDb()

	db.writeNewUser("user", "passwd")
	user, _ := db.getUserAndPasswordForUsername("user")
	assert.Equal(t, "user", user.user.Username)
	assert.Equal(t, "passwd", *(user.password))
}

func TestDbGetNonExistingUserAndPassword(t *testing.T) {
	db := newInMemoryDb()

	userWithPassword, err := db.getUserAndPasswordForUsername("user")
	assert.Error(t, err)
	assert.Nil(t, userWithPassword.password)
}
