package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
)

//User contains all possible user profile information
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	//Profilepic
	Phone string   `json:"phone"`
	Sites []string `json:"sites,omitempty"`
	//Itemlist
}

//GetUserRes serves as intermediary data structure for getting user data
type GetUserRes struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	//Profilepic
	Phone string `json:"phone"`
}

/*Eventually want to change this method to be findUser and determine what the
search key and return typeare using passed input. Currently, the query
is hardcoded and the return type is a slice containing all the users currently
in the bucket*/

//TODO: Change back to array of users

func (cb *Corkboard) findUsers() ([]GetUserRes, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT email, firstname, id, lastname, phone FROM `%s` WHERE _type = 'User'", cb.Bucket.Name())) // nolint: gas
	res, err := cb.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return nil, err
	}
	defer res.Close() // nolint: errcheck

	var user = new(User)
	var userres = new(GetUserRes)
	var users []GetUserRes
	for res.Next(user) {
		userres.ID = user.ID
		userres.Email = user.Email
		userres.Firstname = user.Firstname
		userres.Lastname = user.Lastname
		users = append(users, *userres)

		//log.Println(user)
		//log.Println("Break")
	}

	return users, nil
}

func (cb *Corkboard) findUserByID(id string) (*GetUserRes, error) {

	//TODO: Make sure there is a user found by that id or throw error
	key := "user:" + id
	user := new(User)
	_, err := cb.Bucket.Get(key, user)
	if err != nil {
		log.Println("Unable to get user.")
		return nil, err
	}
	theuser := new(GetUserRes)
	theuser.Email = user.Email
	theuser.Firstname = user.Firstname
	theuser.Lastname = user.Lastname
	theuser.Phone = user.Phone
	theuser.ID = user.ID

	return theuser, nil
}
