package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserMeNoAuth(t *testing.T) {
	buildWebservice()

	getHttpReq, _ := http.NewRequest("GET", "/me", nil)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	assert.Equal(t, 401, httpWriter.Code)
}

func basicAuthEncode(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func TestGetUserMeWrongUsername(t *testing.T) {
	buildWebservice()

	getHttpReq, _ := http.NewRequest("GET", "/me", nil)
	//TODO: move into method

	getHttpReq.Header.Set("Authorization", basicAuthEncode("nonExistingUser", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	assert.Equal(t, 401, httpWriter.Code)
}

func TestPostUser(t *testing.T) {
	buildWebservice()
	registerUser(t)

	// var user User
	// json.Unmarshal(httpWriter.Body.Bytes(), &user)
	// assert.Equal(t, user, User{"viktor"})
}

func registerUser(t *testing.T) {
	bodyReader := strings.NewReader(`{"username":"viktor", "password":"pass"}`)
	httpRequest, _ := http.NewRequest("POST", "/register", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	require.Equal(t, 200, httpWriter.Code)
}

func TestGetMe(t *testing.T) {

	buildWebservice()
	registerUser(t)

	//GET /me
	getHttpReq, _ := http.NewRequest("GET", "/me", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 200, httpWriter.Code)

	var user User
	json.Unmarshal(httpWriter.Body.Bytes(), &user)
	assert.Equal(t, "viktor", user.Username)
}

func TestGetMeWithWrongPassword(t *testing.T) {

	buildWebservice()
	registerUser(t)

	//GET /me
	getHttpReq, _ := http.NewRequest("GET", "/me", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("viktor", "wrongPass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 401, httpWriter.Code)
}
