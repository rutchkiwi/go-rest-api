package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNonExistingUser(t *testing.T) {
	buildWebservice()

	getHttpReq, _ := http.NewRequest("GET", "/9999", nil)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	assert.Equal(t, 404, httpWriter.Code)
	body := httpWriter.Body.String()
	assert.Equal(t, `"no such user"`, body)
}

func TestPostUser(t *testing.T) {
	buildWebservice()

	bodyReader := strings.NewReader(`{"username":"viktor"}`)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	assert.Equal(t, 201, httpWriter.Code)
	location := httpWriter.Header().Get("location")
	assert.Regexp(t, `\d+`, location)

	body := httpWriter.Body.String()
	assert.Equal(t, `{"username":"viktor"}`, body)
}

func TestPostAndGetUSer(t *testing.T) {

	buildWebservice()

	bodyReader := strings.NewReader(`{"username":"viktor"}`)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	require.Equal(t, 201, httpWriter.Code)

	var location string
	location = httpWriter.Header().Get("location")

	getHttpReq, err := http.NewRequest("GET", "/"+location, nil)
	require.NoError(t, err, "error in test")

	httpWriter = httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)
	assert.Equal(t, 200, httpWriter.Code)
	body := httpWriter.Body.String()
	assert.Equal(t, `{"username":"viktor"}`, body)

}
