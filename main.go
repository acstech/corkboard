package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/couchbase/gocb"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
}

type N1qlPerson struct {
	Person Person `json:"person"`
}

// Represents the couchbase data store
var bucket *gocb.Bucket

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var n1qlParams []interface{}
	query := gocb.NewN1qlQuery("SELECT * FROM `default` AS person WHERE META(person).id = $1")
	n1qlParams = append(n1qlParams, ps.ByName("id"))
	rows, _ := bucket.ExecuteN1qlQuery(query, n1qlParams)
	var row N1qlPerson
	rows.One(&row)
	json.NewEncoder(w).Encode(row.Person)
}

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var person []Person
	query := gocb.NewN1qlQuery("SELECT * FROM `default` AS person")
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row N1qlPerson
	for rows.Next(&row) {
		person = append(person, row.Person)
	}
	json.NewEncoder(w).Encode(person)
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var person Person
	var n1qlParams []interface{}
	_ = json.NewDecoder(req.Body).Decode(&person)
	query := gocb.NewN1qlQuery("INSERT INTO `default` (KEY, VALUE) VALUES ($1, {'firstname': $2, 'lastname': $3, 'email': $4})")
	n1qlParams = append(n1qlParams, uuid.NewV4().String())
	n1qlParams = append(n1qlParams, person.Firstname)
	n1qlParams = append(n1qlParams, person.Lastname)
	n1qlParams = append(n1qlParams, person.Email)
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(person)
}

func UpdatePersonEndpoint(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var person Person
	var n1qlParams []interface{}
	_ = json.NewDecoder(req.Body).Decode(&person)
	query := gocb.NewN1qlQuery("UPDATE `default` USE KEYS $1 SET firstname = $2, lastname = $3, email = $4")
	n1qlParams = append(n1qlParams, ps.ByName("id"))
	n1qlParams = append(n1qlParams, person.Firstname)
	n1qlParams = append(n1qlParams, person.Lastname)
	n1qlParams = append(n1qlParams, person.Email)
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(person)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var n1qlParams []interface{}
	query := gocb.NewN1qlQuery("DELETE FROM `default` AS person WHERE META(person).id = $1")
	n1qlParams = append(n1qlParams, ps.ByName("id"))
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&Person{})
}

//Index ...
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

//Hello ...
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()

	cluster, _ := gocb.Connect("couchbase://127.0.0.1")
	bucket, _ = cluster.OpenBucket("default", "")

	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	router.POST("/people/{id}", UpdatePersonEndpoint)
	router.GET("/people", GetPeopleEndpoint)
	router.GET("/people/{id}", GetPersonEndpoint)
	router.PUT("/people", CreatePersonEndpoint)
	router.DELETE("/people/{id}", DeletePersonEndpoint)

	log.Fatal(http.ListenAndServe(":8081", router))
}
