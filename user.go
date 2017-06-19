package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
)

//User contains all necessary information about users
type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	//Profilepic
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	SiteID string `json:"siteid"`
	//Itemlist
}

type NewUserReq struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

//User2 is uses for querying and adding fields
type User2 struct {
	User User `json:"user"`
}

/*Eventually want to change this method to be findUser and determine what the
search key and return typeare using passed input. Currently, the query
is hardcoded and the return type is a slice containing all the users currently
in the bucket*/

//TODO: Change back to array of users

func (cb *Corkboard) findUsers() ([]User, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT email, firstname, id, lastname, phone FROM `%s` WHERE _type = 'User'", cb.Bucket.Name())) // nolint: gas
	res, err := cb.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return nil, err
	}
	defer res.Close() // nolint: errcheck

	var user = new(User)
	var users []User
	for res.Next(user) {
		users = append(users, *user)

		//log.Println(user)
		//log.Println("Break")
	}

	return users, nil
}

//TODO: Change id param to uuid to match corkboard-auth format

//findUserByID ...
func (cb *Corkboard) findUserByID(id string) (*User, error) {
	//Neither of the injected elements come from user input,
	//should be relatively safe

	// query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT email, firstname, id, lastname, phone FROM `%s` WHERE _type = 'User' AND id = '%s'", cb.Bucket.Name(), id)) //nolint: gas
	// res, err := cb.Bucket.ExecuteN1qlQuery(query, []interface{}{id})
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Close() //nolint: errcheck
	//TODO: Make sure there is a user found by that id or throw error
	key := "user:" + id
	user := new(User)
	_, err := cb.Bucket.Get(key, user)
	if err != nil {
		log.Println("Unable to get user.")
		return nil, err
	}
	return user, nil

}
