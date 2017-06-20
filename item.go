package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
	uuid "github.com/satori/go.uuid"
)

//Item struct contains properties for a standard item, not all properties are required
type Item struct {
	Type     string `json:"type,omitempty"`
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

	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT itemid, itemname, itemdesc, itemcat FROM `%s` WHERE type = 'item'", corkboard.Bucket.Name())) //nolint: gas
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

	item := new(Item)
	itemkey := "item:" + itemID
	fmt.Println(itemkey)
	_, err := corkboard.Bucket.Get(itemkey, item)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return item, nil

}

//getUserKey . . .
func getItemKey(id uuid.UUID) string {
	return fmt.Sprintf("item:%s", id.String())
}

//createNewItem . . .
func (corkboard *Corkboard) createNewItem(newitem NewItemReq) {
	var name = newitem.Itemname
	var desc = newitem.Itemdesc
	var cat = newitem.Itemcat
	var price = newitem.Price

	newID := uuid.NewV4()
	uID := newID.String()
	_, err := corkboard.Bucket.Insert(getItemKey(newID), Item{ItemID: uID, ItemName: name, ItemDesc: desc, Category: cat, Price: price}, 0)
	if err != nil {
		log.Println("error:", err)
	}
}
