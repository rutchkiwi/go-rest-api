package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//http://golang-basic.blogspot.co.uk/2014/06/golang-database-step-by-step-guide-on.html

type User struct {
	username string
	//	id       int64
}

func writeNewUser(username string) int64 {
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

	return id
}

func getUser(id int64) User {
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	defer db.Close()

	row := db.QueryRow("SELECT username FROM user WHERE id = ?", id)

	var username string
	if err := row.Scan(&username); err != nil {
		log.Fatal(err)
	}

	return User{username}
}

func checkErr(err error) {
	//TODO: handle
	if err != nil {
		panic(err)
	}
}
