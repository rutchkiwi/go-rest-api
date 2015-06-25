package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This example show how to test one particular RouteFunction (getIt)
// It uses the httptest.ResponseRecorder to capture output

// func TestPostUser(t *testing.T) {
// 	httpReq, _ := http.NewRequest("POST", "/", nil)
// 	getReq := restful.NewRequest(httpReq)

// 	recorder := new(httptest.ResponseRecorder)
// 	resp := restful.NewResponse(recorder)

// 	postUser(req, resp)

// 	assert.Equal(t, 201, recorder.Code)
// }

func TestGetNonExistingUser(t *testing.T) {
	buildWebservice()

	getHttpReq, _ := http.NewRequest("GET", "/9999", nil)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	fmt.Println(httpWriter)
	assert.Equal(t, 404, httpWriter.Code)

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
	fmt.Println("location:")
	fmt.Println("/" + location)

	getHttpReq, err := http.NewRequest("GET", "/"+location, nil)
	require.NoError(t, err, "error in test")

	httpWriter = httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)
	assert.Equal(t, 200, httpWriter.Code)

}
