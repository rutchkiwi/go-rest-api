package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestListConnectionsEmpty(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")

	connections := listConnections(t)

	require.Len(t, connections, 0)
}

func listConnections(t *testing.T) []User {
	httpReq, _ := http.NewRequest("GET", "/connection", nil)
	httpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, httpReq)

	require.Equal(t, 200, httpWriter.Code)

	var connections []User
	json.Unmarshal(httpWriter.Body.Bytes(), &connections)
	return connections
}

func TestAddConnection(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")

	addConnection(t, "viktor", "pass", user2Id)
}

func TestAddConnectionIdempotent(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")

	addConnection(t, "viktor", "pass", user2Id)
	addConnection(t, "viktor", "pass", user2Id)
}

func TestAddConnectionBadInput(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")

	bodyString := fmt.Sprintf(`{badjson}`)
	bodyReader := strings.NewReader(bodyString)
	httpReq, _ := http.NewRequest("PUT", "/connection", bodyReader)
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	httpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, httpReq)

	require.Equal(t, 400, httpWriter.Code)
}

func TestAddConnectionBadId(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")

	bodyString := fmt.Sprintf(`{"id":%d}`, 999)
	bodyReader := strings.NewReader(bodyString)
	httpReq, _ := http.NewRequest("PUT", "/connection", bodyReader)
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	httpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, httpReq)

	require.Equal(t, 400, httpWriter.Code)
}

func addConnection(t *testing.T, fromUsername, fromPassword string, toUserId int64) {
	bodyString := fmt.Sprintf(`{"id":%d}`, toUserId)
	bodyReader := strings.NewReader(bodyString)
	httpReq, _ := http.NewRequest("PUT", "/connection", bodyReader)
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	httpReq.Header.Set("Authorization", basicAuthEncode(fromUsername, fromPassword))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, httpReq)

	require.Equal(t, 200, httpWriter.Code)
}

func TestListConnections(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")

	addConnection(t, "viktor", "pass", user2Id)

	connections := listConnections(t)
	require.Len(t, connections, 1)
	assert.Equal(t, "user2", connections[0].Username)
}

func TestListOnlyMyConnections(t *testing.T) {
	buildWebservice(true)
	viktorId := registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")
	user3Id := registerUser(t, "user3", "pass")

	addConnection(t, "user2", "pass", viktorId)
	addConnection(t, "user2", "pass", user2Id)
	addConnection(t, "user2", "pass", user3Id)
	addConnection(t, "user3", "pass", user2Id)

	addConnection(t, "viktor", "pass", user2Id)
	addConnection(t, "viktor", "pass", user3Id)

	connections := listConnections(t)
	require.Len(t, connections, 2)
	assert.Contains(t, connections, User{user2Id, "user2"})
	assert.Contains(t, connections, User{user3Id, "user3"})
}
