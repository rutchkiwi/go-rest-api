package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
)

type AuthenticatedUser struct {
	user    User
	isAdmin bool
}

func BasicAuth(req *restful.Request, db Database) (AuthenticatedUser, error) {
	auth := req.HeaderParameter("Authorization")
	if len(auth) < 6 || auth[:6] != "Basic " {
		return AuthenticatedUser{}, fmt.Errorf("invalid basic auth syntax header")
	}
	b, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return AuthenticatedUser{}, err
	}
	tokens := strings.SplitN(string(b), ":", 2)
	if len(tokens) != 2 {
		return AuthenticatedUser{}, err
	}
	loggedInUser, err := authfn(tokens[0], tokens[1], database)
	if err != nil {
		return AuthenticatedUser{}, err
	}

	//TODO: actual id!
	return loggedInUser, nil
}

func authfn(username, givenPassword string, db Database) (AuthenticatedUser, error) {
	userWithPassword, err := db.getUserAndPasswordForUsername(username)
	//TODO: insercure that we return an empty user? (easy to mess up)
	if err != nil {
		return AuthenticatedUser{}, fmt.Errorf("Invalid credentials")
	}
	if givenPassword == *(userWithPassword.password) { //TODO: secure compare
		return AuthenticatedUser{userWithPassword.user, userWithPassword.isAdmin}, nil
	} else {
		return AuthenticatedUser{}, fmt.Errorf("Invalid credentials")
	}
}

func unauthorized(resp *restful.Response) {
	resp.WriteHeader(http.StatusUnauthorized)
	resp.AddHeader("WWW-Authenticate", "Basic realm=\""+"realm"+"\"")
	resp.WriteEntity("Unauthorized")
}
