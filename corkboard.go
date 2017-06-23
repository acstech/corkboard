package main

import (
	"log"

	corkboardauth "github.com/acstech/corkboard-auth"
	"github.com/couchbase/gocb"
	"github.com/jasonmoore30/madhatter"
	"github.com/julienschmidt/httprouter"
)

//Corkboard is an instance of the Corkboard server
type Corkboard struct {
	Bucket        *gocb.Bucket
	CorkboardAuth *corkboardauth.CorkboardAuth
}

//CBConfig is all the necessary input values to configure a new CB Connection
type CBConfig struct {
	Connection string
	BucketName string
	BucketPass string
	PrivateRSA string
}

//NewCorkboard creates a Corkboard and connects to the  CBConfig passed to it
func NewCorkboard(config *CBConfig) (*Corkboard, error) {
	cluster, err := gocb.Connect(config.Connection)
	if err != nil {
		return nil, err
	}
	log.Println("Able to connect!")
	bucket, err := cluster.OpenBucket(config.BucketName, config.BucketPass)
	if err != nil {
		return nil, err
	}
	log.Println("successfully opened bucket: ", config.BucketName)

	cba, err := corkboardauth.New(&corkboardauth.Config{
		CBConnection:   config.Connection,
		CBBucket:       config.BucketName,
		CBBucketPass:   config.BucketPass,
		PrivateRSAFile: config.PrivateRSA,
	})
	if err != nil {
		return nil, err
	}
	return &Corkboard{
		Bucket: bucket, CorkboardAuth: cba}, nil

}

//Router returns the router containing the Corkboard endpoints
func (cb *Corkboard) Router() *httprouter.Router {
	router := httprouter.New()
	stdChain := madhatter.New(cb.corsPreflight, cb.defaultHeaders, cb.authToken)
	noAuthChain := madhatter.New(cb.corsPreflight, cb.defaultHeaders)

	router.GET("/api/items", stdChain.Then(cb.GetItems))
	router.GET("/api/items/:id", stdChain.Then(cb.GetItemByID))
	router.POST("/api/items/new", stdChain.Then(cb.NewItem))
	router.PUT("/api/items/edit/:id", stdChain.Then(cb.EditItem))
	router.DELETE("/api/items/delete/:id", stdChain.Then(cb.DeleteItem))
	router.GET("/api/users", (stdChain.Then(cb.GetUsers)))
	router.GET("/api/users/:id", stdChain.Then(cb.GetUser))
	router.PUT("/api/users/edit/:id", stdChain.Then(cb.UpdateUser))
	router.POST("/api/users/register", noAuthChain.Then(cb.CorkboardAuth.RegisterUser()))
	router.POST("/api/users/auth", noAuthChain.Then(cb.CorkboardAuth.AuthUser()))

	return router
}
