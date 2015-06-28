package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginInAsAdmin(t *testing.T) {
	buildWebservice()

	getHttpReq, _ := http.NewRequest("GET", "/me", nil)

	getHttpReq.Header.Set("Authorization", basicAuthEncode("admin", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	assert.Equal(t, 200, httpWriter.Code)
}

func TestAdminUsersPage(t *testing.T) {
	buildWebservice()

	getHttpReq, _ := http.NewRequest("GET", "/admin/users", nil)

	getHttpReq.Header.Set("Authorization", basicAuthEncode("admin", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 200, httpWriter.Code)

	var users []UserWithConnections
	json.Unmarshal(httpWriter.Body.Bytes(), &users)
	require.Len(t, users, 1)
	assert.Equal(t, "admin", users[0].User.Username)
	assert.Len(t, users[0].Connections, 0)
}

func TestAdminUsersPageListsConnections(t *testing.T) {
	buildWebservice()

	//add some connections
	registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")
	addConnection(t, "viktor", "pass", user2Id)

	getHttpReq, _ := http.NewRequest("GET", "/admin/users", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("admin", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 200, httpWriter.Code)

	var users []UserWithConnections
	json.Unmarshal(httpWriter.Body.Bytes(), &users)
	require.Len(t, users, 3)

	assert.Equal(t, "admin", users[0].User.Username)
	assert.Len(t, users[0].Connections, 0)

	assert.Equal(t, "viktor", users[1].User.Username)
	require.Len(t, users[1].Connections, 1)
	assert.Equal(t, "user2", users[1].Connections[0].Username)
	assert.Equal(t, user2Id, users[1].Connections[0].Id)
}
