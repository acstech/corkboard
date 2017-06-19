package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
)

type Item struct {
	ItemID   string `json:"itemid"`
	ItemName string `json:"itemname"`
	ItemDesc string `json:"itemdesc"`
	Category string `json:"itemcat" `
	//itempic
	Price      string `json:"price"`
	DatePosted string `json:"date"`
	SaleStatus string `json:"status"`
	SellerId   string `json:"sellerid"`
}

type Item2 struct {
	Item Item `json:"itemname"`
}

type NewItemReq struct {
	itemname string `json:"itemname"`
	itemcat  string `json:itemcat`
	itemdesc string `json:itemdesc`
	price    string `json:itemprice`
	//item picture coming up
}

//findItems takes a corkboard object and queries couchbase
func (corkboard *Corkboard) findItems() ([]Item, error) {

	//META(default).id = Items1
	//TODO: verify object is of type item
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT name FROM `%s` WHERE itemname = 'bike'", corkboard.Bucket.Name()))
	log.Println(corkboard.Bucket.Name())
	rows, err := corkboard.Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var item = new(Item)
	var items []Item
	for rows.Next(item) {
		items = append(items, *item)
	}
	return items, err

}

//findItemById queries for a specific item by id key
func (corkboard *Corkboard) findItemById(itemId string) (*Item, error) {

	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT * FROM `%s` WHERE itemid = '%s'", corkboard.Bucket.Name(), itemId))
	rows, err := corkboard.Bucket.ExecuteN1qlQuery(query, []interface{}{itemId})
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	item := new(Item)
	for rows.Next(item) {
		return item, nil
	}
	return nil, err

}

//createNewItem . . .
/*
func (corkboard *Corkboard) createNewItem(newitem NewItemReq) error {
	var name string = newitem.itemname
	var desc string = newitem.itemdesc
	var cat string = newitem.itemcat
	var price string = newitem.price

	bucket := corkboard.Bucket.Name()
	bucket.Upsert("Item:" + strconv.Itoa())

}*/
