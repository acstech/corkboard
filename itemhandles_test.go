package corkboard_test

/*
import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

var (
	newitemsURL string

	itemsURL      string
	itemidURL     string
	edititemURL   string
	deleteitemURL string
	baditemsURL   string
)

func init() {

	//server = httptest.NewServer(cork.Router())
	newitemsURL = fmt.Sprintf("%s/api/items/new", server.URL)
	itemsURL = fmt.Sprintf("%s/api/items", server.URL)
	itemidURL = fmt.Sprintf("%s/api/items/15b27e85-7497-4f57-a4ca-1d93443509aa", server.URL)
	edititemURL = fmt.Sprintf("%s/api/items/edit/15b27e85-7497-4f57-a4ca-1d93443509aa", server.URL)
	deleteitemURL = fmt.Sprintf("%s/api/items/delete/4f57", server.URL)
	baditemsURL = fmt.Sprintf("%s/api/items/15b27e85", server.URL)
}

//-----------------------------------------
//PASSING TESTS GO HERE
//-----------------------------------------

//TestCreateItemPass tests NewItem, should always pass
//A new item will be inserted into CB every time the test is run.
func TestCreateItemPass(t *testing.T) {

	itemJSON := `{ "itemname": "word@whip.com", "itemdesc": "finesse", "itemprice": "345"}`
	reader := strings.NewReader(itemJSON)

	req, err := http.NewRequest("POST", newitemsURL, reader)

	if err != nil {
		t.Error(err)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

//TestGetItemPass tests GetItems, should always pass
func TestGetItemPass(t *testing.T) {

	req, err := http.NewRequest("GET", itemsURL, nil)
	if err != nil {
		t.Error(err)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

//TestGetItemIDPass tests GetItemByID, should always pass
func TestGetItemIDPass(t *testing.T) {

	req, err := http.NewRequest("GET", itemidURL, nil)
	if err != nil {
		t.Error(err)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

//TestUpdateItemPass tests EditItem, should always pass
func TestUpdateItemPass(t *testing.T) {

	itemJSON := `{ "itemname": "WASHINGTON DC", "itemdesc": "finesse", "itemprice": "not even enough"}`
	reader := strings.NewReader(itemJSON)

	req, err := http.NewRequest("PUT", edititemURL, reader)
	if err != nil {
		t.Error(err)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

//-----------------------------------------
//FAILING TESTS GO HERE
//-----------------------------------------

//TestDeleteItemFail attempts to test DeleteItem with an invalid ID string,
// should always fail
func TestDeleteItemFail(t *testing.T) {

	req, err := http.NewRequest("DELETE", deleteitemURL, nil)
	if err != nil {
		t.Error(err)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

//TestGetItemFail tests GetItemByID, is passed an invalid ID, should fail
func TestGetItemFail(t *testing.T) {

	req, err := http.NewRequest("GET", baditemsURL, nil)
	if err != nil {
		t.Error(err)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}
*/
//-----------------------------------------
//ADD SEVERAL MORE TESTS HERE
//-----------------------------------------
