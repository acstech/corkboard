package corkboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	corkboardauth "github.com/acstech/corkboard-auth"
	"github.com/julienschmidt/httprouter"
)

/*UpdateUserReq is a structure used to deal with incoming http request body information
and add it to an existing user in the database*/
type UpdateUserReq struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	PicID     string `json:"picid,omitempty"`
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
	//TODO: GetUsers should return the picID/GetURL of the user
	usersRes := make([]GetUserRes, len(users))

	for i, user := range users {
		var userRes GetUserRes
		userRes.Email = user.Email
		userRes.Firstname = user.Firstname
		userRes.Lastname = user.Lastname
		userRes.ID = user.ID
		userRes.Phone = user.Phone
		userRes.Zipcode = user.Zipcode
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
	userRes.Zipcode = user.Zipcode
	userRes.PicID = user.PicID
	var url string
	if cb.Environment == envDev {
		url = fmt.Sprintf("http://localhost:%s/api/images/%s", os.Getenv("CB_PORT"), user.PicID)
	} else {
		primaryID := user.PicID
		url = cb.getImageURL(primaryID)
	}
	userRes.PicURL = url
	itemIDList, err := cb.findUserItems(id)
	if err != nil {
		log.Println(err)
		return
	}

	var itemList []Item
	for _, element := range itemIDList {
		item, err2 := cb.findItemByID(element.ID)
		if err != nil {
			log.Println(err2)
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
	user, err := cb.findUserByKey(key)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var userRes GetUserRes
	userRes.Email = user.Email
	userRes.Firstname = user.Firstname
	userRes.Lastname = user.Lastname
	userRes.ID = user.ID
	userRes.Phone = user.Phone
	userRes.Zipcode = user.Zipcode
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

	theErr := cb.verify(userReq, id)
	if len(theErr.Errors) != 0 {
		errsRes, _ := json.Marshal(theErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errsRes) //nolint: errcheck
		return
	}

	user.Firstname = userReq.Firstname
	user.Lastname = userReq.Lastname
	user.Phone = userReq.Phone
	user.Email = userReq.Email
	user.Zipcode = userReq.Zipcode
	user.PicID = userReq.PicID

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

	//Check if user has any existing items, if so delete them
	items, err := cb.findUserItems(id)
	if err != nil {
		log.Println(err)
	}
	if len(items) != 0 && err == nil {
		for i := 0; i < len(items); i++ {

			theItems, _ := cb.findItemByID(items[i].ID)
			//delete images for each picture
			for j := 0; j < len(theItems.PictureID); j++ {
				cb.deleteImageID(theItems.PictureID[j]) //nolint: errcheck
			}
			//delete actual item
			var docID = "item:" + items[i].ID
			//is this bucket.remove errcheck important?
			cb.Bucket.Remove(docID, 0) //nolint: errcheck
		}
	}
	//delete user profile photo from S3 storage
	cb.deleteImageID(theuser.PicID) //nolint: errcheck

	//delete user from couchbase
	_, err2 := cb.Bucket.Remove(key, 0)
	if err2 != nil {
		log.Println(err2)
		return
	}

	w.WriteHeader(http.StatusOK)

}
