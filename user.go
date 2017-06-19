package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
)

//User contains all necessary information about users
type User struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	//Profilepic
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	//Itemlist
}

//User2 is uses for querying and adding fields
type User2 struct {
	User User `json:"user"`
}

/*Eventually want to change this method to be findUser and determine what the
search key and return typeare using passed input. Currently, the query
is hardcoded and the return type is a slice containing all the users currently
in the bucket*/
//findUserByEmail queries couchbase and finds users by their email address
func (corkboard *Corkboard) findUsers() []User {
	var users []User
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT * FROM `%s` WHERE _type = 'User'", corkboard.Bucket.Name()))
	log.Println(corkboard.Bucket.Name())
	rows, err := corkboard.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		fmt.Println(err)
	}
	//TODO: THink the error is occuring here
	var row User2
	for rows.Next(&row) {
		users = append(users, row.User)
	}
	return users
}

// func (corkboard *Corkboard) findUserByID(id uuid.UUID) (*User, error) {
//   query :=
//
// }
