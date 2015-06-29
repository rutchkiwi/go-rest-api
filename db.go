package main

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//http://golang-basic.blogspot.co.uk/2014/06/golang-database-step-by-step-guide-on.html

type Database struct {
	db *sql.DB
}

func newFileDb() Database {
	db, err := sql.Open("sqlite3", "app.db")
	checkErr(err)
	database := Database{db}
	database.bootstrap()
	return database
}

func newInMemoryDb() Database {
	db, err := sql.Open("sqlite3", "")
	checkErr(err)
	database := Database{db}
	database.bootstrap()
	return database
}

func (database Database) bootstrap() {
	// Bootstrap db schema

	// TODO: fix password storage
	// Should be properly hashed/salted..
	// maybe using http://godoc.org/golang.org/x/crypto/bcrypt
	// (but doing this in the proper way would take a while)
	sqlStmt := `
	PRAGMA foreign_keys = ON;
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	username VARCHAR(64) UNIQUE NOT NULL,
	 	password VARCHAR(64) NOT NULL, --YOLO!
	 	admin BOOLEAN DEFAULT FALSE NOT NULL
	 	);
	CREATE TABLE IF NOT EXISTS connection (
		fromUser INTEGER NOT NULL REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE,
		toUser   INTEGER NOT NULL REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE,
		PRIMARY KEY (fromUser, toUser)
	);
	CREATE INDEX IF NOT EXISTS usernameIndex ON user(username);
	`
	_, err := database.db.Exec(sqlStmt)
	checkErr(err)
}

type User struct {
	Id       int64
	Username string
}

func (database Database) writeNewUser(username, password string) (User, error) {
	db := database.db

	res, err := db.Exec("INSERT INTO user(username, password) VALUES(?,?)", username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return User{}, errors.New("User " + username + " already exists")
		} else {
			panic(err)
		}

	}

	id, err := res.LastInsertId()
	checkErr(err)

	return User{id, username}, nil
}

func (database Database) makeAdmin(userId int64) {
	// sqllite does not have boolean literals, hence the "= 1"
	_, err := database.db.Exec("UPDATE user SET admin = 1 WHERE id = ?", userId)
	checkErr(err)
}

func (database Database) getUser(id int64) (User, error) {
	db := database.db

	row := db.QueryRow("SELECT username FROM user WHERE id = ?", id)

	var username string
	err := row.Scan(&username)
	checkErr(err)

	return User{id, username}, nil
}

type UserWithPassword struct {
	user     User
	password *string
	isAdmin  bool
}

// Returns password as a *string, so that it can be Nil (otherwise we'd hade to return "", which
// would be error prone when comaring it to other given password strings)
func (database Database) getUserAndPasswordForUsername(username string) (UserWithPassword, error) {
	db := database.db

	row := db.QueryRow("SELECT id, password, admin FROM user WHERE username=?", username)

	var id int64
	var password string
	var isAdmin bool
	err := row.Scan(&id, &password, &isAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			return UserWithPassword{}, err
		} else {
			panic(err)
		}
	}

	return UserWithPassword{User{id, username}, &password, isAdmin}, nil
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

func (database Database) connectedUsers(userId int64) []User {
	res := make([]User, 0)

	rows, err := database.db.Query(`
		SELECT user.id, user.username FROM 
			connection JOIN user ON user.id = connection.toUser 
			WHERE connection.fromUser = ?`, userId)
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

func (database Database) addConnection(from, to int64) error {
	_, err := database.db.Exec(`
		INSERT OR IGNORE INTO connection(fromUser, toUser) VALUES(?,?)`,
		from, to)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return errors.New("No such user")
		} else {
			panic(err)
		}
	}
	return nil
}

type DbUsersWithConnections map[User][]User

func (database Database) listAllConnections() (res DbUsersWithConnections) {
	res = make(DbUsersWithConnections)

	rows, err := database.db.Query(`
		-- Left join because we want users without connections as well
		SELECT u1.id, u1.username, u2.id, u2.username FROM 
			user AS u1 LEFT JOIN connection ON u1.id = connection.fromUser
			LEFT JOIN user AS u2 ON u2.id = connection.toUser
			ORDER BY u1.id
			`)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var fromUser User
		var toUsername sql.NullString
		var toId sql.NullInt64
		err := rows.Scan(&fromUser.Id, &fromUser.Username, &toId, &toUsername)
		checkErr(err)

		if toId.Valid {
			// this user has at least one connection, unpack the nullable values
			toIdValue, _ := toId.Value()
			toUsernameValue, _ := toUsername.Value()
			res[fromUser] = append(res[fromUser], User{toIdValue.(int64), toUsernameValue.(string)})
		} else {
			// this user doesn't have any connections
			res[fromUser] = []User{}
		}
	}
	return res
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
