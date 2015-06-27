package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/require"
)

func TestListConnections(t *testing.T) {
	buildWebservice()
	registerUser(t, "viktor", "pass")

	//GET /me
	getHttpReq, _ := http.NewRequest("GET", "/connection", nil)
	getHttpReq.Header.Set("Authorization", basicAuthEncode("viktor", "pass"))

	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.ServeHTTP(httpWriter, getHttpReq)

	require.Equal(t, 200, httpWriter.Code)

	var actual []User
	json.Unmarshal(httpWriter.Body.Bytes(), &actual)

	require.Len(t, actual, 0)
}
