package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//UsersRes ...
// type GetUsersRes struct {
// 	Users []User `json:"users"`
// }
type GetUserRes struct {
	User User `json:"user"`
}

//GetUsers handles handles "/users"
func (cb *Corkboard) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//TODO: Figure out the issue here
	//users is an aray of User structs
	users, err := cb.findUsers()
	if err != nil {
		log.Println(err)
	}

	//I have an array of users (struct format) and I want to marshal
	//them to JSON and write the response
	//log.Println("Made it!")
	//Can I  marshal the whole array at once or do I need to do it
	//one at a time
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
	}
	w.Write(usersJSON)
}

//GetUser ..
func (cb *Corkboard) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id string = ps.ByName("id")
	log.Println(id)
	//THE ID IS MAKING IT TO HERE :
	//findUserByID is not working, panicking when serving
	user, err := cb.findUserByID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//log.Println(user.Firstname)

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	w.Write(userJSON)
}
