package corkboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	corkboardauth "github.com/acstech/corkboard-auth"
	"github.com/julienschmidt/httprouter"
)

/*type GetItemReq struct {
	Itemname string    `json:"itemname,omitempty"`
	Itemcat  string    `json:"itemcat,omitempty"`
	Itemdesc string    `json:"itemdesc,omitempty"`
	Price    string    `json:"itemprice,omitempty"`
	Status   string    `json:"salestatus,omitempty"`
	Date     time.Time `json:"date,omitempty"`
}*/

//GetItems is an http handler for finding an array of items and storing them in the items array
//Currently, GetItems finds items with a N1QLQuery that searches for a "type" variable
func (corkboard *Corkboard) GetItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	items, _ := corkboard.findItems()
	if items == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var itemsRes []GetItemRes
	for _, item := range items {

		itemRes := new(GetItemRes)
		var primaryID string
		var url string
		if corkboard.Environment == envDev {
			if item.PictureID != nil {
				primaryID = item.PictureID[0]
				url = fmt.Sprintf("http://localhost:%s/api/images/%s", os.Getenv("CB_PORT"), primaryID)
			}
		} else {
			if item.PictureID != nil {
				primaryID = item.PictureID[0]
				url = corkboard.getImageURL(primaryID)
			}
		}
		itemRes.Category = item.Category
		itemRes.DatePosted = item.DatePosted
		itemRes.ItemDesc = item.ItemDesc
		itemRes.ItemID = item.ItemID
		itemRes.ItemName = item.ItemName
		if url != "" {
			itemRes.PicURL = url
		}
		itemRes.PictureID = primaryID
		itemRes.Price = item.Price
		itemRes.Status = item.Status
		itemRes.UserID = item.UserID

		itemsRes = append(itemsRes, *itemRes)
	}

	// array of items is marshalled to JSONobject
	JSONobject, err := json.Marshal(itemsRes)
	if err != nil {
		log.Println(err)
	}
	//All items are written out
	_, err2 := w.Write(JSONobject)
	if err != nil {
		log.Println(err2)
	}
	//w.WriteHeader(http.StatusOK)
}

//GetItemsByCat provides searchability for items by category
func (corkboard *Corkboard) GetItemsByCat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	category := p.ByName("key")
	items, _ := corkboard.findItemsByCat(category)
	if items == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	JSONobj, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
	}

	_, err2 := w.Write(JSONobj)
	if err2 != nil {
		log.Println(err2)
	}
	w.WriteHeader(http.StatusOK)
}

//GetItemByID uses the httprouter params to find the item by id, then Marshal & Write it in JSON
func (corkboard *Corkboard) GetItemByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	theid := p.ByName("id")
	item, _ := corkboard.findItemByID(theid)
	if item == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// var newitem Item
	// //NewReqTransfer(item, &itemreq)
	// newitem.ItemName = item.ItemName
	// newitem.Category = item.Category
	// newitem.ItemDesc = item.ItemDesc
	// newitem.Price = item.Price
	// newitem.Status = item.Status
	// newitem.PictureID = item.PictureID
	// newitem.DatePosted = item.DatePosted

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

	claims, ok := r.Context().Value(ReqCtxClaims).(corkboardauth.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uid := claims.UID

	item.UserID = uid
	fmtErrs := corkboard.createNewItem(item)
	if len(fmtErrs.Errors) != 0 {
		errsRes, _ := json.Marshal(fmtErrs)
		_, err := w.Write(errsRes)
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

//EditItem finds an item to be updated, creates a new item with new info, then appends new info to original item
func (corkboard *Corkboard) EditItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//reqitem stores information from update request
	var reqitem NewItemReq
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Could eventually break this out into middleware!!!
	//would be much more organized!
	claims, ok := r.Context().Value(ReqCtxClaims).(corkboardauth.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uid := claims.UID
	if uid != item.UserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//TODO: I bet changing an items price to 0 still gives NaN
	//original item has new data appended to its variables
	item.ItemName = reqitem.Itemname
	item.ItemDesc = reqitem.Itemdesc
	item.Category = reqitem.Itemcat

	var priceSplit = strings.TrimPrefix(reqitem.Price, "$ ")
	priceSplit = strings.Replace(priceSplit, ",", "", -1)
	var price, error = strconv.ParseFloat(priceSplit, 64)
	if error != nil {
		log.Println(error)
		return
	}
	if priceSplit == "0.00" {
		price = 0.00
	}
	if item.Price > 10000000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item.Price = price
	item.Status = reqitem.Status
	if len(reqitem.PictureID) > 5 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item.PictureID = reqitem.PictureID
	item.DatePosted = reqitem.Date
	//The URL generated in "findItemByID" is not needed here
	item.PicURL = nil

	//call to updateItem appends item to couchbase
	errs := corkboard.updateItem(item)
	if len(errs.Errors) != 0 {
		errsRes, _ := json.Marshal(errs)
		_, err := w.Write(errsRes)
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

//DeleteItem calls removeItemByID to delete couchbase document containing item information
func (corkboard *Corkboard) DeleteItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	theid := p.ByName("id")
	item, err := corkboard.findItemByID(theid)
	if err != nil {
		log.Println(err)
	}
	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(ReqCtxClaims).(corkboardauth.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uid := claims.UID
	if uid != item.UserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	for i := 0; i < len(item.PictureID); i++ {
		err2 := corkboard.deleteImageID(item.PictureID[i])
		if err2 != nil {
			log.Println(err2)
		}
	}

	err2 := corkboard.removeItem(item)
	if err2 != nil {
		log.Println(err2)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
