package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/appengine"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// User struct will be used for the json params.
type User struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Age       int    `json:"age,omitempty"`
}

var users []User

func main() {

	//Local array to have something to test.
	users = append(users, User{ID: "1", Firstname: "James", Lastname: "Hetfield", Age: 56})
	users = append(users, User{ID: "2", Firstname: "Lars", Lastname: "Ulrich", Age: 55})
	users = append(users, User{ID: "3", Firstname: "Kirk", Lastname: "Hammett", Age: 56})
	users = append(users, User{ID: "4", Firstname: "Robert", Lastname: "Trujillo", Age: 55})

	var apirouter = mux.NewRouter()
	apirouter.HandleFunc("/", health).Methods("GET")
	apirouter.HandleFunc("/users", GetUsers).Methods("GET")
	apirouter.HandleFunc("/users/{id}", GetUserID).Methods("GET")
	apirouter.HandleFunc("/users/{id}", DelUser).Methods("DELETE")
	apirouter.HandleFunc("/users/{id}", CreateUser).Methods("POST")

	fmt.Println("API up and running on port 8080")

	//Allowing all CORS calls currently.
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(apirouter)))
	appengine.Main()
}

// Health check. Polling / outputs "ok"
func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("ok")
}

// GetUsers will get all the users from the array.
func GetUsers(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// GetUserID will get users by ID.
func GetUserID(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range users {

		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// if the id is not found, still empty object
	json.NewEncoder(w).Encode(&User{})
}

// CreateUser will create a user.
func CreateUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	user.ID = params["id"]
	users = append(users, user)

	json.NewEncoder(w).Encode(users)

}

// DelUser will remove a user.
func DelUser(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}
