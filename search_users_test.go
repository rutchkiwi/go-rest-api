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

func TestSearchNoHits(t *testing.T) {

	buildWebservice()
	registerUser(t, "viktor", "pass")

	//GET /me
	getHttpReq, _ := http.NewRequest("GET", "/search?q=wierdQuery", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 200, httpWriter.Code)

	expected := SearchResults{make([]User, 0)}

	var actual SearchResults
	json.Unmarshal(httpWriter.Body.Bytes(), &actual)

	assert.Equal(t, expected, actual)
}

func TestSearch(t *testing.T) {

	buildWebservice()
	registerUser(t, "viktor", "pass")
	registerUser(t, "user2", "pass")

	//GET /me
	getHttpReq, _ := http.NewRequest("GET", "/search?q=user", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 200, httpWriter.Code)

	var actual SearchResults
	json.Unmarshal(httpWriter.Body.Bytes(), &actual)

	require.Len(t, actual.Results, 1)
	assert.Equal(t, "user2", actual.Results[0].Username)
	assert.True(t, actual.Results[0].Id > 0)
}
