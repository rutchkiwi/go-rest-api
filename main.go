package main

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

var (
	database Database
)

func main() {
	buildWebservice()
	http.ListenAndServe(":8080", nil)
}

func buildWebservice() {
	ws := new(restful.WebService)

	restful.EnableTracing(true)

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	ws.Route(ws.GET("/me").To(getUser))
	ws.Route(ws.POST("/register").To(postUserRegistration))
	ws.Route(ws.GET("/search").To(searchForUsers))
	ws.Route(ws.GET("/connection").To(listConnectedUsers))
	ws.Route(ws.PUT("/connection/{id}").To(connectToUser))

	database = newInMemoryDb()

	restful.Add(ws)
}

func getUser(req *restful.Request, resp *restful.Response) {
	user, err := BasicAuth(req, database)
	if err != nil {
		unauthorized(resp)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.WriteEntity(user)
}

type UserRegistration struct {
	Username, Password string
}

func postUserRegistration(req *restful.Request, resp *restful.Response) {
	userRegistration := new(UserRegistration)
	//TODO: invalid input causes panic
	err := req.ReadEntity(userRegistration)
	checkErr(err)

	// TODO: remove return vals if not needed
	var newUser User
	newUser = database.dbWriteNewUser(userRegistration.Username, userRegistration.Password)
	resp.WriteEntity(newUser)
}

type SearchResults struct {
	Results []User
}

func searchForUsers(req *restful.Request, resp *restful.Response) {
	q := req.QueryParameter("q")
	users := database.searchForUsers(q)
	resp.WriteEntity(SearchResults{users})
}

func listConnectedUsers(req *restful.Request, resp *restful.Response) {
	user, err := BasicAuth(req, database)
	if err != nil {
		unauthorized(resp)
		return
	}

	connections := database.connectedUsers(user.Id)
	resp.WriteEntity(connections)
}

type ConnectionRequest struct {
	Id int64
}

func connectToUser(req *restful.Request, resp *restful.Response) {
	// loggedInUser, err := BasicAuth(req, database)
	// if err != nil {
	// 	unauthorized(resp)
	// 	return
	// }
	from, err := BasicAuth(req, database)
	if err != nil {
		unauthorized(resp)
		return
	}

	var connectionRequest ConnectionRequest
	err = req.ReadEntity(&connectionRequest)
	if err != nil {
		checkErr(err) //TODO:
	}

	database.addConnection(from.Id, connectionRequest.Id)

	// resp.WriteHeader(http.StatusOK)
	resp.WriteEntity("Connected or was already connected!") //TODO:
}
