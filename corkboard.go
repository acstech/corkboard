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
	//SiteID? Or is this auto-generated?

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
	//Connection opens successful
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
		Bucket: bucket, CorkboardAuth: cba}, nil //no error if we get this far

}

//Router returns all the router containing the Corkboard endpoints
func (cb *Corkboard) Router() *httprouter.Router {
	router := httprouter.New()
	stdChain := madhatter.New(cb.authToken)

	router.GET("/api/items", cb.GetItems)
	router.GET("/api/items/:id", cb.GetItemByID)
	router.POST("/api/items/new", cb.NewItem)
	router.PUT("/api/items/edit/:id", cb.EditItem)
	router.DELETE("/api/items/delete/:id", cb.DeleteItem)
	router.GET("/api/users", (stdChain.Then(cb.GetUsers)))
	//router.GET("/api/users", cb.GetUsers)
	router.GET("/api/users/:id", stdChain.Then(cb.GetUser))
	router.PUT("/api/users/edit/:id", stdChain.Then(cb.UpdateUser))
	router.DELETE("/api/users/delete/:id", cb.DeleteUser)
	//router.PUT("/api/users/edit/:id", cb.UpdateUser)

	router.HandlerFunc("POST", "/api/users/register", cb.CorkboardAuth.RegisterUser())
	router.HandlerFunc("POST", "/api/users/auth", cb.CorkboardAuth.AuthUser())

	return router
}
