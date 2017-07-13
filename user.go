package corkboard

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/couchbase/gocb"
)

//User contains all possible user profile information
type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Zipcode   string   `json:"zipcode"`
	PicID     string   `json:"picid"`
	Phone     string   `json:"phone"`
	Sites     []string `json:"sites"`
}

//ItemID is used to unmarshal userItems queries
type ItemID struct {
	ID string `json:"itemid"`
}

//FakeUser is a dummy struct used to add the "_type" field to users
type FakeUser User

//GetUserRes serves as intermediary data structure for getting user data
type GetUserRes struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	PicID     string `json:"picid,omitempty"`
	PicURL    string `json:"url,omitempty"`
	Items     []Item `json:"items,omitempty"`
}

func (cb *Corkboard) findUsers() ([]User, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT email, firstname, id, lastname, phone, sites, zipcode FROM `%s` WHERE _type = 'User'", cb.Bucket.Name())) // nolint: gas
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
	return users, nil
}

func (cb *Corkboard) findUserByID(id string) (*User, error) {
	key := "user:" + id
	user := new(User)
	_, err := cb.Bucket.Get(key, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (cb *Corkboard) findUserByKey(key string) (*User, error) {
	parseKey := strings.Split(key, "=")
	value := parseKey[1]
	var searchKey string

	if parseKey[0] == "email" {
		searchKey = "email"
	} else if parseKey[0] == "firstname" {
		searchKey = "firstname"
	} else if parseKey[0] == "lastname" {
		searchKey = "lastname"
	} else {
		log.Println("Request incorrectly formatted")
		return nil, nil
	}
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT email, firstname, id, lastname, phone, sites, zipcode FROM `%s` WHERE %s = '%s'", cb.Bucket.Name(), searchKey, value)) //nolint: gas
	res, err := cb.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return nil, err
	}

	defer res.Close() // nolint: errcheck
	user := new(User)
	userBytes := res.NextBytes()
	err = json.Unmarshal(userBytes, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (cb *Corkboard) findUserItems(userID string) ([]ItemID, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT itemid FROM `%s` WHERE type = 'item' AND userid = '%s'", cb.Bucket.Name(), userID)) //nolint: gas
	res, err := cb.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return nil, err
	}

	defer res.Close() //nolint:errcheck

	var items []ItemID
	itemID := new(ItemID)
	for res.Next(itemID) {
		items = append(items, *itemID)
		itemID = new(ItemID)
	}
	log.Println(len(items))
	return items, nil
}
