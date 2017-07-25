package corkboard_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	bucket        string
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
	globalprice  float64
	globalimage  string
	globalimage1 string

	newitemsURL   string
	itemsURL      string
	itemidURL     string
	edititemURL   string
	deleteitemURL string
	baditemsURL   string
	badedititems  string
	/*deleteuserURL string*/
	baduserURL string
	dev        bool

	newimageurl string
)

type Token struct {
	Token string `json:"token"`
}
type UserValues struct {
	TheUserID    string `json:"id"`
	TheUserEmail string `json:"email"`
}

type ItemValues struct {
	TheItemID  string  `json:"id"`
	ItemUserID string  `json:"userid"`
	PriceType  float64 `json:"price"`
	PicID      string  `json:"picid"`
}

func init() {

	//Set up connection for tests to run on
	cork, err := corkboard.NewCorkboard(&corkboard.CBConfig{
		Connection:  os.Getenv("CB_CONNECTION"),
		BucketName:  os.Getenv("CB_BUCKET"),
		BucketPass:  os.Getenv("CB_BUCKET_PASS"),
		PrivateRSA:  os.Getenv("CB_PRIVATE_RSA"),
		Environment: os.Getenv("CB_ENVIRONMENT"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(cork.Environment) == 0 {
		dev = false
	} else {
		dev = true
	}

	server := httptest.NewServer(cork.Router())
	bucket = cork.Bucket.Name()

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

	//Connection strings (image)
	newimageurl = fmt.Sprintf("%s/api/image/new", server.URL)
}

//doTest takes all necessary parameters to make test request, returns a bool indicating necessary data cleanup
func doTest(method string, url string, json string, token string, code int) bool {

	reader := strings.NewReader(json)

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		log.Println(err) //t.Error(err)
	}
	if token != "" {
		bearer := "Bearer " + token
		req.Header.Set("authorization", bearer)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		log.Println(err2) //t.Error(err2)
	}
	if res.StatusCode != code {
		log.Println("Success expected: %d, received: %d", code, res.StatusCode)
		return true
	}
	return false
}

//-----------------------------------------
//SETUP USER, RUN PASSING USER TESTS
//-----------------------------------------

//TestCreateUserPass tests out the RegisterUser function, should pass
//AND add a user to CB
func TestCreateUserPass(t *testing.T) {
	emailaddress = "Ma98nfbjh6734vdSa223b"

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"ca12341t", "confirm":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	doTest("POST", newuserURL, userJSON, "", 201)
}

//TestGetUserFail attempts to call GetUsers before authorization
func TestGetUsersFail(t *testing.T) {
	doTest("GET", usersURL, "", "", 401)
}

//TestAuthPass authorizes user and stores token for future test functions
func TestAuthPass(t *testing.T) {
	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	reader := strings.NewReader(userJSON)

	req, err := http.NewRequest("POST", authStr, reader)
	if err != nil {
		t.Error(err)
	}

	//one second delay to allow DB time to update
	timer := time.NewTimer(time.Second * 1)
	<-timer.C
	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	var theTok Token
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&theTok) //nolint: errcheck
	//theToken stores token from response for future use
	theToken = theTok.Token

	if res.StatusCode != 200 {
		t.Errorf("Success expected: 200, received: %d", res.StatusCode)
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

	//put token in request header
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)
	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}

	//store fields from returned users into array
	var Arr []UserValues
	body, _ := ioutil.ReadAll(res.Body)
	errre := json.Unmarshal(body, &Arr)
	if errre != nil {
		log.Println(errre)
	}

	//iterate through array and find authorized user by email
	for i := 0; i < len(Arr); i++ {
		email := Arr[i].TheUserEmail
		if email == "Ma98nfbjh6734vdSa223b@ROCKWELL.com" {
			globaluserid = Arr[i].TheUserID //assign globaluserid for future use
		}
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: 200, received: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestGetUserPass tests GetUser, should always pass
func TestGetUserPass(t *testing.T) {

	useridURL = fmt.Sprintf("%s/api/users/%s", serveURL, globaluserid)
	doTest("GET", useridURL, "", theToken, 200)
}

//TestEditUserPass adds first and lastname to user
func TestEditUserPass(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI", "phone":"(803) 431 - 6820"}`, emailaddress)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	doTest("PUT", edituserURL, userJSON, theToken, 200)
}

//TestSearchUserPass1 query by email
func TestSearchUserPass1(t *testing.T) {

	searchuserURL = fmt.Sprintf("%s/api/search/email=%s@ROCKWELL.com", serveURL, emailaddress)

	doTest("GET", searchuserURL, "", theToken, 200)
}

//TestSearchUserPass2 query by firstname
func TestSearchUserPass2(t *testing.T) {
	searchuserURL = fmt.Sprintf("%s/api/search/firstname=MARCO", serveURL)

	doTest("GET", searchuserURL, "", theToken, 200)
}

//TestSearchUserPass3 query by lastname
func TestSearchUserPass3(t *testing.T) {
	//searchuserURL2 := fmt.Sprintf("%s/api/search/fds=fd", serveURL)
	searchuserURL = fmt.Sprintf("%s/api/search/lastname=BELLINELI", serveURL)

	doTest("GET", searchuserURL, "", theToken, 200)
}

//-----------------------------------------
//FAILING USER TESTS GO HERE
//-----------------------------------------

//TestGetUsersFailAuth attempts to pass an invalid token
func TestGetUsersFailAuth(t *testing.T) {
	doTest("GET", usersURL, "", "123", 401)
}

//TestCreateUserFail passes an invalid email address
func TestCreateUserFail(t *testing.T) {
	emailaddress = "Ma98nfbjh6734vdSa223b"

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@RO", "password":"ca12341t", "confirm":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	doTest("POST", newuserURL, userJSON, "", 400)
}

//TestCreateUserFail2 has mismatched passwords
func TestCreateUserFail2(t *testing.T) {
	emailaddress = "Ma98nfbjh6734vdSa223b"

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"asdf", "confirm":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	doTest("POST", newuserURL, userJSON, "", 400)
}

//TestCreateUserFail3 fails due to oversized data fields
func TestCreateUserFail3(t *testing.T) {
	emailaddress = "Ma98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223ba98nfbjh6734vdSa223b"

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"ca1ca12341tca12341tca12341t2341t", "confirm":"ca1ca12341tca12341tca12341t2341t", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	doTest("POST", newuserURL, userJSON, "", 400)
}

//TestCreateUserFail4 fails due to missing data fields
func TestCreateUserFail4(t *testing.T) {
	emailaddress = "Ma734vdSa223b"

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "confirm":"ca1ct", "siteId":"12341234-1234-1234-1234-123412341234"}`, emailaddress)
	doTest("POST", newuserURL, userJSON, "", 400)
}

//TestSearchUserFail2 due to invalid search value
func TestSearchUserFail2(t *testing.T) {
	searchuserURL = fmt.Sprintf("%s/api/search/email=2345", serveURL)

	doTest("GET", searchuserURL, "", theToken, 500)
}

//FAILURE due to malformed JSON request
func TestEditUserFail(t *testing.T) {
	userJSON :=
		fmt.Sprintf(`{ "em:"%s@ROCKWELL", "password":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO BELLINELLI"}`, emailaddress)
	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestGetUsersFail2 fails due to malformed header
func TestGetUsersFail2(t *testing.T) {

	doTest("GET", usersURL, "", theToken, 401)
}

//TestEditUserFail3 lacks required field
func TestEditUserFail3(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI", "phone":"(803) 431 - 6820"}`)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestEditUserFail5 has a lastname that is too long
func TestEditUserFail5(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELIB
			ELLINELIBELLINELIBELLINELIBELLINELIBELLINELIBELLINELIBELLINELIBELLINELI", "phone":"(803) 431 - 6820"}`, emailaddress)

	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestEditUserFail6 has a firstname that is too long
func TestEditUserFail6(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
			OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO", "lastname":"BELLINELI", "phone":"(803) 431 - 6820"}`, emailaddress)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestEditUserFail7 has an email that is too long
func TestEditUserFail7(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI", "phone":"(803) 431 - 6820"}`, emailaddress)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestEditUserFail8 has a bad phone number input
func TestEditUserFail8(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI", "phone":"(803) 431"}`, emailaddress)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestEditUserFail9 has an email that invalid
func TestEditUserFail9(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI", "phone":"(803) 431 - 6820"}`, emailaddress)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//TestEditUserFail10 has an email that is too long
func TestEditUserFail10(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELLLLLLLLLLL.com", "password":"ca12341t", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO", "lastname":"BELLINELI", "phone":"(803) 431 - 6820", "zipcode":"34"}`, emailaddress)
	doTest("PUT", edituserURL, userJSON, theToken, 400)
}

//-----------------------------------------
//PASSING ITEM/IMAGE TESTS GO HERE
//-----------------------------------------

//TestCreateImageURLPass will create URL and we will store it for future use
func TestCreateImageURLPass(t *testing.T) {
	if dev {
		itemJSON := `{"checksum": "jf893qfnbsj", "extension": "jpg"}`
		reader := strings.NewReader(itemJSON)

		req, err := http.NewRequest("POST", newimageurl, reader)
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

		var Arr ItemValues
		body, _ := ioutil.ReadAll(res.Body)
		errre := json.Unmarshal(body, &Arr)
		if errre != nil {
			log.Println(errre)
		}

		//iterate through array and find images under user
		globalimage1 = Arr.PicID

		if res.StatusCode != 200 {
			t.Errorf("Success expected: 200, received: %d", res.StatusCode)
		}
	}
}

//TestNewImagePass uses the Url from TestCreateImageURLPass and puts our image in local storage
func TestNewImagePass(t *testing.T) {
	if dev {
		imageurl := fmt.Sprintf("%s/api/image/post/%s", serveURL, globalimage1)

		path := "./testassets/cat.jpg"

		pic, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
		}

		reader := bytes.NewReader(pic)
		req, err := http.NewRequest("PUT", imageurl, reader)

		if err != nil {
			log.Println(err)
		}

		res, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			t.Error(err2)
		}
		if res.StatusCode != 201 {
			t.Errorf("Success expected: 201, received: %d", res.StatusCode)
		}
	}
}

//TestCreateItemPass creates an item with multiple fields, should always pass
func TestCreateItemPass(t *testing.T) {

	if dev {
		itemJSON := fmt.Sprintf(`{ "name": "helmet", "description": "hard hat", "category": "sports", "price": "$ 2", "salestatus": "4sale", "picid": [ "%s" ] }`, globalimage1)
		doTest("POST", newitemsURL, itemJSON, theToken, 201)
	}
}

//TestGetItemsPass tests GetItems, should always pass
func TestGetItemsPass(t *testing.T) {
	if dev {
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

		var Arr []ItemValues
		body, _ := ioutil.ReadAll(res.Body)
		errre := json.Unmarshal(body, &Arr)
		if errre != nil {
			log.Println(errre)
		}

		//iterate through array and find items under user
		for i := 0; i < len(Arr); i++ {
			id := Arr[i].ItemUserID
			if id == globaluserid {
				globalitemid = Arr[i].TheItemID
				globalprice = Arr[i].PriceType
			}
		}
		intprice := int(globalprice)
		strprice := string(intprice)

		if len(strprice) == 0 {
			t.Error("Unexpected empty field")
		}

		if res.StatusCode != 200 {
			t.Errorf("Success expected: 200, received: %d", res.StatusCode)
		}
		res.Body.Close() //nolint: errcheck
	}
}

//TestDeleteItemPass cleans up database, deletes item.
func TestDeleteItemPass(t *testing.T) {
	if dev {
		deleteitemURL = fmt.Sprintf("%s/api/items/delete/%s", serveURL, globalitemid)
		if doTest("DELETE", deleteitemURL, "", theToken, 200) {
			log.Println("Warning: Potential untracked item object in CouchBase resulting from modified code:", globalitemid)
		}
	}
}

// //----------------------------------------------------------
// //Images Tests Go Here
// //----------------------------------------------------------

//TestCreateImageURLPass1 will create URL and we will store it for future use
func TestCreateImageURLPass1(t *testing.T) {
	if dev {
		itemJSON := `{"checksum": "h892y93g4rf", "extension": "jpg"}`
		reader := strings.NewReader(itemJSON)

		req, err := http.NewRequest("POST", newimageurl, reader)
		if err != nil {
			t.Error(err)
		}

		bearer := "Bearer " + theToken
		req.Header.Set("authorization", bearer)
		// timer := time.NewTimer(time.Second * 1)
		// <-timer.C
		res, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			t.Error(err2)
		}

		defer res.Body.Close() //nolint: errcheck

		var Arr ItemValues
		body, _ := ioutil.ReadAll(res.Body)
		errre := json.Unmarshal(body, &Arr)
		if errre != nil {
			log.Println(errre)
		}

		//iterate through array and find images under user
		globalimage = Arr.PicID

		if res.StatusCode != 200 {
			t.Errorf("Success expected: 200, received: %d", res.StatusCode)
		}
	}
}

//TestNewImagePass1 uses the Url from TestCreateImageURLPass and puts our image in local storage
func TestNewImagePass1(t *testing.T) {
	if dev {
		imageurl := fmt.Sprintf("%s/api/image/post/%s", serveURL, globalimage)

		path := "./testassets/cat.jpg"

		pic, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
		}

		reader := bytes.NewReader(pic)
		req, err := http.NewRequest("PUT", imageurl, reader)

		if err != nil {
			log.Println(err)
		}

		res, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			t.Error(err2)
		}
		if res.StatusCode != 201 {
			t.Errorf("Success expected: 201, received: %d", res.StatusCode)
			log.Println("Warning: Potential untracked image object in image database resulting from modified code:", globalimage)
		}
	}
}

//TestGetImagePass calls GetImage
func TestGetImagePass(t *testing.T) {
	if dev {
		geturl := fmt.Sprintf("%s/api/images/%s", serveURL, globalimage)
		doTest("GET", geturl, "", theToken, 200)
	}
}

//TestDeleteImagePass removes our image from the local folder
func TestDeleteImagePass(t *testing.T) {
	if dev {
		url := fmt.Sprintf("%s/api/images/delete/%s", serveURL, globalimage)
		if doTest("DELETE", url, "", theToken, 200) {
			log.Println("Warning: Potential untracked image object in image storage resulting from modified code:", globalimage)
		}
	}
}

//TestDeleteImageFail attempts to delete a image with an invalid url
func TestDeleteImageFail(t *testing.T) {
	if dev {
		url := fmt.Sprintf("%s/api/images/delete/%s", serveURL, "IDONTEXISTS")
		doTest("DELETE", url, "", theToken, 400)
	}
}

//TestCreateItemPass1 creates an item with multiple fields, should always pass
func TestCreateItemPass1(t *testing.T) {

	itemJSON := `{ "name": "helmet", "description": "hard hat", "category": "sports", "price": "$ 2", "salestatus": "4sale" }`
	doTest("POST", newitemsURL, itemJSON, theToken, 201)
}

//TestGetItemsPass2 tests GetItems, should always pass
func TestGetItemsPass2(t *testing.T) {

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

	var Arr []ItemValues
	body, _ := ioutil.ReadAll(res.Body)
	errre := json.Unmarshal(body, &Arr)
	if errre != nil {
		log.Println(errre)
	}

	//iterate through array and find items under user
	for i := 0; i < len(Arr); i++ {
		id := Arr[i].ItemUserID
		if id == globaluserid {
			globalitemid = Arr[i].TheItemID
			globalprice = Arr[i].PriceType
		}
	}
	intprice := int(globalprice)
	strprice := string(intprice)

	if len(strprice) == 0 {
		t.Error("Unexpected empty field")
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: 200, received: %d", res.StatusCode)
	}
	res.Body.Close() //nolint: errcheck
}

//TestGetUserPass2 will also check the User Items array
func TestGetUserPass2(t *testing.T) {
	useridURL = fmt.Sprintf("%s/api/users/%s", serveURL, globaluserid)
	doTest("GET", useridURL, "", theToken, 200)
}

//TestGetItemIDPass tests GetItemByID, should always pass
func TestGetItemIDPass(t *testing.T) {

	itemidURL = fmt.Sprintf("%s/api/items/%s", serveURL, globalitemid)
	doTest("GET", itemidURL, "", theToken, 200)
}

//TestGetItemsByCatPass will do this
func TestGetItemsByCatPass(t *testing.T) {
	caturl := fmt.Sprintf("%s/api/category/%s", serveURL, "sports")
	doTest("GET", caturl, "", theToken, 200)
}

//TestUpdateItemPass tests EditItem, should always pass
func TestUpdateItemPass(t *testing.T) {

	itemJSON := `{ "name": "WASHINGTON DC", "description": "finesse", "category": "states", "price": "$ 345" }`
	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 200)
}

//TestUpdateItemFail2 tries to edit with missing field
func TestUpdateItemFail2(t *testing.T) {
	itemJSON := `{ "name": "WASHINGTON DC", "description": "finesse", "price": "$ 345" }`

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 400)
}

//TestUpdateItemFail3 exceeds data standards (name)
func TestUpdateItemFail3(t *testing.T) {

	itemJSON := `{ "name": "WASHINGTON DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCD
		DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCD", "description": "finesse", "category": "states", "price": "$ 345" }`
	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 400)
}

//TestUpdateItemFail4 is missing description field
func TestUpdateItemFail4(t *testing.T) {
	itemJSON := `{ "name": "WASHINGTON DC", "description": "", "price": "$ 345" }`

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 400)
}

//TestUpdateItemFail5 has bad price field and oversized status field
func TestUpdateItemFail5(t *testing.T) {
	itemJSON := `{ "name": "WASHINGTON DC", "description": "finesse", "price": "$ 3f3f4", "salestatus":"1234123412341234123412341234" }`

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 400)
}

//TestUpdateItemFail6 is missing all fields
func TestUpdateItemFail6(t *testing.T) {
	itemJSON := `{ "name": "", "description": "", "price": "", "category":"" }`

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 400)
}

//TestUpdateItemFail7 fails due to too many pictures
func TestUpdateItemFail7(t *testing.T) {
	itemJSON := `{ "name": "WASHINGTON DC", "description": "asdf", "price": "$ 345", "picid": [ "1", "2", "3", "4", "5", "6" ]  }`

	edititemURL = fmt.Sprintf("%s/api/items/edit/%s", serveURL, globalitemid)
	doTest("PUT", edititemURL, itemJSON, theToken, 400)
}

//TestDeleteItemPass2 cleans up database, deletes item.
func TestDeleteItemPass2(t *testing.T) {
	deleteitemURL = fmt.Sprintf("%s/api/items/delete/%s", serveURL, globalitemid)
	if doTest("DELETE", deleteitemURL, "", theToken, 200) {
		log.Println("Warning: Potential untracked item object in CouchBase resulting from modified code:", globalitemid)
	}
}

//-----------------------------------------
//FAILING ITEM TESTS GO HERE
//-----------------------------------------

//TestCreateItemFail malforms the price field
func TestCreateItemFail(t *testing.T) {
	itemJSON := `{ "name": "helmet", "description": "hard hat", "category": "sports", "price": "dollars", "salestatus": "4sale" }`
	doTest("POST", newitemsURL, itemJSON, theToken, 400)
}

//TestCreateItemFail2 has empty fields
func TestCreateItemFail2(t *testing.T) {
	itemJSON := `{ "name": "", "description": "", "category": "", "price": "", "salestatus": "4sale" }`
	doTest("POST", newitemsURL, itemJSON, theToken, 400)
}

//TestCreateItemFail3 fields exceed data limits
func TestCreateItemFail3(t *testing.T) {
	itemJSON := `{ "name": "DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCD
		DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCD", "description": "", "category": "DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCD
			DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCD", "price": "$ 123049810239481029384019238401923840192384", "salestatus": "DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCD
				DCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCDCDCDCDCDDCDCD" }`
	doTest("POST", newitemsURL, itemJSON, theToken, 400)
}

//TestGetItemsByCatFail searches for nonexistent category
func TestGetItemsByCatFail(t *testing.T) {
	caturl := fmt.Sprintf("%s/api/category/%s", serveURL, "i dont live")
	doTest("GET", caturl, "", theToken, 204)
}

//TestDeleteItemFail attempts to test DeleteItem with an invalid ID string,
// should always fail
func TestDeleteItemFail(t *testing.T) {

	doTest("DELETE", deleteitemURL, "", theToken, 404)
}

//TestGetItemFail tests GetItemByID, is passed an invalid ID, should fail
func TestGetItemFail(t *testing.T) {
	doTest("GET", baditemsURL, "", theToken, 204)
}

//TestUpdateItemFail fails becaues item is gone at this point.
func TestUpdateItemFail(t *testing.T) {
	doTest("PUT", badedititems, "", theToken, 404)
}

//TestCreateItemPass2 creates an item to be deleted alongside the user
func TestCreateItemPass2(t *testing.T) {

	itemJSON := `{ "name": "helmet", "description": "hard hat", "category": "sports", "price": "$ 2", "salestatus": "4sale"}`
	doTest("POST", newitemsURL, itemJSON, theToken, 201)
}

//TestDeleteUserPass cleans up user
func TestDeleteUserPass(t *testing.T) {

	deleteuserURL = fmt.Sprintf("%s/api/users/delete/%s", serveURL, globaluserid)
	req, err := http.NewRequest("DELETE", deleteuserURL, nil)
	if err != nil {
		t.Error(err)
	}
	timer := time.NewTimer(time.Second * 1)
	<-timer.C
	bearer := "Bearer " + theToken
	req.Header.Set("authorization", bearer)

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Error(err2)
	}
	if res.StatusCode != 200 {
		t.Errorf("Success expected: 200, received: %d", res.StatusCode)
		log.Println("Warning: Potential untracked user object in CouchBase resulting from modified code:", globaluserid)
	}
	res.Body.Close() //nolint :errcheck
}

//-----------------------------------------
//FAILING USER TESTS GO HERE
//-----------------------------------------

//TestSearchUserFail fails because of invalid value
func TestSearchUserFail(t *testing.T) {
	searchuserURL = fmt.Sprintf("%s/api/search/email=%sROCKWELL", serveURL, emailaddress)
	doTest("GET", searchuserURL, "", theToken, 500)
}

//TestEditUserFail4 fails due to non-existent user
func TestEditUserFail4(t *testing.T) {
	userJSON :=
		fmt.Sprintf(`{ "email:"%s@ROCKWELL", "password":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO BELLINELLI"}`, emailaddress)

	edituserURL = fmt.Sprintf("%s/api/users/edit/%s", serveURL, globaluserid)
	doTest("PUT", edituserURL, userJSON, theToken, 404)
}

//TestGetUserFail fails due to non-existent user
func TestGetUserFail(t *testing.T) {

	doTest("GET", baduserURL, "", theToken, 404)
}

//TestDeleteUserFail fails because user does not exist
func TestDeleteUserFail(t *testing.T) {
	deleteuserURL = fmt.Sprintf("%s/api/users/delete/%s", serveURL, globaluserid)
	doTest("DELETE", deleteuserURL, "", theToken, 204)
}
