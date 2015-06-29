package main

import (
	"net/http"
	"sort"

	"github.com/emicklei/go-restful"
)

var (
	database Database
)

func main() {
	buildWebservice(false)
	http.ListenAndServe(":8080", nil)
}

func buildWebservice(inMemoryDb bool) {
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
	ws.Route(ws.PUT("/connection").To(connectToUser))
	ws.Route(ws.GET("/admin/users").To(listAllUsersWithConnections))

	if inMemoryDb {
		database = newInMemoryDb()
	} else {
		database = newFileDb()
	}

	// Add admin user
	admin, _ := database.writeNewUser("admin", "pass")
	database.makeAdmin(admin.Id)

	restful.Add(ws)
}

func getUser(req *restful.Request, resp *restful.Response) {
	authenticatedUser, err := BasicAuth(req, database)
	if err != nil {
		unauthorized(resp)
		return
	}

	resp.WriteEntity(authenticatedUser.user)
}

type ErrorMsg struct {
	Message string
}

type UserRegistration struct {
	Username, Password string
}

func postUserRegistration(req *restful.Request, resp *restful.Response) {
	userRegistration := new(UserRegistration)
	err := req.ReadEntity(userRegistration)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.WriteEntity(ErrorMsg{err.Error()})
		return
	}

	// TODO: remove return vals if not needed
	var newUser User
	newUser, err = database.writeNewUser(userRegistration.Username, userRegistration.Password)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.WriteEntity(ErrorMsg{"Username " + userRegistration.Username + " is already taken"})
		return
	} else {
		resp.WriteEntity(newUser)
		return
	}
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
	authenticatedUser, err := BasicAuth(req, database)
	if err != nil {
		unauthorized(resp)
		return
	}

	connections := database.connectedUsers(authenticatedUser.user.Id)
	resp.WriteEntity(connections)
}

type ConnectionRequest struct {
	Id int64
}

func writeBadRequestMsg(resp *restful.Response, err error) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.WriteEntity(ErrorMsg{err.Error()})
}

// Takes a ConnectionRequest json as PUT body
// Kind of stupid, it would (maybe) be nicer to use the LINK verb (but its obscure)
// Or at least use PUT connection/<toId> with empty body. Sadly, Go-Restful cant handle
// empty PUT requests... Kind of regretting my choice of REST framework.
func connectToUser(req *restful.Request, resp *restful.Response) {
	authenticatedUser, err := BasicAuth(req, database)
	if err != nil {
		unauthorized(resp)
		return
	}

	var connectionRequest ConnectionRequest
	err = req.ReadEntity(&connectionRequest)
	if err != nil {
		writeBadRequestMsg(resp, err)
		return
	}

	err = database.addConnection(authenticatedUser.user.Id, connectionRequest.Id)
	if err != nil {
		writeBadRequestMsg(resp, err)
		return
	}
	// resp.WriteHeader(http.StatusOK)
	resp.WriteEntity("Connected or was already connected!") //TODO:
}

type UserWithConnections struct {
	User        User
	Connections []User
}

// Enables sorting by id in a list of UserWithConnections
type ById []UserWithConnections

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].User.Id < a[j].User.Id }

func listAllUsersWithConnections(req *restful.Request, resp *restful.Response) {
	//TODO: remove duplication of auth stuff
	//TODO: check admin privileges
	authenticatedUser, err := BasicAuth(req, database)
	if err != nil || !authenticatedUser.isAdmin {
		unauthorized(resp)
		return
	}

	allConnections := database.listAllConnections()

	res := make([]UserWithConnections, 0)
	for fromUser, toUsers := range allConnections {
		res = append(res, UserWithConnections{fromUser, toUsers})
	}

	// Go randomizes iteration order in maps, so we need to sort here
	// (since its nice for users to be sorted by id)
	sort.Sort(ById(res))

	resp.WriteEntity(res)
}
