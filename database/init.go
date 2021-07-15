package database

import (
	"database/sql"
)

type DBUser struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TableMethods interface {
	createTable () error
	InsertUsers(string, string) error
	GetUsers() ([]*DBUser, error)
	GetUserByID(int) (*DBUser, error)
	GetUserPasswordByUsername(string)(string, error)
	UpdateUser(int, string, string) error
	DeleteUser(int)error
}

type DB struct {
	DBConn *sql.DB
}

func NewDB() (*DB, error) {	//this is a constructor to create a database
	dbConn, err := sql.Open("sqlite3", "./users-database.db")
	if err != nil {
		return nil, err
	}

	newDBConnection := &DB{
		DBConn: dbConn,
	}

	//this is to prevent empty database when created
	if err = newDBConnection.createTable(); err != nil {
		return nil, err
	}

	return newDBConnection, nil
}

func (db *DB) createTable () error {
	statement, err := db.DBConn.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,username TEXT,password TEXT)")
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertUsers(username string, password string) error {
	statement, err := db.DBConn.Prepare("INSERT INTO users (username, password) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(username, password)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUsers() ([]*DBUser, error) {
	dbUserRows := make([]*DBUser, 0)
	row, err := db.DBConn.Query("SELECT * FROM users")
	if err != nil {
		return nil, nil
	}

	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			//return
		}
	}(row)

	for row.Next() { // Iterate and fetch the records from result cursor
		dbUser := &DBUser{}
		if err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password); err != nil {
			return nil, err
		}
		dbUserRows = append(dbUserRows, dbUser)
	}
	return dbUserRows, nil
}

func (db *DB) GetUserByID(id int) (*DBUser, error) {
	row := db.DBConn.QueryRow("SELECT * FROM users WHERE id=?", id)
	//Create a DBUser instance to store every scanned value into each of the variable
	dbUser := &DBUser{}
	if err := row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password); err != nil {
		return nil, err
	}
	return dbUser, nil
}

func (db *DB) GetUserPasswordByUsername (username string)(*string, error){
	row := db.DBConn.QueryRow("SELECT password FROM users WHERE username=?", username)
	dbUser := &DBUser{}
	if err := row.Scan(&dbUser.Password); err != nil {
		return nil, err
	}
	return &dbUser.Password, nil
}

func (db *DB) UpdateUser(id int, username string, password string) error{
	targetRow, err := db.GetUserByID(id)
	if err != nil {
		return err
	}
	if len(username) > 0 && targetRow.Username != username{
		targetRow.Username = username
	}
	if len(password) > 0 && targetRow.Password != password{
		targetRow.Password = password
	}
	statement, err := db.DBConn.Prepare("UPDATE users SET username = ? , password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(targetRow.Username, targetRow.Password, targetRow.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteUser(id int) error{
	targetRow, err := db.GetUserByID(id)
	if err != nil {
		return err
	}
	statement, err := db.DBConn.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(targetRow.ID)
	if err != nil {
		return err
	}
	return nil
}

//
//func CreateTable (database *sql.DB){
//	//Create table users if doest not exists
//	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,firstname TEXT,lastname TEXT)")
//	if err != nil {
//		fmt.Println(err)
//	}
//	exec, err := statement.Exec()
//	if err != nil {
//		return
//	}
//}
//
//func InsertUsers(database *sql.DB, firstname string, lastname string)  {
//	//Insert record into table users
//	statement, err := database.Prepare("INSERT INTO users (firstname, lastname) VALUES (?,?)")
//	if err != nil {
//		log.Fatalln(err.Error())
//	}
//	_, err = statement.Exec(firstname, lastname)
//	if err != nil {
//		log.Fatalln(err.Error())
//	}
//}
//
//func DisplayUsers(database *sql.DB) {
//	row, err := database.Query("SELECT * FROM users")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer row.Close()
//	for row.Next() { // Iterate and fetch the records from result cursor
//		var id int
//		var firstname string
//		var lastname string
//		row.Scan(&id, &firstname, &lastname)
//		fmt.Println(strconv.Itoa(id) + ":" + firstname + " " + lastname)
//	}


//defer func(database *sql.DB) {
//	err := database.Close()
//	if err != nil {
//		fmt.Println(err)
//	}
//}(database)