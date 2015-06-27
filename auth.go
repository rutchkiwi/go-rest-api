package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
)

func BasicAuth(req *restful.Request, db Database) (User, error) {
	auth := req.HeaderParameter("Authorization")
	if len(auth) < 6 || auth[:6] != "Basic " {
		return User{}, fmt.Errorf("invalid basic auth syntax header")
	}
	b, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return User{}, err
	}
	tokens := strings.SplitN(string(b), ":", 2)
	if len(tokens) != 2 {
		return User{}, err
	}
	loggedInUser, err := authfn(tokens[0], tokens[1], database)
	if err != nil {
		return User{}, err
	}

	//TODO: actual id!
	return loggedInUser, nil
}

func authfn(username, givenPassword string, db Database) (User, error) {
	user, actualPassword, err := db.dbGetUserAndPasswordForUsername(username)
	//TODO: insercure that we return an empty user? (easy to mess up)
	if err != nil {
		return User{}, fmt.Errorf("Invalid credentials")
	}
	if givenPassword == *actualPassword { //TODO: secure compare
		return user, nil
	} else {
		return User{}, fmt.Errorf("Invalid credentials")
	}
}

func unauthorized(resp *restful.Response) {
	resp.WriteHeader(http.StatusUnauthorized)
	resp.AddHeader("WWW-Authenticate", "Basic realm=\""+"realm"+"\"")
	resp.WriteEntity("Unauthorized")
}
