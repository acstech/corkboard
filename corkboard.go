package corkboard

import (
  "errors"
  "github.com/couchbase/gocb"
  "github.com/julienschmidt/httprouter"
)

//Corkboard is an instance of the Corkboard server
type Corkboard struct {
  bucket    *gocb.bucket
  //SiteID? Or is this auto-generated?

}

//CBConfig is all the neccessary input values to configure a new CB Connection
type CBConfig struct {
  Connection string
  BucketName string
  BucketPass string
}


/*NewCorkBoard creates a Corkboard and connects to the  CBConfig passed
to it*/
func NewCorkboard(config *CBConfig) (*Corkboard, error){
  cluster, err := gocb.Connect(CBConfig.Connection)
if err != nil {
  return nil, error
}
  bucket, err = cluster.OpenBucket(CBConfig.BucketName, CBConfig.BucketPass)
  if err != nil {
    return nil, error
  }
  return &Corkboard{
    bucket: bucket
    //SiteID
  }, nil //no error if we get this far
}

//Router returns all the router containing the Corkboard endpoints
func (corkboard *Corkboard) Router() *httprouter.Router {
  router := httprouter.New()
  router.HandlerFunc("GET", "/users", corkboard.GetUsers())
  //router.HandlerFunc("GET", "/users/:id/profile", corkboard.GetUser())
  //router.HandlerFunc("POST", "/users/new", corkboard.NewUser())
  //router.HandlerFunc("PUT", "/users/:id/profile/edit", corkboard.EditUser())
  // router.HandlerFunc("DELETE", "/users/:id/profile/delete", corkboard.DeleteUser())
  // router.HandlerFunc("GET", "/items", corkboard.GetItems())
  // router.HandlerFunc("GET", "/items/:id/details", corkboard.GetItem())
  // router.HandlerFunc("POST", "items/new", corkboard.NewItem())
  // router.HandlerFunc("PUT", "/items/:id/details/edit", corkboard.EditItem())
 return router

}
