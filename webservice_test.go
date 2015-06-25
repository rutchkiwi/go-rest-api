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

	// bodyReader := strings.NewReader("<Sample><Value>42</Value></Sample>")
	// 	httpRequest, _ := http.NewRequest("GET", "/test/THIS", bodyReader)
	// 	httpRequest.Header.Set("Content-Type", restful.MIME_XML)
	// 	httpWriter := httptest.NewRecorder()

	bodyReader := strings.NewReader(`{"username":"viktor"}`)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)
	// req := restful.NewRequest(httpRequest)

	// recorder := new(httptest.ResponseRecorder)
	// resp := restful.NewResponse(recorder)

	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	require.Equal(t, 201, httpWriter.Code)

	var location string
	location = httpWriter.Header().Get("location")
	fmt.Println("location:")
	fmt.Println("/" + location)

	getHttpReq, err := http.NewRequest("GET", "/"+location, nil)
	require.NoError(t, err, "test error")
	// getReq := restful.NewRequest(getHttpReq)

	fmt.Println("222")

	httpWriter = httptest.NewRecorder()

	fmt.Println("3333")
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)
	fmt.Println("4444")

	assert.Equal(t, 200, httpWriter.Code)

}
