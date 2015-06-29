package main

import (
	"encoding/base64"
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

func TestGetUserMeNoAuth(t *testing.T) {
	buildWebservice(true)

	getHttpReq, _ := http.NewRequest("GET", "/me", nil)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	assert.Equal(t, 401, httpWriter.Code)
}

func basicAuthEncode(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func TestGetUserMeWrongUsername(t *testing.T) {
	buildWebservice(true)

	getHttpReq, _ := http.NewRequest("GET", "/me", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("nonExistingUser", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	assert.Equal(t, 401, httpWriter.Code)
}

func TestPostUser(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")
}

func TestPostUserTwice(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")

	bodyString := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, "viktor", "pass")
	bodyReader := strings.NewReader(bodyString)
	httpRequest, _ := http.NewRequest("POST", "/register", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)
	require.Equal(t, 400, httpWriter.Code)

	var errorMsg ErrorMsg
	json.Unmarshal(httpWriter.Body.Bytes(), &errorMsg)
	assert.Equal(t, "Username viktor is already taken", errorMsg.Message)
}

func registerUser(t *testing.T, username, password string) int64 {
	bodyString := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	bodyReader := strings.NewReader(bodyString)
	httpRequest, _ := http.NewRequest("POST", "/register", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	require.Equal(t, 200, httpWriter.Code)

	var newUser User
	json.Unmarshal(httpWriter.Body.Bytes(), &newUser)
	assert.Equal(t, username, newUser.Username)
	assert.True(t, newUser.Id > 0)

	return newUser.Id
}

func TestRegisterUserBadInput(t *testing.T) {
	bodyString := fmt.Sprintf(`sdfsdfs`)
	bodyReader := strings.NewReader(bodyString)
	httpRequest, _ := http.NewRequest("POST", "/register", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	assert.Equal(t, 400, httpWriter.Code)
}

func TestGetMe(t *testing.T) {
	buildWebservice(true)
	registerUser(t, "viktor", "pass")
	registerUser(t, "user2", "pass")

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
	buildWebservice(true)
	registerUser(t, "viktor", "pass")

	//GET /me
	getHttpReq, _ := http.NewRequest("GET", "/me", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("viktor", "wrongPass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 401, httpWriter.Code)
}
