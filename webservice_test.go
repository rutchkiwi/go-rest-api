package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
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
	getHttpReq, _ := http.NewRequest("GET", "9999", nil)
	getReq := restful.NewRequest(getHttpReq)

	getRecorder := new(httptest.ResponseRecorder)
	getResp := restful.NewResponse(getRecorder)

	getUser(getReq, getResp)

	assert.Equal(t, 404, getRecorder.Code)

}

func TestPostAndGetUSer(t *testing.T) {

	// bodyReader := strings.NewReader("<Sample><Value>42</Value></Sample>")
	// 	httpRequest, _ := http.NewRequest("GET", "/test/THIS", bodyReader)
	// 	httpRequest.Header.Set("Content-Type", restful.MIME_XML)
	// 	httpWriter := httptest.NewRecorder()

	bodyReader := strings.NewReader(`{"username":"viktor"}`)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)
	req := restful.NewRequest(httpRequest)

	recorder := new(httptest.ResponseRecorder)
	resp := restful.NewResponse(recorder)

	postUser(req, resp)

	assert.Equal(t, 201, recorder.Code)
	var location string
	location = recorder.Header().Get("location")

	getHttpReq, _ := http.NewRequest("GET", location, nil)
	getReq := restful.NewRequest(getHttpReq)

	getRecorder := new(httptest.ResponseRecorder)
	getResp := restful.NewResponse(getRecorder)

	postUser(getReq, getResp)

	assert.Equal(t, 200, getRecorder.Code)

}
