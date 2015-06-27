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
	//TODO: you need to use real one in the actual app
	db, err := sql.Open("sqlite3", "")
	checkErr(err)
	// Bootstrap
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	username VARCHAR(64) UNIQUE NOT NULL,
	 	password VARCHAR(64) NOT NULL --YOLO!
	 	);
	CREATE TABLE "connections" (
		"from" INTEGER NOT NULL REFERENCES "user"("id") ON UPDATE CASCADE ON DELETE CASCADE,
		"to"   INTEGER NOT NULL REFERENCES "user"("id") ON UPDATE CASCADE ON DELETE CASCADE,
		PRIMARY KEY ("from", "to")
	);
	`
	_, err = db.Exec(sqlStmt)
	checkErr(err)
	return Database{db}
	// defer db.Close() todo: neccessary?
}

type User struct {
	Id       int64
	Username string
}

func (database Database) dbWriteNewUser(username, password string) User {
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

	return User{id, username}
}

func (database Database) dbGetUser(id int64) (User, error) {
	db := database.db

	row := db.QueryRow("SELECT username FROM user WHERE id = ?", id)

	var username string
	if err := row.Scan(&username); err != nil {
		// TODO: bad error handling here
		return User{}, err
	}

	return User{id, username}, nil
}

// Returns password as a *string, so that it can be Nil (otherwise we'd hade to return "", which could
// cause security holes when comaring it to other given password strings)
func (database Database) dbGetUserAndPasswordForUsername(username string) (User, *string, error) {
	//TODO: create index on username
	db := database.db

	//TODO: handle when this doesnt find anything
	row := db.QueryRow("SELECT id, password FROM user WHERE username=?", username)

	var id int64
	var password string
	err := row.Scan(&id, &password)

	if err != nil {
		return User{}, nil, err
	}

	//TODO: This needs to be handled in a more secure way
	// return "WTF", fmt.Errorf("could not find password in db")
	return User{id, username}, &password, nil
}

func (database Database) searchForUsers(query string) []User {
	res := make([]User, 0)

	query = "%" + query + "%"
	rows, err := database.db.Query("SELECT id, username FROM user WHERE username LIKE ?", query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username)
		checkErr(err)
		res = append(res, user)
	}

	return res
}

func (d Database) connectedUsers(userId int64) []User {
	res := make([]User, 0)

	// rows, err := database.db.Query("SELECT id, username FROM user WHERE username LIKE ?", query)
	// checkErr(err)
	// defer rows.Close()
	// for rows.Next() {
	// 	var user User
	// 	err := rows.Scan(&user.Id, &user.Username)
	// 	checkErr(err)
	// 	res = append(res, user)
	// }

	return res
}

func checkErr(err error) {
	//TODO: handle
	if err != nil {
		panic(err)
	}
}
