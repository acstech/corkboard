package corkboard


import (
  "io"
	"net/http"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/acstech/corkboard"
	uuid "github.com/satori/go.uuid"
)
type UsersRes Struct {
Users []User `json:"users"`
}

//GetUsers is an HTTP Router Handle to get a full list of users
func (corkboard *Corkboard) GetUsers() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {

  }

}
