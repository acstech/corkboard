package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
)

//Item struct contains properties for a standard item, not all properties are required
type Item struct {
	ItemID   string `json:"itemid,omitempty"`
	ItemName string `json:"itemname,omitempty"`
	ItemDesc string `json:"itemdesc,omitempty"`
	Category string `json:"itemcat" `
	//itempic
	Price      string `json:"price,omitempty"`
	DatePosted string `json:"date,omitempty"`
	SaleStatus string `json:"status,omitempty"`
	UserID     string `json:"userid,omitempty"`
}

//NewItemReq struct for creating new items
type NewItemReq struct {
	Itemname string `json:"itemname"`
	Itemcat  string `json:"itemcat"`
	Itemdesc string `json:"itemdesc"`
	Price    string `json:"itemprice"`
	//item picture coming up
}

//findItems takes a corkboard object and queries couchbase
func (corkboard *Corkboard) findItems() ([]Item, error) {

	//META(default).id = Items1
	//TODO: verify object is of type item
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT name FROM `%s` WHERE itemname = 'bike'", corkboard.Bucket.Name())) //nolint: gas
	log.Println(corkboard.Bucket.Name())
	rows, err := corkboard.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		fmt.Println("caught error: ", err)
		return nil, err
	}

	defer rows.Close() //nolint: errcheck

	var item = new(Item)
	var items []Item
	for rows.Next(item) {
		items = append(items, *item)
	}
	return items, nil

}

//findItemById queries for a specific item by id key
func (corkboard *Corkboard) findItemByID(itemID string) (*Item, error) {

	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT * FROM `%s` WHERE itemid = '%s'", corkboard.Bucket.Name(), itemID)) //nolint: gas
	rows, err := corkboard.Bucket.ExecuteN1qlQuery(query, []interface{}{itemID})
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close() //nolint: errcheck

	item := new(Item)
	for rows.Next(item) {
		return item, nil
	}
	return nil, err

}

//createNewItem . . .

func (corkboard *Corkboard) createNewItem(newitem NewItemReq) {
	var name = newitem.Itemname
	var desc = newitem.Itemdesc
	var cat = newitem.Itemcat
	var price = newitem.Price

	_, err := corkboard.Bucket.Upsert("Item:test", Item{ItemID: "9", ItemName: name, ItemDesc: desc, Category: cat, Price: price}, 0)
	if err != nil {
		log.Println("error:", err)
	}

}
