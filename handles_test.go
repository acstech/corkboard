package corkboard_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	server        *httptest.Server //nolint: megacheck
	reader        io.Reader
	newuserURL    string
	usersURL      string
	useridURL     string
	authStr       string
	edituserURL   string
	deleteuserURL string
	searchuserURL string

	serveURL     string
	theToken     string
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
	TheUserID    string `json:"id"`
	TheUserEmail string `json:"email"`
	TheItemID    string `json:"itemid"`
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
	baduserURL = fmt.Sprintf("%s/api/users/15b27e85", server.URL)

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
	emailaddress = "Ma98nfbjh6734vdSa223b"

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
	res.Body.Close() //nolint: errcheck
}

//TestAuthPass authorizes user and stores token for future test functions
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

	//var claims corkboardauth.CustomClaims
	//var parse jwt.Parser

	var theTok Token
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&theTok) //nolint: errcheck
	//header stores token from response for future use
	theToken = theTok.Token

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestGetUserPass tests GetUsers, should always pass
//Stores user ID from first user in array form GetUsers call for future use

func TestGetUsersPass(t *testing.T) {

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)
	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	//store fields from returned users into array
	var Arr []Values
	body, _ := ioutil.ReadAll(res.Body)
	errre := json.Unmarshal(body, &Arr)
	if errre != nil {
		log.Println(errre)
	}

	//iterate through array and find authorized user by email
	for i := 0; i < len(Arr); i++ {
		email := Arr[i].TheUserEmail
		if email == "Ma98nfbjh6734vdSa223b@ROCKWELL" {
			log.Print("equality")
			globaluserid = Arr[i].TheUserID //assign globaluserid for future use
		}
	}
	//globaluserid = Arr[0].TheUserID
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestGetUserPass tests GetUser, should always pass
func TestGetUserPass(t *testing.T) {

	useridURL = fmt.Sprintf("%s/api/users/%s", serveURL, globaluserid)
	req, err := http.NewRequest("GET", useridURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestEditUserPass adds first and lastname to user
func TestEditUserPass(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL", "password":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	req, err := http.NewRequest("PUT", edituserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//FAILURE due to malformed JSON request
func TestEditUserFail(t *testing.T) {
	userJSON :=
		fmt.Sprintf(`{ "em:"%s@ROCKWELL", "password":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO BELLINELLI"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	req, err := http.NewRequest("PUT", edituserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 400 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestSearchUserPass1 query by email
func TestSearchUserPass1(t *testing.T) {

	searchuserURL = fmt.Sprintf("%s/api/search/email=%s@ROCKWELL", serveURL, emailaddress)

	req, err := http.NewRequest("GET", searchuserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestSearchUserPass2 query by firstname
func TestSearchUserPass2(t *testing.T) {
	searchuserURL = fmt.Sprintf("%s/api/search/firstname=MARCO", serveURL)

	req, err := http.NewRequest("GET", searchuserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestSearchUserPass3 query by lastname
func TestSearchUserPass3(t *testing.T) {
	searchuserURL = fmt.Sprintf("%s/api/search/lastname=BELLINELI", serveURL)

	req, err := http.NewRequest("GET", searchuserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestDeleteUserPass cleans up user
func TestDeleteUserPass(t *testing.T) {

	deleteuserURL = fmt.Sprintf("%s/api/users/delete/%s", serveURL, globaluserid)
	req, err := http.NewRequest("DELETE", deleteuserURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint :errcheck
}

//-----------------------------------------
//FAILING USER TESTS GO HERE
//-----------------------------------------

//TestSearchUserFail fails because of invalid value
func TestSearchUserFail(t *testing.T) {
	searchuserURL = fmt.Sprintf("%s/api/search/email=%sROCKWELL", serveURL, emailaddress)

	req, err := http.NewRequest("GET", searchuserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 500 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestEditUserFail2 fails due to non-existent user
func TestEditUserFail2(t *testing.T) {
	userJSON :=
		fmt.Sprintf(`{ "email:"%s@ROCKWELL", "password":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO BELLINELLI"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	req, err := http.NewRequest("PUT", edituserURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestGetUserFail fails due to non-existent user
func TestGetUserFail(t *testing.T) {

	req, err := http.NewRequest("GET", baduserURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestDeleteUserFail fails because user does not exist
func TestDeleteUserFail(t *testing.T) {
	deleteuserURL = fmt.Sprintf("%s/api/users/delete/%s", serveURL, globaluserid)
	req, err := http.NewRequest("DELETE", deleteuserURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//-----------------------------------------
//PASSING ITEM TESTS GO HERE
//-----------------------------------------

//TestCreateItemPass creates an item with multiple fields, should always pass
func TestCreateItemPass(t *testing.T) {

	itemJSON := `{ "itemname": "helmet", "itemdesc": "hard hat", "itemcat": "sports", "itemprice": "$ 344", "salestatus": "4sale" }`
	reader := strings.NewReader(itemJSON)

	req, err := http.NewRequest("POST", newitemsURL, reader)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	defer res.Body.Close() //nolint: errcheck

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

//TestGetItemsPass tests GetItems, should always pass
func TestGetItemsPass(t *testing.T) {

	req, err := http.NewRequest("GET", itemsURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)
	timer := time.NewTimer(time.Second * 1)
	<-timer.C
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
	res.Body.Close() //nolint: errcheck
}

//TestGetItemIDPass tests GetItemByID, should always pass
func TestGetItemIDPass(t *testing.T) {

	itemidURL = fmt.Sprintf("%s/api/items/%s", serveURL, globalitemid)
	req, err := http.NewRequest("GET", itemidURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestUpdateItemPass tests EditItem, should always pass
func TestUpdateItemPass(t *testing.T) {

	itemJSON := `{ "itemname": "WASHINGTON DC", "itemdesc": "finesse", "itemprice": "$ 345"}`
	reader := strings.NewReader(itemJSON)

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	req, err := http.NewRequest("PUT", edititemURL, reader)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestDeleteItemPass cleans up database, deletes item.
func TestDeleteItemPass(t *testing.T) {
	deleteitemURL = fmt.Sprintf("%s/api/items/delete/%s", serveURL, globalitemid)
	req, err := http.NewRequest("DELETE", deleteitemURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//-----------------------------------------
//FAILING ITEM TESTS GO HERE
//-----------------------------------------

//Test on empty DB
func TestGetItemsFail(t *testing.T) {
	req, err := http.NewRequest("GET", itemsURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestDeleteItemFail attempts to test DeleteItem with an invalid ID string,
// should always fail
func TestDeleteItemFail(t *testing.T) {

	req, err := http.NewRequest("DELETE", deleteitemURL, nil)
	if err != nil {
		t.Error(err)
	}

	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestGetItemFail tests GetItemByID, is passed an invalid ID, should fail
func TestGetItemFail(t *testing.T) {
	req, err := http.NewRequest("GET", baditemsURL, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	if res.StatusCode != 204 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestUpdateItemFail fails becaues item is gone at this point.
func TestUpdateItemFail(t *testing.T) {
	req, err := http.NewRequest("PUT", badedititems, nil)
	if err != nil {
		t.Error(err)
	}
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	//204???
	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}
