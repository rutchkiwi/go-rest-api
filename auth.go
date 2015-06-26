package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
)

type LoggedInUser struct {
	id       int64
	username string
}

func BasicAuth(req *restful.Request) (int64, error) {
	auth := req.HeaderParameter("Authorization")
	if len(auth) < 6 || auth[:6] != "Basic " {
		return -1, fmt.Errorf("invalid basic auth syntax header")
	}
	b, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return -1, err
	}
	tokens := strings.SplitN(string(b), ":", 2)
	if len(tokens) != 2 {
		return -1, err
	}
	loggedInUser, err := authfn(tokens[0], tokens[1])
	if err != nil {
		return -1, err
	}

	//TODO: actual id!
	return loggedInUser.id, nil
}

func authfn(username, password string) (LoggedInUser, error) {
	//TODO: hook to db
	if password != "guessme" {
		return LoggedInUser{}, fmt.Errorf("Invalid credentials")
	} else {
		return LoggedInUser{1, username}, nil
	}
}

func unauthorized(resp *restful.Response) {
	resp.WriteHeader(http.StatusUnauthorized)
	resp.AddHeader("WWW-Authenticate", "Basic realm=\""+"realm"+"\"")
	resp.WriteEntity("Unauthorized")
}
