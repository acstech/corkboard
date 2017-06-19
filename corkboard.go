package main

import (
	"log"

	"github.com/couchbase/gocb"
	"github.com/jasonmoore30/madhatter"
	"github.com/julienschmidt/httprouter"
)

//Corkboard is an instance of the Corkboard server
type Corkboard struct {
	Bucket *gocb.Bucket
	//SiteID? Or is this auto-generated?

}

//CBConfig is all the necessary input values to configure a new CB Connection
type CBConfig struct {
	Connection string
	BucketName string
	BucketPass string
}

//NewCorkboard creates a Corkboard and connects to the  CBConfig passed to it
func NewCorkboard(config *CBConfig) (*Corkboard, error) {
	cluster, err := gocb.Connect(config.Connection)
	if err != nil {
		return nil, err
	}
	log.Println("Able to connect!")
	//Connection opens successful
	bucket, err := cluster.OpenBucket(config.BucketName, config.BucketPass)
	if err != nil {
		return nil, err
	}
	log.Println("successfully opened bucket: ", config.BucketName)
	return &Corkboard{
		Bucket: bucket}, nil //no error if we get this far
}

//Router returns all the router containing the Corkboard endpoints
func (cb *Corkboard) Router() *httprouter.Router {
	router := httprouter.New()
	stdChain := madhatter.New(testMiddleware2)

	router.GET("/api/users", (stdChain.Then(cb.GetUsers)))
	router.GET("/api/users/:id", cb.GetUser)
	router.HandlerFunc("POST", "/api/users/register", cb.RegisterUser())
	return router
}
