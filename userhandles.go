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

//GetUsers is an HTTP Router Handle to get a full list of users
func (corkboard *Corkboard) GetUsers() httprouter.Handle {
	//TODO: Figure out the issue
	users := corkboard.findUsers()
	if users == nil {
		log.Println("Not finding Users")
	}
	log.Println(users)
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		json.NewEncoder(w).Encode(users)
	}

}
