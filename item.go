package corkboard

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/gocb"
	uuid "github.com/satori/go.uuid"
)

//Item struct contains properties for a standard item, not all properties are required
type Item struct {
	Type       string    `json:"type,omitempty"`
	ItemID     string    `json:"id,omitempty"`
	ItemName   string    `json:"name,omitempty"`
	ItemDesc   string    `json:"description,omitempty"`
	Category   string    `json:"category,omitempty" `
	PictureID  []string  `json:"picid,omitempty"`
	PicURL     []string  `json:"url,omitempty"`
	Price      float64   `json:"price"`
	DatePosted time.Time `json:"date,omitempty"`
	Status     string    `json:"salestatus,omitempty"`
	UserID     string    `json:"userid,omitempty"`
}

//NewItemReq struct for creating new items
type NewItemReq struct {
	Type      string    `json:"type,omitempty"`
	Itemname  string    `json:"name,omitempty"`
	Itemcat   string    `json:"category,omitempty"`
	Itemdesc  string    `json:"description,omitempty"`
	Price     string    `json:"price,omitempty"`
	PictureID []string  `json:"picid,omitempty"`
	Status    string    `json:"salestatus,omitempty"`
	Date      time.Time `json:"date,omitempty"`
	UserID    string    `json:"userid,omitempty"`
	//item picture coming up
}

//GetItemRes is used to return all image data and include the Id and url of the primary pic
type GetItemRes struct {
	ItemID     string    `json:"id,omitempty"`
	ItemName   string    `json:"name,omitempty"`
	ItemDesc   string    `json:"description,omitempty"`
	Category   string    `json:"category,omitempty"`
	PictureID  string    `json:"picid,omitempty"`
	PicURL     string    `json:"url, omitempty"`
	Price      float64   `json:"price"`
	DatePosted time.Time `json:"date,omitempty"`
	Status     string    `json:"salestatus,omitempty"`
	UserID     string    `json:"userid,omitempty"`
}

//ErrorRes contains the error message thrown by a given error
type ErrorRes struct {
	Message string `json:"message"`
}

//ErrorsRes contains an array of all the error Responses from the errors in a data access method
type ErrorsRes struct {
	Errors []ErrorRes `json:"errors,omitempty"`
}

//getUserKey concatenates the uuid with the "item" prefix
func getItemKey(id uuid.UUID) string {
	return fmt.Sprintf("item:%s", id.String())
}

//findItems takes a corkboard object and queries couchbase
func (corkboard *Corkboard) findItems() ([]Item, error) {

	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT itemid, name, description, price, category, picid, date, userid FROM `%s` WHERE type = 'item'", corkboard.Bucket.Name())) //nolint: gas
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
		item = new(Item)
	}
	return items, nil

}

//findItemsByCat queries couchbase for items by category matching the parameter
func (corkboard *Corkboard) findItemsByCat(itemCat string) ([]Item, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT id, name, description, price, category, date, userid FROM `%s` WHERE category = '%s'", corkboard.Bucket.Name(), itemCat)) //nolint: gas

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

	if corkboard.Environment == envDev {
		for _, id := range item.PictureID {
			url := fmt.Sprintf("http://localhost:%s/api/images/%s", os.Getenv("CB_PORT"), id)
			item.PicURL = append(item.PicURL, url)
		}
	} else {
		for _, id := range item.PictureID {
			url := corkboard.getImageURL(id)
			item.PicURL = append(item.PicURL, url)
		}
	}
	return item, nil
}

