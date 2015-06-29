package main

import (
	"crypto/subtle"
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

// Call this to check which, if any, user authenticating.
// It's a bit bad that you need to check the results yourself in the REST
// endpoints.. We could implement this as a function wrapper as well, but this
// might be confusing instead.

//TODO: Basic auth without SSL/TLS = 😱
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

	return loggedInUser, nil
}

func authfn(username, givenPassword string, db Database) (AuthenticatedUser, error) {
	userWithPassword, err := db.getUserAndPasswordForUsername(username)
	if err != nil {
		return AuthenticatedUser{}, fmt.Errorf("Invalid credentials")
	}
	// Use constant time compare to protect against timing attacks
	passwordsMatch := subtle.ConstantTimeCompare(
		[]byte(givenPassword),
		[]byte(*(userWithPassword.password))) == 1
	if passwordsMatch {
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
