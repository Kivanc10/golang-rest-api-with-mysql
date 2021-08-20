package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	dbop "kivancaydogmus.com/apps/userApp/dbOp"
	"kivancaydogmus.com/apps/userApp/middleware"
)

// func MainPage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to my app : %s", "user app")
// }

var counter int

var mutex = &sync.Mutex{}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, "Counter : %d", counter)
	mutex.Unlock()
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var persons []dbop.Person = dbop.GetAllUsers()
	json.NewEncoder(w).Encode(persons)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	person := dbop.AddUser(reqBody)
	if (dbop.Person{}) == person {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("This firstname is used already please choose different one")
	} else {
		//person.Token = r.Header.Get("Token") //
		r.Header.Set("Token", person.Token)
		json.NewEncoder(w).Encode(r.Header)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	person := dbop.Login(reqBody)
	if (dbop.Person{}) == person {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Please provide the correct credentials to login")
	} else {
		r.Header.Set("Token", person.Token)
		json.NewEncoder(w).Encode(r.Header)
		w.WriteHeader(http.StatusCreated)
		fmt.Println("author.. --> ", r.Header.Get("Authorization"))
	}
}

func deleteMe(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if myMap, ok := props.(jwt.MapClaims); ok {
		username := myMap["user_name"]
		if v, e := username.(string); e {
			if id := dbop.DeleteMe(v); id != 0 {
				json.NewEncoder(w).Encode("User is deleted succesfully")
			} else {
				w.WriteHeader(http.StatusNotFound)
				log.Print("An error occured during the delete the user")
			}
		}
	}

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	reqBody, _ := ioutil.ReadAll(r.Body)
	if myMap, ok := props.(jwt.MapClaims); ok {
		username := myMap["user_name"]
		//log.Print("username up ", username)
		if v, e := username.(string); e {
			w.WriteHeader(http.StatusConflict)
			if person, err := dbop.UpdateUser(reqBody, v); err != nil {
				fmt.Fprintf(w, "Failed to update user  !!!")
			} else {
				r.Header.Set("Token", person.Token)
				//fmt.Fprintf(w, "User is updated succesfully\n")
				json.NewEncoder(w).Encode(r.Header)
			}
		}
	}
}

func getMe(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if myMap, ok := props.(jwt.MapClaims); ok {
		username := myMap["user_name"]
		if v, e := username.(string); e {
			if v == "" {
				json.NewEncoder(w).Encode("please login again")
			}
			person := dbop.GetMe(v)
			person.Token = dbop.GetLastLoginToken(v)
			if person.UserName == "" || person.Token == "" || len(dbop.IfTokenIsValid(person.Token)) == 0 {
				r.Header.Set("Authorization", "")
				fmt.Println("del uth --> ", r.Header.Get("Authorization"))
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("please authenticate")
			} else {
				json.NewEncoder(w).Encode(person)
			}
		}
	} else {
		json.NewEncoder(w).Encode("we cannot get the user")
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	reqBody, _ := ioutil.ReadAll(r.Body)
	if myMap, ok := props.(jwt.MapClaims); ok {
		username := myMap["user_name"]
		if v, e := username.(string); e {
			person := dbop.GetMe(v)
			if person.Token == "" || person.UserName == "" {
				w.WriteHeader(http.StatusUnauthorized)
				r.Header.Set("Authorization", "")
				json.NewEncoder(w).Encode("please authenticate")
			} else {
				dbop.AddTodo(v, reqBody)
				json.NewEncoder(w).Encode("Your todo is saved succesfully")
			}
		}
	} else {
		json.NewEncoder(w).Encode("unable to create todo for the user")
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if myMap, ok := props.(jwt.MapClaims); ok {
		username := myMap["user_name"]
		if v, e := username.(string); e {
			todos := dbop.GetTodo(v)
			json.NewEncoder(w).Encode(todos)
		}
	}
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	todos := dbop.GetAllTodos()
	json.NewEncoder(w).Encode(todos)
}

func logout(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if myMap, ok := props.(jwt.MapClaims); ok {
		username := myMap["user_name"]
		if v, e := username.(string); e {
			_, err := dbop.LogOutFromAllSession(v)
			if err != nil {
				log.Print("an error occured during the delete tokens main f ", err)
			}
			person := dbop.GetMe(v)
			authR := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			person.Token = authR[1]
			r.Header.Set("Token", person.Token)
			json.NewEncoder(w).Encode(r.Header)
		}
	}
}

// func logoutFromlastSession(w http.ResponseWriter, r *http.Request) {
// 	//
// 	props := r.Context().Value("props")
// 	if myMap, ok := props.(jwt.MapClaims); ok {
// 		username := myMap["user_name"]
// 		if v, e := username.(string); e {
// 			person := dbop.LogOutFromLastSession(v)
// 			r.Header.Set("Authorization", "")
// 			json.NewEncoder(w).Encode(person)
// 		}
// 	}
// }

func HandleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	//	myRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.HandleFunc("/counter", incrementCounter)
	myRouter.HandleFunc("/users", getUsers)
	myRouter.HandleFunc("/signUp", addUser).Methods("POST")
	myRouter.HandleFunc("/signIn", login).Methods("POST")
	myRouter.Handle("/users/me", middleware.MiddleWare(http.HandlerFunc(getMe)))
	myRouter.Handle("/todo", middleware.MiddleWare(http.HandlerFunc(addTodo))).Methods("POST")
	myRouter.Handle("/todos/me", middleware.MiddleWare(http.HandlerFunc(getTodos)))
	myRouter.HandleFunc("/todos", getAllTodos)
	myRouter.Handle("/user/me", middleware.MiddleWare(http.HandlerFunc(deleteMe))).Methods("DELETE")
	myRouter.Handle("/users/update/me", middleware.MiddleWare(http.HandlerFunc(updateUser))).Methods("PUT")
	myRouter.Handle("/users/logout/me", middleware.MiddleWare(http.HandlerFunc(logout)))
	//myRouter.Handle("/users/logout", middleware.MiddleWare(http.HandlerFunc(logoutFromlastSession))).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
