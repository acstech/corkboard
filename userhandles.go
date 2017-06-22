package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*UpdateUserReq is a structure used to deal with incoming http request body information
and add it to an existing user in the database*/
type UpdateUserReq struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

//GetUsers handles GET requests and responds with a slice of all users from couchbase
func (cb *Corkboard) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//TODO: Figure out the issue here
	//users is an aray of User structs
	users, err := cb.findUsers()
	if err != nil {
		log.Println(err)
		return
	}
	if users == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//I have an array of users (struct format) and I want to marshal
	//them to JSON and write the response
	//log.Println("Made it!")
	//Can I  marshal the whole array at once or do I need to do it
	//one at a time
	usersRes := make([]GetUserRes, len(users))

	//TODO: Debug. Any user who does not have firstname/lastname/phone
	//is returning values from whatever random user was returned last.
	for i, user := range users {
		var userRes GetUserRes
		userRes.Email = user.Email
		userRes.Firstname = user.Firstname
		userRes.Lastname = user.Lastname
		userRes.ID = user.ID
		userRes.Phone = user.Phone
		usersRes[i] = userRes
	}

	usersJSON, err := json.Marshal(usersRes)
	if err != nil {
		log.Println(err)
		return
	}
	_, error := w.Write(usersJSON)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//GetUser handles GET requests and responds with the user identified by the url param
func (cb *Corkboard) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id string = ps.ByName("id")
	log.Println(id)
	//THE ID IS MAKING IT TO HERE :
	//findUserByID is not working, panicking when serving
	user, err := cb.findUserByID(id)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//log.Println(user.Firstname)
	var userRes GetUserRes
	userRes.Email = user.Email
	userRes.Firstname = user.Firstname
	userRes.Lastname = user.Lastname
	userRes.ID = user.ID
	userRes.Phone = user.Phone
	userJSON, err := json.Marshal(userRes)
	if err != nil {
		log.Println(err)
		return
	}
	_, error := w.Write(userJSON)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
}

//UpdateUser handles POST requests with UpdateUserReq body data
func (cb *Corkboard) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var id string = ps.ByName("id")

	user, err1 := cb.findUserByID(id)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userKey := fmt.Sprintf("user:%s", id)

	userReq := new(UpdateUserReq)
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		log.Println(err)
		return
	}

	user.Firstname = userReq.Firstname
	user.Lastname = userReq.Lastname
	user.Phone = userReq.Phone
	user.Email = userReq.Email

	//TODO: Figure out how to keep Upsert from deleting the _type field without adding it
	//to the User struct

	_, error := cb.Bucket.Upsert(userKey, &struct {
		Type string `json:"_type"`
		FakeUser
	}{
		Type:     "User",
		FakeUser: FakeUser(*user),
	}, 0)
	if error != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//DeleteUser removes a user from couchbase bucket by id query
func (cb *Corkboard) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var id string = ps.ByName("id")
	key := "user:" + id

	theuser, _ := cb.findUserByID(id)
	if theuser == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err := cb.Bucket.Remove(key, 0)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
