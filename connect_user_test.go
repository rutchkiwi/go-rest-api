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

func TestListConnections(t *testing.T) {
	buildWebservice()
	registerUser(t, "viktor", "pass")

	connections := listConnections(t)

	require.Len(t, connections, 0)
}

func listConnections(t *testing.T) []User {
	//TODO: its wierd that usernbame password are not constants
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
	buildWebservice()
	registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")

	addConnection(t, user2Id)
}

func addConnection(t *testing.T, toUserId int64) {
	bodyString := fmt.Sprintf(`{"id":%d}`, toUserId)
	fmt.Println(bodyString)
	bodyReader := strings.NewReader(bodyString)
	httpReq, _ := http.NewRequest("PUT", fmt.Sprint("/connection/", toUserId), bodyReader)
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	httpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, httpReq)

	require.Equal(t, 200, httpWriter.Code)
}

func TestAddAndListConnections(t *testing.T) {
	buildWebservice()
	registerUser(t, "viktor", "pass")
	user2Id := registerUser(t, "user2", "pass")

	addConnection(t, user2Id)

	connections := listConnections(t)
	assert.Len(t, connections, 1)
}
