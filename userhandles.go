package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//UsersRes ...
type UsersRes struct {
	Users []User `json:"users"`
}

//GetUsers handles handles "/users"
func (cb *Corkboard) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//TODO: Figure out the issue here
	//users is an aray of User structs
	users := cb.findUsers()
	if users == nil {
		log.Println("Not finding Users")
	}

	//json.NewEncoder(w).Encode(users)
	//I have an array of users (struct format) and I want to marshal
	//them to JSON and write the response
	log.Println("Made it!")
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
	}
	log.Println("UsersJSON: ", usersJSON)
	w.Write(usersJSON)

}
