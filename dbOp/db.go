package dbop

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"kivancaydogmus.com/apps/userApp/middleware"
)

type Person struct {
	Id       int    `json:"PersonID"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	Token    string `json:"Token"`
}

type Token struct {
	OwnerId int    `json:"OwnerID"`
	Context string `json:"Context"`
}

type Todo struct {
	OwnerId int    `json:"OwnerID"`
	Context string `json:"Context"`
}

const (
	username = "username"
	password = "password"
	hostname = "127.0.0.1:3306"
	dbname   = "article"
)

var db *sql.DB

func init() {
	db = prepareDb(dbname)
	defer db.Close()

}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func connectToDb(db *sql.DB) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("error % s when creating db\n", err)
		return
	}
	no, err := res.RowsAffected()

	if err != nil {
		log.Printf("error %s when fetching db\n", err)
	}
	log.Printf("rows affected %d\n", no)
}

func prepareDb(dbname string) *sql.DB {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("error %s during the open db\n", err)
	}
	connectToDb(db)
	return db
}

func GetAllUsers() []Person {
	db := prepareDb(dbname)
	defer db.Close()
	results, err := db.Query("select * from Persons")
	if err != nil {
		panic(err.Error())
	}
	var personArr []Person
	for results.Next() {
		var temp Person
		err = results.Scan(&temp.Id, &temp.UserName, &temp.Password, &temp.Token)
		if err != nil {
			panic(err.Error())
		}
		GetLastLoginToken(temp.UserName)
		personArr = append(personArr, temp)
	}
	return personArr
}

func GetTodo(username string) []string {
	db := prepareDb(dbname)
	defer db.Close()
	person := getPersonToDelToken(username)
	res, err := db.Query("SELECT * FROM Todo WHERE OwnerID = ?", person.Id)
	if err != nil {
		log.Fatal("an error occcured during the get todo from db ", err)
	}
	var todos []string

	for res.Next() {
		var tempTodo Todo
		err := res.Scan(&tempTodo.OwnerId, &tempTodo.Context)
		if err != nil {
			log.Fatal("an error occured during the scan db to get todo ", err)
		}
		todos = append(todos, tempTodo.Context)
	}
	return todos
}

func getIdByName(username string) int {
	db := prepareDb(dbname)
	defer db.Close()
	results, err := db.Query("select PersonID from Persons WHERE UserName = ?", username)
	if err != nil {
		log.Fatal("An error occured during the query db to get id by name ", err)
	}
	var id int
	for results.Next() {
		err = results.Scan(&id)
		if err != nil {
			log.Fatal("an error occured during the scan db to get id by name ", err)
		}
	}
	fmt.Println("got id --> ", id)
	return id
}

func AddUser(reqBody []byte) Person {
	var person Person
	db := prepareDb(dbname)
	defer db.Close()
	json.Unmarshal(reqBody, &person)
	addTokenToPerson(&person)
	id, err := insert(db, person)
	if err != nil {
		//log.Fatal("Failed to insert into db ", err)
		log.Println("Failed to insert into db ", err)
		person = Person{}
	}
	_, err = addTokenToDb(&person)
	//getLastLoginToken(person.UserName)
	if err != nil {
		log.Println("An error occured during the add token to db main func ", err)
	}
	log.Printf("Inserted row with ID of %d\n", id)
	return person
}

func insert(db *sql.DB, person Person) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO Persons VALUES (?,?,?,?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(person.Id, person.UserName, person.Password, person.Token)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func Login(reqBody []byte) Person {
	var person Person
	db := prepareDb(dbname)
	defer db.Close()
	json.Unmarshal(reqBody, &person)
	temp := getPersonToLogin(person.UserName, person.Password)
	if temp.UserName != person.UserName || temp.Password != person.Password {
		log.Print("please provide the correct credentials ")
		person = Person{}
	} else {
		addTokenToPerson(&person)
		addTokenToDb(&person)
		//getLastLoginToken(person.UserName)
	}
	return person
}

func getPersonToLogin(username, password string) Person {
	var person Person
	db := prepareDb(dbname)
	defer db.Close()
	res, err := db.Query("SELECT * FROM Persons WHERE UserName = ? and Password = ?", username, password)

	if err != nil {
		log.Fatal("an error occured during the get user to login ", err)
	}
	for res.Next() {
		err := res.Scan(&person.Id, &person.UserName, &person.Password, &person.Token)
		if err != nil {
			log.Fatal("an error occured during the scan db to get user for login ", err)
		}
	}
	return person
}

func getPersonToDelToken(username string) Person {
	var person Person
	db := prepareDb(dbname)
	defer db.Close()
	res, err := db.Query("SELECT * FROM Persons WHERE UserName = ?", username)
	if err != nil {
		log.Fatal("an error occured during the get user to delete token ", err)
	}
	for res.Next() {
		err := res.Scan(&person.Id, &person.UserName, &person.Password, &person.Token)
		if err != nil {
			log.Fatal("an error occured during the scan db to get user to delete token ", err)
		}
	}
	return person
}

func addTokenToPerson(person *Person) {
	userDbID := getIdByName(person.UserName)
	token, err := middleware.CreateToken(uint64(userDbID), person.UserName)
	if err != nil {
		log.Fatal("An error occured during the produce token ", err)
	}
	// person.Token.OwnerId = person.Id
	// person.Token.Context = token
	person.Token = token
}

func addTokenToDb(person *Person) (int64, error) {
	userDbID := getIdByName(person.UserName)
	var token = Token{userDbID, person.Token}
	if token.Context == "" {
		log.Print("Unvalid token")
		return 0, errors.New("Unvalid token error")
	}
	db := prepareDb(dbname)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO Tokens VALUES (?,?)")
	if err != nil {
		log.Fatal("An error occured during the insert token into db ", err)
		return 0, err
	}
	res, err := stmt.Exec(token.OwnerId, token.Context)
	if err != nil {
		log.Fatal("An error occured during the exec db to add token ", err)
		return 0, err
	}
	return res.RowsAffected()
}

func AddTodo(username string, reqContext []byte) Person {
	person := getPersonToDelToken(username)
	var todo Todo
	json.Unmarshal(reqContext, &todo)
	todo.OwnerId = person.Id

	_, err := addTodoSubFunc(person, todo)
	if err != nil {
		log.Fatal("an error occured during the add todo the user main func ", err)
	}
	return person
}

func addTodoSubFunc(person Person, todo Todo) (int64, error) {
	// userId := getIdByName(person.UserName)
	// var todo = Todo{userId, context}
	db := prepareDb(dbname)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO Todo VALUES(?,?)")
	if err != nil {
		log.Fatal("An error occured during the insert todo into db ", err)
		return 0, err
	}
	res, err := stmt.Exec(todo.OwnerId, todo.Context)
	if err != nil {
		log.Fatal("an error occured during the exec db to add todo ", err)
	}
	return res.RowsAffected()
}

func DeleteMe(username string) int64 {
	db := prepareDb(dbname)
	defer db.Close()
	person := getPersonToDelToken(username)
	id, err := deleteUserById(db, int64(person.Id))
	if err != nil {
		log.Print("Failed to delete into db ", err)
		//os.Exit(1)
	}
	deleteAllTokens(&person)
	deleteAllTodos(&person)
	log.Printf("deleted row with ID of %d\n", id)
	return id
}

func deleteUserById(db *sql.DB, id int64) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM Persons WHERE PersonID = ?")
	if err != nil {
		log.Print("An error occured during delete the user v1")
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		log.Print("An error occured during delete the user v2")
		return 0, err
	}
	return res.RowsAffected()
}

func deleteAllTokens(person *Person) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM Tokens WHERE OwnerID = ?")
	if err != nil {
		log.Print("an error occured during the delete tokens belong to user ", err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(person.Id)
	if err != nil {
		log.Print("an error occured during the exec db to delete all tokens belong to the user ", err)
		return 0, err
	}
	return res.RowsAffected()
}

func deleteAllTodos(person *Person) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM Todo WHERE OwnerID = ?")
	if err != nil {
		log.Print("an error occured during the delete todos belong to user ", err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(person.Id)
	if err != nil {
		log.Print("an error occured during the exec db to delete all todos belong to the user ", err)
		return 0, err
	}
	return res.RowsAffected()
}

func LogOutFromAllSession(username string) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close()
	person := getPersonToDelToken(username)
	stmt, err := db.Prepare("DELETE FROM Tokens WHERE OwnerID = ?")
	if err != nil {
		log.Print("an error occured during the delete tokens belong to user ", err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(person.Id)
	if err != nil {
		log.Print("an error occured during the exec db to delete all tokens belong to the user ", err)
		return 0, err
	}
	//person.Token = ""
	return res.RowsAffected()
}

func UpdateUser(reqBody []byte, username string) (Person, error) {
	db := prepareDb(dbname)
	defer db.Close()

	var newPerson Person
	json.Unmarshal(reqBody, &newPerson)
	newPerson.Id = getIdByName(username)
	_, err := updateUser(newPerson, int64(getIdByName(username)))
	if err != nil {
		log.Fatal("Failed to update user ", err)
		return Person{}, err
	}
	addTokenToPerson(&newPerson)
	addTokenToDb(&newPerson)
	//getLastLoginToken(new.UserName)
	return newPerson, nil
}

func updateUser(person Person, oldPersonId int64) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close()
	stmt, err := db.Prepare("UPDATE Persons SET UserName = ?,Password = ? WHERE PersonID = ?")
	if err != nil {
		log.Fatal("An error occured during the update user %w", err)
		return 0, err
	}
	defer stmt.Close()
	//userId := getIdByName(person.UserName)
	//newToken, err := middleware.CreateToken(uint64(oldPersonId), person.UserName)

	if err != nil {
		log.Fatal("an error occured during the create token to update ", err)
	}
	//addTokenToPerson(&person)

	res, err := stmt.Exec(person.UserName, person.Password, oldPersonId)
	fmt.Println("token new update", person.Token)
	if err != nil {
		log.Fatal("an error occured during the exec db to update : %w", err)
		return 0, err
	}
	//addTokenToDb(&person)
	return res.RowsAffected()
}

func GetMe(username string) Person {
	person := getPersonToDelToken(username)
	//getLastLoginToken(username)
	return person
}

func GetAllTodos() []Todo {
	var allTodos []Todo
	db := prepareDb(dbname)
	defer db.Close()
	res, err := db.Query("SELECT * FROM Todo")
	if err != nil {
		log.Print("an error occured during the get todo from db ", err)
	}
	for res.Next() {
		var temp Todo
		err := res.Scan(&temp.OwnerId, &temp.Context)
		if err != nil {
			log.Print("an error occured during the scan db to get all todos ", err)
		}
		allTodos = append(allTodos, temp)
	}
	return allTodos
}

func IfTokenIsValid(token string) []string {
	var temp []string
	db := prepareDb(dbname)
	defer db.Close()
	res, err := db.Query("select * from Tokens WHERE Context = ?", token)
	if err != nil {
		log.Print("an error occured during the get token from db ", err)
		return make([]string, 0)
	} else {
		for res.Next() {
			var tempToken Token
			err = res.Scan(&tempToken.OwnerId, &tempToken.Context)
			if err != nil {
				log.Print("an error occured during the scan db to get token ", err)
				return make([]string, 0)
			} else {
				temp = append(temp, tempToken.Context)
			}
		}
	}
	return temp
}

func GetLastLoginToken(username string) string {
	db := prepareDb(dbname)
	defer db.Close()
	person := getPersonToDelToken(username)
	res, err := db.Query("SELECT Context from Tokens WHERE OwnerID = ?", person.Id)
	if err != nil {
		log.Print("an error occured during the get last token ", err)
		return ""
	}
	var tokens []string
	for res.Next() {
		var tempToken Token
		err := res.Scan(&tempToken.Context)
		if err != nil {
			log.Print("an error occured during the scan db to get last token ", err)
			break
		}
		tokens = append(tokens, tempToken.Context)
	}
	if len(tokens) == 0 {
		return ""
	} else {
		return tokens[len(tokens)-1]
	}

}
