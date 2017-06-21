package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//GetItems is an http handler for finding an array of items and storing them in the items array
//Currently, GetItems finds items with a N1QLQuery that searches for a "type" variable
func (corkboard *Corkboard) GetItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	items, _ := corkboard.findItems()
	if items == nil {
		log.Println("Items not found")
	}
	// array of items is marshalled to JSONobject
	JSONobject, err := json.Marshal(items)
	if err != nil {
		log.Println("could not marshal items")
	}
	//All items are written out
	_, err2 := w.Write(JSONobject)
	if err != nil {
		log.Println(err2)
	}
}

//GetItemByID uses the httprouter params to find the item by id, then Marshal & Write it in JSON
func (corkboard *Corkboard) GetItemByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	theid := p.ByName("id")
	item, _ := corkboard.findItemByID(theid)
	if item == nil {
		log.Println("Item could not be found")
	}

	JSONitem, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
	}
	_, err3 := w.Write(JSONitem)
	if err3 != nil {
		log.Println(err3)
	}
}

//NewItem endpoint decodes http request, calls createNewItem
func (corkboard *Corkboard) NewItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var item NewItemReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)
	if err != nil {
		log.Println(err)
	}
	err2 := corkboard.createNewItem(item)
	if err2 != nil {
		log.Println(err2)
	}

}

//EditItem finds an item to be updated, creates a new item with new info, then appends new info to original item
func (corkboard *Corkboard) EditItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//reqitem stores information from update request
	var reqitem Item
	decoder := json.NewDecoder(r.Body)
	err2 := decoder.Decode(&reqitem)
	if err2 != nil {
		log.Println(err2)
	}

	theid := p.ByName("id")
	item, err := corkboard.findItemByID(theid)
	if err != nil {
		log.Println(err)
	}
	if item == nil {
		log.Println("Item could not be found")
	}
	//original item has new data appended to its variables
	item.ItemName = reqitem.ItemName
	item.ItemDesc = reqitem.ItemDesc
	item.Category = reqitem.Category
	item.Price = reqitem.Price
	item.Status = reqitem.Status

	//call to updateItem inserts item to couchbase
	err3 := corkboard.updateItem(item)
	if err3 != nil {
		log.Println(err3)
	}
}

//DeleteItem calls removeItemByID to delete couchbase document containing item information
func (corkboard *Corkboard) DeleteItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	theid := p.ByName("id")
	var docID = "item:" + theid
	_, err := corkboard.Bucket.Remove(docID, 0)
	if err != nil {
		log.Println(err)
	}

}
