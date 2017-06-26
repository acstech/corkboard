package corkboard_test

import (
	"encoding/json"
	"fmt"
	"io"
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
	server     *httptest.Server
	reader     io.Reader
	newuserURL string
	//token        string
	header       string
	usersURL     string
	useridURL    string
	emailaddress string
	authStr      string
	edituserURL  string
	/*deleteuserURL string
	baduserURL   string*/
)

type Token struct {
	Token string `json:"token"`
}

func init() {

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

	server = httptest.NewServer(cork.Router())
	newuserURL = fmt.Sprintf("%s/api/users/register", server.URL)
	authStr = fmt.Sprintf("%s/api/users/auth", server.URL)
	usersURL = fmt.Sprintf("%s/api/users", server.URL)
	useridURL = fmt.Sprintf("%s/api/users/14951238-1a93-473d-b0bd-a986ea0bf76c", server.URL)
	edituserURL = fmt.Sprintf("%s/api/users/edit/14951238-1a93-473d-b0bd-a986ea0bf76c", server.URL)
	//deleteuserURL = fmt.Sprintf("%s/api/items/delete/", server.URL)
	//baduserURL = fmt.Sprintf("%s/api/items/15b27e85", server.URL)
}

//-----------------------------------------
//PASSING TESTS GO HERE
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

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

//TestGetUserPass tests GetUser, should always pass
func TestGetUserPass(t *testing.T) {

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
}

func TestEditUserPass(t *testing.T) {

	userJSON :=
		fmt.Sprintf(`{ "email":"%s@ROCKWELL", "password":"cat", "confirm":"cat", "siteId":"12341234-1234-1234-1234-123412341234", "firstname":"MARCO BELLINELLI"}`, emailaddress)
	reader := strings.NewReader(userJSON)

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
}

//-----------------------------------------
//FAILING TESTS GO HERE
//-----------------------------------------

//-----------------------------------------
//ADD SEVERAL MORE TESTS HERE
//-----------------------------------------
