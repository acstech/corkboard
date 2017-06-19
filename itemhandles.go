package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//ItemsRes does all this
type ItemsRes struct {
	Items []Item `json:"items"`
}

//GetItems is an http handler for finding an array of items and storing them in the items array
func (corkboard *Corkboard) GetItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	items, _ := corkboard.findItems()
	if items == nil {
		log.Println("Items not found")
	}

	log.Println(items)
	JSONobject, err := json.Marshal(items)
	if err != nil {
		log.Println("could not marshal items")
	}
	w.Write(JSONobject)
}

//GetItemByID is a function for finding a specific Item obect by ID
func (corkboard *Corkboard) GetItemByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	theid := p.ByName("id")
	item, _ := corkboard.findItemByID(theid)
	if item == nil {
		log.Println("Item could not be found")
	}

	JSONuser, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
	}
	w.Write(JSONuser)

}

//NewItem . . .
func (corkboard *Corkboard) NewItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var item NewItemReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)
	if err != nil {
		log.Println("issues")
	}
	corkboard.createNewItem(item)

}
