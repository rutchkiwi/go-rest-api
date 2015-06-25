package main

import (
	"io"
	"net/http"

	"github.com/emicklei/go-restful"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)
	http.ListenAndServe(":8080", nil)
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}

func postUser(req *restful.Request, resp *restful.Response) {
	restful.DefaultResponseContentType(restful.MIME_JSON)
	resp.AddHeader("location", "todo")
	resp.WriteHeader(http.StatusCreated)
	user := User{"todo"}
	resp.WriteEntity(user)
}