//createNewItem is called by NewItem, takes a new item request and inserts it into the database
func (corkboard *Corkboard) createNewItem(newitem NewItemReq) ErrorsRes {
	//Might be best to return a list of all the errors that occur at the end of the method

	errs := newitem.verify()
	// if len(errs.Errors) != 0 {
	// 	return errs
	// }
	var name = newitem.Itemname
	var desc = newitem.Itemdesc
	var cat = newitem.Itemcat
	var picid = newitem.PictureID

	var priceSplit = strings.TrimPrefix(newitem.Price, "$ ")
	priceSplit = strings.Replace(priceSplit, ",", "", -1)
	price, err := strconv.ParseFloat(priceSplit, 64)
	if err != nil {
		errs.Errors = append(errs.Errors, ErrorRes{Message: "Parsing price failed. Allowed characters: $ , . 0-9"})
	}
	if priceSplit == "0.00" || priceSplit == "" {
		price = 0.00
	}
	var status = newitem.Status
	newID := uuid.NewV4()
	uID := newID.String()
	if len(errs.Errors) != 0 {
		return errs
	}
	_, err = corkboard.Bucket.Insert(getItemKey(newID), Item{ItemID: uID, Type: "item", ItemName: name, ItemDesc: desc, Category: cat, PictureID: picid, Price: price, Status: status, UserID: newitem.UserID, DatePosted: time.Now()}, 0)
	if err != nil {
		errs.Errors = append(errs.Errors, ErrorRes{Message: err.Error()})
	}
	return errs
}

//TODO: Clean this data up as well, item checks
//updateItem upserts updated item object to couchbase document
func (corkboard *Corkboard) updateItem(item *Item) ErrorsRes {

	var theID = "item:" + item.ItemID
	thetime := time.Now()
	item.DatePosted = thetime
	errs := item.verify()
	_, err := corkboard.Bucket.Upsert(theID, item, 0)
	if err != nil {
		errs.Errors = append(errs.Errors, ErrorRes{Message: err.Error()})
	}
	return errs
}

func (corkboard *Corkboard) removeItem(item *Item) error {

	id := item.ItemID
	var theKey = "item:" + id
	_, err := corkboard.Bucket.Remove(theKey, 0)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (newitem *NewItemReq) verify() ErrorsRes {
	var errs []ErrorRes
	if len(newitem.Itemname) > 180 {
		errs = append(errs, ErrorRes{Message: "Item name greater than 180 characters."})
	} else if newitem.Itemname == "" {
		errs = append(errs, ErrorRes{Message: "Must include an item name."})
	}
	if len(newitem.Itemcat) > 25 {
		errs = append(errs, ErrorRes{Message: "Category too long."})
	} else if newitem.Itemcat == "" {
		errs = append(errs, ErrorRes{Message: "Must enter a category."})
	}

	if len(newitem.Itemdesc) > 500 {
		errs = append(errs, ErrorRes{Message: "Item description greater than 500 characters."})
	}
	if len(newitem.Price) > 12 {
		errs = append(errs, ErrorRes{Message: "Price is too large."})
	}
	if len(newitem.PictureID) > 5 {
		errs = append(errs, ErrorRes{Message: "Too many pictures uploaded. (Max 5)"})
	}
	if len(newitem.Status) > 20 {
		errs = append(errs, ErrorRes{Message: "Invalid Status."})
	}
	var fmtErrs ErrorsRes
	fmtErrs.Errors = errs
	return fmtErrs
}

func (item *Item) verify() ErrorsRes {
	var errs []ErrorRes
	if len(item.ItemName) > 180 {
		errs = append(errs, ErrorRes{Message: "Item name greater than 180 characters."})
	} else if item.ItemName == "" {
		errs = append(errs, ErrorRes{Message: "Must include an item name."})
	}
	if len(item.Category) > 25 {
		errs = append(errs, ErrorRes{Message: "Category too long."})
	} else if item.Category == "" {
		errs = append(errs, ErrorRes{Message: "Must enter a category."})
	}

	if len(item.ItemDesc) > 500 {
		errs = append(errs, ErrorRes{Message: "Item description greater than 500 characters."})
	}
	if item.Price > 10000000 {
		errs = append(errs, ErrorRes{Message: "Price is too large."})
	}
	if len(item.PictureID) > 5 {
		errs = append(errs, ErrorRes{Message: "Too many pictures uploaded. (Max 5)"})
	}
	if len(item.Status) > 20 {
		errs = append(errs, ErrorRes{Message: "Invalid Status."})
	}
	var fmtErrs ErrorsRes
	fmtErrs.Errors = errs
	return fmtErrs
}
