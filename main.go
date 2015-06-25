package main

import (
	"fmt"
	"net/http"
	"strconv"

	// "lib/db"

	"github.com/emicklei/go-restful"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

// ws.Route(ws.GET("/{user-id}").To(u.findUser))  // u is a UserResource

// ...

// // GET http://localhost:8080/users/1
// func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
// 	id := request.PathParameter("user-id")
// 	...
// }

func main() {
	buildWebservice()
	http.ListenAndServe(":8080", nil)
}

func buildWebservice() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/{user-id}").To(getUser))
	ws.Route(ws.POST("/").To(postUser))
	restful.Add(ws)
}

func getUser(req *restful.Request, resp *restful.Response) {
	// fmt.Println("getting user")
	restful.DefaultResponseContentType(restful.MIME_JSON)
	// fmt.Printf("getuser by string id %v\n", req.PathParameter("user-id"))

	id, err := strconv.ParseInt(req.PathParameter("user-id"), 0, 64)
	if err != nil {
		// fmt.Println(err)
		//TODO: bad error handling here
		resp.WriteHeader(http.StatusBadRequest)
		resp.WriteEntity("malformed id") //TODO: json?	}
		return
	}

	// fmt.Printf("getuser by id %v\n", id)
	user, err := dbGetUser(id)
	if err != nil {
		fmt.Println(err)
		//TODO: bad error handling here
		resp.WriteHeader(http.StatusNotFound)
		resp.WriteEntity("no such user") //TODO: json?
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.WriteEntity(user)
	}
}

func postUser(req *restful.Request, resp *restful.Response) {
	inputUser := new(User)
	err := req.ReadEntity(inputUser)
	checkErr(err)

	restful.DefaultResponseContentType(restful.MIME_JSON) //TODO: move this
	id, user := dbWriteNewUser(inputUser.Username)
	resp.WriteHeader(http.StatusCreated)
	resp.AddHeader("location", strconv.FormatInt(id, 10))
	resp.WriteEntity(user)
}
