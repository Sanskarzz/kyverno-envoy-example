// Test HTTP server

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// The person Type
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

var people []Person

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to my HTTP server</h1>")
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = strconv.Itoa(rand.Intn(1000000000))
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}

}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe"})
	people = append(people, Person{ID: "2", Firstname: "Jane", Lastname: "Doe"})
	people = append(people, Person{ID: "3", Firstname: "Jaden", Lastname: "Smith"})
	router.HandleFunc("/", homePage)
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeleteMovie).Methods("DELETE")
	fmt.Printf("Starting server at port 8080\n")

	log.Fatal(http.ListenAndServe(":8080", router))
}
