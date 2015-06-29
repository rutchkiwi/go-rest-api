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
	//TODO: fix password storage
	sqlStmt := `
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

	PRAGMA foreign_keys = ON;
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

//TODO: remove the db prefix everywhere
func (database Database) dbWriteNewUser(username, password string) (User, error) {
	db := database.db

	//TODO: merge stmt and exec
	stmt, err := db.Prepare("INSERT INTO user(username, password) VALUES(?,?)")
	checkErr(err)

	res, err := stmt.Exec(username, password)
	if err != nil {
		return User{}, err //TODO: should filter to only right error
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

type UserWithPassword struct {
	user     User
	password *string
	isAdmin  bool
}

// Returns password as a *string, so that it can be Nil (otherwise we'd hade to return "", which could
// cause security holes when comaring it to other given password strings)
func (database Database) dbGetUserAndPasswordForUsername(username string) (UserWithPassword, error) {
	//TODO: create index on username
	db := database.db

	//TODO: handle when this doesnt find anything
	row := db.QueryRow("SELECT id, password, admin FROM user WHERE username=?", username)

	var id int64
	var password string
	var isAdmin bool
	err := row.Scan(&id, &password, &isAdmin)

	if err != nil {
		return UserWithPassword{}, err
	}

	//TODO: This needs to be handled in a more secure way
	// return "WTF", fmt.Errorf("could not find password in db")
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

func (database Database) addConnection(from, to int64) {
	//TODO: handle trying to connect to user who doesnt exist
	_, err := database.db.Exec(`
		INSERT OR IGNORE INTO connection(fromUser, toUser) VALUES(?,?)`,
		from, to)
	checkErr(err)
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
	//TODO: handle
	if err != nil {
		panic(err)
	}
}
