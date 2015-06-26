package main

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

// "lib/db"

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

var (
	database Database
)

func main() {
	buildWebservice()
	http.ListenAndServe(":8080", nil)
}

func buildWebservice() {
	ws := new(restful.WebService)

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	ws.Route(ws.GET("/me").To(getUser))
	ws.Route(ws.POST("/register").To(postUser))

	database = newInMemoryDb()

	restful.Add(ws)
}

func getUser(req *restful.Request, resp *restful.Response) {
	id, err := BasicAuth(req)
	if err != nil {
		unauthorized(resp)
		return
	}

	// fmt.Println("getting user")
	// fmt.Printf("getuser by string id %v\n", req.PathParameter("user-id"))

	// id, err := strconv.ParseInt(req.PathParameter("user-id"), 0, 64)
	// if err != nil {
	// 	// fmt.Println(err)
	// 	//TODO: bad error handling here
	// 	resp.WriteHeader(http.StatusBadRequest)
	// 	resp.WriteEntity("malformed id") //TODO: json?	}
	// 	return
	// }

	// fmt.Printf("getuser by id %v\n", id)
	user, err := database.dbGetUser(id)
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

type UserRegistration struct {
	Username, Password string
}

func postUser(req *restful.Request, resp *restful.Response) {
	userRegistration := new(UserRegistration)
	err := req.ReadEntity(userRegistration)
	checkErr(err)

	// TODO: remove return vals if not needed
	database.dbWriteNewUser(userRegistration.Username, userRegistration.Password)
	resp.WriteHeader(http.StatusOK)
	resp.WriteEntity("success! See /me to see your info")
}
