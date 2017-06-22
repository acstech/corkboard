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
	Sites []string `json:"sites"`
	//Itemlist
}

//FakeUser is a dummy struct used to add the "_type" field to users
type FakeUser User

//GetUserRes serves as intermediary data structure for getting user data
type GetUserRes struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	//Profilepic
	Phone string `json:"phone,omitempty"`
}

//MarshalJSON marshals user objects into JSON and adds the appropriate "_type"
// func (user *User) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(&struct {
// 		Type string `json:"_type"`
// 		FakeUser
// 	}{
// 		Type:     "User",
// 		FakeUser: FakeUser(*user),
// 	})
// }

func (cb *Corkboard) findUsers() ([]User, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT email, firstname, id, lastname, phone, sites FROM `%s` WHERE _type = 'User'", cb.Bucket.Name())) // nolint: gas
	res, err := cb.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return nil, err
	}

	defer res.Close() // nolint: errcheck

	var users []User
	user := new(User)
	for res.Next(user) {
		users = append(users, *user)
		user = new(User)
	}
	//spew.Dump(users)
	return users, nil
}

func (cb *Corkboard) findUserByID(id string) (*User, error) {

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
