package corkboard

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	Price      float64   `json:"itemprice,omitempty"`
	DatePosted time.Time `json:"date,omitempty"`
	Status     string    `json:"salestatus,omitempty"`
	UserID     string    `json:"userid,omitempty"`
}

//NewItemReq struct for creating new items
type NewItemReq struct {
	Type     string    `json:"type,omitempty"`
	Itemname string    `json:"itemname,omitempty"`
	Itemcat  string    `json:"itemcat,omitempty"`
	Itemdesc string    `json:"itemdesc,omitempty"`
	Price    string    `json:"itemprice,omitempty"`
	Status   string    `json:"salestatus,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	UserID   string    `json:"userid,omitempty"`
	//item picture coming up
}

//getUserKey concatenates the uuid with the "item" prefix
func getItemKey(id uuid.UUID) string {
	return fmt.Sprintf("item:%s", id.String())
}

//findItems takes a corkboard object and queries couchbase
func (corkboard *Corkboard) findItems() ([]Item, error) {

	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT itemid, itemname, itemdesc, itemprice, itemcat, date, userid FROM `%s` WHERE type = 'item'", corkboard.Bucket.Name())) //nolint: gas
	//log.Println(corkboard.Bucket.Name())
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

	var priceSplit = strings.Split(newitem.Price, " ")
	var price, error = strconv.ParseFloat(priceSplit[1], 64)
	if error != nil {
		log.Println(error)
		return error
	}
	var status = newitem.Status

	newID := uuid.NewV4()
	uID := newID.String()
	//log.Println("User ID is: ", newitem.UserID)
	_, err := corkboard.Bucket.Insert(getItemKey(newID), Item{ItemID: uID, Type: "item", ItemName: name, ItemDesc: desc, Category: cat, Price: price, Status: status, UserID: newitem.UserID, DatePosted: time.Now()}, 0)
	return err
}

//updateItem upserts updated item object to couchbase document
func (corkboard *Corkboard) updateItem(item *Item) error {

	var theID = "item:" + item.ItemID
	thetime := time.Now()
	item.DatePosted = thetime
	_, err := corkboard.Bucket.Upsert(theID, item, 0)
	return err

}
