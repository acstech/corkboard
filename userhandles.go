package corkboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	corkboardauth "github.com/acstech/corkboard-auth"
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
	users, err := cb.findUsers()
	if err != nil {
		log.Println(err)
		return
	}
	if users == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	usersRes := make([]GetUserRes, len(users))

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
	//log.Println(id)
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

	var userRes GetUserRes
	userRes.Email = user.Email
	userRes.Firstname = user.Firstname
	userRes.Lastname = user.Lastname
	userRes.ID = user.ID
	userRes.Phone = user.Phone

	itemIDList, err := cb.findUserItems(id)

	var itemList []Item
	for _, element := range itemIDList {
		item, err := cb.findItemByID(element.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		itemList = append(itemList, *item)
	}
	userRes.Items = itemList

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
}

//SearchUser handles request to search users by email, firstname and lastname
func (cb *Corkboard) SearchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var key string = ps.ByName("key")
	log.Println(key)
	user, err := cb.findUserByKey(key)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	claims, ok := r.Context().Value(ReqCtxClaims).(corkboardauth.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uid := claims.UID
	if uid != user.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userKey := fmt.Sprintf("user:%s", id)

	userReq := new(UpdateUserReq)
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Firstname = userReq.Firstname
	user.Lastname = userReq.Lastname
	user.Phone = userReq.Phone
	user.Email = userReq.Email

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
	claims, ok := r.Context().Value(ReqCtxClaims).(corkboardauth.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uid := claims.UID
	if uid != theuser.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err := cb.Bucket.Remove(key, 0)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
