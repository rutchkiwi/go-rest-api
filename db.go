package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//http://golang-basic.blogspot.co.uk/2014/06/golang-database-step-by-step-guide-on.html

type User struct {
	username string
	//	id       int64
}

func dbWriteNewUser(username string) int64 {
	// todo: make configurable, so test can use in memory db
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	defer db.Close()

	// insert
	stmt, err := db.Prepare("INSERT INTO user(username) VALUES(?)")
	checkErr(err)

	res, err := stmt.Exec(username)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Printf("put db with id %v\n", id)

	return id
}

func dbGetUser(id int64) (User, error) {
	fmt.Printf("getting db with id %v\n", id)
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	defer db.Close()

	row := db.QueryRow("SELECT username FROM user WHERE id = ?", id)

	var username string
	if err := row.Scan(&username); err != nil {
		// TODO: bad error handling here
		return User{}, err
	}

	return User{username}, nil
}

func checkErr(err error) {
	//TODO: handle
	if err != nil {
		panic(err)
	}
}
