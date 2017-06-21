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
	Category string `json:"itemcat,omitempty" `
	//itempic
	Price      string `json:"itemprice,omitempty"`
	DatePosted string `json:"date,omitempty"`
	Status     string `json:"salestatus,omitempty"`
	UserID     string `json:"userid,omitempty"`
}

//NewItemReq struct for creating new items
type NewItemReq struct {
	Itemname string `json:"itemname,omitempty"`
	Itemcat  string `json:"itemcat,omitempty"`
	Itemdesc string `json:"itemdesc,omitempty"`
	Price    string `json:"itemprice,omitempty"`
	Status   string `json:"salestatus,omitempty"`
	//item picture coming up
}

//getUserKey concatenates the uuid with the "item" prefix
func getItemKey(id uuid.UUID) string {
	return fmt.Sprintf("item:%s", id.String())
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
	_, err := corkboard.Bucket.Get(itemkey, item)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return item, nil

}

//createNewItem is called by NewItem, takes a new item request and inserts it into the database
func (corkboard *Corkboard) createNewItem(newitem NewItemReq) error {
	//Add more fields later???
	var name = newitem.Itemname
	var desc = newitem.Itemdesc
	var cat = newitem.Itemcat
	var price = newitem.Price
	var status = newitem.Status

	//generate uuid for new item
	newID := uuid.NewV4()
	uID := newID.String()
	_, err := corkboard.Bucket.Insert(getItemKey(newID), Item{ItemID: uID, ItemName: name, ItemDesc: desc, Category: cat, Price: price, Status: status}, 0)
	return err
}

//updateItem upserts updated item object to couchbase document
func (corkboard *Corkboard) updateItem(item *Item) error {

	var theID = "item:" + item.ItemID
	_, err := corkboard.Bucket.Upsert(theID, item, 0)
	return err

}
