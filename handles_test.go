package corkboard_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/acstech/corkboard"
	_ "github.com/joho/godotenv/autoload"
)

var (
	//Should be able to use Server and Reader from itemhandles_test.go
	server        *httptest.Server
	reader        io.Reader
	newuserURL    string
	usersURL      string
	useridURL     string
	authStr       string
	edituserURL   string
	deleteuserURL string

	serveURL     string
	header       string
	emailaddress string
	globaluserid string
	globalitemid string

	newitemsURL   string
	itemsURL      string
	itemidURL     string
	edititemURL   string
	deleteitemURL string
	baditemsURL   string
	badedititems  string
	/*deleteuserURL string*/
	baduserURL string
)

type Token struct {
	Token string `json:"token"`
}

//new struct???
type Values struct {
	TheUserID string `json:"id"`
	TheItemID string `json:"itemid"`
}

func init() {

	//Set up connection for tests to run on
	cork, err := corkboard.NewCorkboard(&corkboard.CBConfig{
		Connection: os.Getenv("CB_CONNECTION"),
		BucketName: os.Getenv("CB_BUCKET"),
		BucketPass: os.Getenv("CB_BUCKET_PASS"),
		PrivateRSA: os.Getenv("CB_PRIVATE_RSA"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := httptest.NewServer(cork.Router())
	//Connection strings (user)
	serveURL = server.URL
	newuserURL = fmt.Sprintf("%s/api/users/register", server.URL)
	authStr = fmt.Sprintf("%s/api/users/auth", server.URL)
	usersURL = fmt.Sprintf("%s/api/users", server.URL)
	baduserURL = fmt.Sprintf("%s/api/items/15b27e85", server.URL)

	//Connection strings (item)
	newitemsURL = fmt.Sprintf("%s/api/items/new", server.URL)
	itemsURL = fmt.Sprintf("%s/api/items", server.URL)
	baditemsURL = fmt.Sprintf("%s/api/items/15b27e85", server.URL)
	badedititems = fmt.Sprintf("%s/api/items/edit/15b27e85", server.URL)
}

//-----------------------------------------
//PASSING USER TESTS GO HERE
//-----------------------------------------

//TestCreateUserPass tests out the RegisterUser function, should pass
//AND add a user to CB
func TestCreateUserPass(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var email = []rune("abcdefghijklmonpqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	for i := range b {
		b[i] = email[rand.Intn(len(email))]
	}
	emailaddress = string(b)

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL", "password":"cat", "confirm":"cat", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	req, err := http.NewRequest("POST", newuserURL, reader)

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
	res.Body.Close()
}

func TestAuthPass(t *testing.T) {
	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL", "password":"cat", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	req, err := http.NewRequest("POST", authStr, reader)
	if err != nil {
		t.Error(err)
	}

	timer := time.NewTimer(time.Second * 1)
	<-timer.C
	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	var theTok Token
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&theTok)
	header = theTok.Token

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//TestGetUserPass tests GetUsers, should always pass
func TestGetUsersPass(t *testing.T) {

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	var Arr []Values
	body, _ := ioutil.ReadAll(res.Body)
	errre := json.Unmarshal(body, &Arr)
	if errre != nil {
		log.Println(errre)
	}
	globaluserid = Arr[0].TheUserID

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//TestGetUserPass tests GetUser, should always pass
func TestGetUserPass(t *testing.T) {

	useridURL = fmt.Sprintf("%s/api/users/%s", serveURL, globaluserid)
	req, err := http.NewRequest("GET", useridURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestEditUserPass(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL", "password":"cat", "confirm":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO BELLINELLI"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	req, err := http.NewRequest("PUT", edituserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestDeleteUserPass(t *testing.T) {

	/*deleteuserURL = fmt.Sprintf("%s/api/items/delete/%s", serveURL, globaluserid)
	req, err := http.NewRequest("DELETE", deleteuserURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()*/

}

//-----------------------------------------
//FAILING USER TESTS GO HERE
//-----------------------------------------

func TestGetUserFail(t *testing.T) {

	req, err := http.NewRequest("GET", baduserURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//-----------------------------------------
//ADD SEVERAL MORE TESTS HERE
//-----------------------------------------

//-----------------------------------------
//PASSING ITEM TESTS GO HERE
//-----------------------------------------

func TestCreateItemPass(t *testing.T) {

	itemJSON := `{ "itemname": "word@whip.com", "itemdesc": "finesse", "itemprice": "345"}`
	reader := strings.NewReader(itemJSON)

	req, err := http.NewRequest("POST", newitemsURL, reader)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//TestGetItemPass tests GetItems, should always pass
func TestGetItemPass(t *testing.T) {

	req, err := http.NewRequest("GET", itemsURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	var Arr []Values
	body, _ := ioutil.ReadAll(res.Body)
	errre := json.Unmarshal(body, &Arr)
	if errre != nil {
		log.Println(errre)
	}
	globalitemid = Arr[0].TheItemID

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//TestGetItemIDPass tests GetItemByID, should always pass
func TestGetItemIDPass(t *testing.T) {

	itemidURL = fmt.Sprintf("%s/api/items/%s", serveURL, globalitemid)
	req, err := http.NewRequest("GET", itemidURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//TestUpdateItemPass tests EditItem, should always pass
func TestUpdateItemPass(t *testing.T) {

	itemJSON := `{ "itemname": "WASHINGTON DC", "itemdesc": "finesse", "itemprice": "not even enough"}`
	reader := strings.NewReader(itemJSON)

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	req, err := http.NewRequest("PUT", edititemURL, reader)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestDeleteItemPass(t *testing.T) {
	deleteitemURL = fmt.Sprintf("%s/api/items/delete/%s", serveURL, globalitemid)
	req, err := http.NewRequest("DELETE", deleteitemURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//-----------------------------------------
//FAILING ITEM TESTS GO HERE
//-----------------------------------------

//Test on empty DB
func TestGetItemsFail(t *testing.T) {
}

//TestDeleteItemFail attempts to test DeleteItem with an invalid ID string,
// should always fail
func TestDeleteItemFail(t *testing.T) {

	req, err := http.NewRequest("DELETE", deleteitemURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//TestGetItemFail tests GetItemByID, is passed an invalid ID, should fail
func TestGetItemFail(t *testing.T) {

	req, err := http.NewRequest("GET", baditemsURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestEditItemFail(t *testing.T) {
	req, err := http.NewRequest("PUT", badedititems, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + header
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	//204???
	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close()
}

//-----------------------------------------
//ADD SEVERAL MORE TESTS HERE
//-----------------------------------------