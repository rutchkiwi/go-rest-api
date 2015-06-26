package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//http://golang-basic.blogspot.co.uk/2014/06/golang-database-step-by-step-guide-on.html

type Database struct {
	db *sql.DB
}

func newInMemoryDb() Database {
	db, err := sql.Open("sqlite3", "")
	checkErr(err)
	// Bootstrap
	sqlStmt := `
	create table if not exists user (
		id integer primary key autoincrement,
	 	username varchar(64) unique not null,
	 	password varchar(64) not null --YOLO!
	 	);
	`
	_, err = db.Exec(sqlStmt)
	checkErr(err)
	return Database{db}
	// defer db.Close() todo: neccessary?
}

type User struct {
	Username string
	//	id       int64
}

func (database Database) dbWriteNewUser(username, password string) (int64, User) {
	db := database.db
	// todo: make configurable, so test can use in memory db
	// sqlStmt := `
	// create table user (
	// 	id integer primary key autoincrement,
	//  	username varchar(64) unique not null,
	//  	password varchar(64) not null --YOLO!
	//  	);
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	// insert
	//TODO: password shouldnt be clear text
	stmt, err := db.Prepare("INSERT INTO user(username, password) VALUES(?,?)")
	checkErr(err)

	res, err := stmt.Exec(username, password)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id, User{username}
}

func (database Database) dbGetUser(id int64) (User, error) {
	db := database.db

	row := db.QueryRow("SELECT username FROM user WHERE id = ?", id)

	var username string
	if err := row.Scan(&username); err != nil {
		// TODO: bad error handling here
		return User{}, err
	}

	return User{username}, nil
}

func (database Database) dbGetPasswordForUsername(username string) (string, error) {
	//TODO: create index on username
	db := database.db

	row := db.QueryRow("SELECT password FROM user WHERE username=?", username)

	var password string
	err := row.Scan(&password)
	checkErr(err)
	//TODO: This needs to be handled in a more secure way
	// return "WTF", fmt.Errorf("could not find password in db")
	return password, nil

}

func checkErr(err error) {
	//TODO: handle
	if err != nil {
		panic(err)
	}
}
