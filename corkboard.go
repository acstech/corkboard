package main

import (
	"github.com/couchbase/gocb"
	"github.com/julienschmidt/httprouter"
)

//Corkboard is an instance of the Corkboard server
type Corkboard struct {
	bucket *gocb.Bucket
	//SiteID? Or is this auto-generated?

}

//CBConfig is all the neccessary input values to configure a new CB Connection
type CBConfig struct {
	Connection string
	BucketName string
	BucketPass string
}

//TODO: Figure out why CBConfig.Connection and other CBConfig elements are not
//accessible and giving back "undefined" or "CBConfig does not have method Connection"

//NewCorkboard creates a Corkboard and connects to the  CBConfig passed to it
func NewCorkboard(config *CBConfig) (*Corkboard, error) {
	cluster, err := gocb.Connect("CBConfig.Connection")
	if err != nil {
		return nil, err
	}
	newbucket, err := cluster.OpenBucket("CBConfig.BucketName", "CBConfig.BucketPass")
	if err != nil {
		return nil, err
	}
	return &Corkboard{
		bucket: newbucket}, nil //no error if we get this far
}

//Router returns all the router containing the Corkboard endpoints
func (corkboard *Corkboard) Router() *httprouter.Router {
	router := httprouter.New()

	router.GET("/users", corkboard.GetUsers())

	return router

}
