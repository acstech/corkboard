package corkboard

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	corkboardauth "github.com/acstech/corkboard-auth"
	"github.com/couchbase/gocb"
	"github.com/jasonmoore30/madhatter"
	"github.com/julienschmidt/httprouter"
)

const (
	envDev = "dev"
)

//Corkboard is an instance of the Corkboard server
type Corkboard struct {
	Bucket        *gocb.Bucket
	CorkboardAuth *corkboardauth.CorkboardAuth
	Environment   string
}

//CBConfig is all the necessary input values to configure a new CB Connection
type CBConfig struct {
	Connection  string
	BucketName  string
	BucketPass  string
	PrivateRSA  string
	Environment string
}

//NewCorkboard creates a Corkboard and connects to the  CBConfig passed to it
func NewCorkboard(config *CBConfig) (*Corkboard, error) {
	cluster, err := gocb.Connect(config.Connection)
	if err != nil {
		return nil, err
	}
	bucket, err := cluster.OpenBucket(config.BucketName, config.BucketPass)
	if err != nil {
		return nil, err
	}
	if err = createIndexes(bucket); err != nil {
		return nil, err
	}

	if config.PrivateRSA == "" {
		config.PrivateRSA = "id_rsa"
	}
	if _, err = os.Stat(config.PrivateRSA); os.IsNotExist(err) {

		//IF we dont have an RSA, make one
		key, err2 := rsa.GenerateKey(rand.Reader, 2048)
		if err2 != nil {
			log.Println(err2)
		}
		//if err = privKey.Validate(
		marshalKey := x509.MarshalPKCS1PrivateKey(key)
		privPem := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: marshalKey,
		}
		file, err2 := os.OpenFile(config.PrivateRSA, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0400)
		if err2 != nil {
			log.Println(err2)
		}
		defer file.Close() //nolint: errcheck
		err = pem.Encode(file, privPem)
		if err != nil {
			log.Println(err)
		}
	}

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
		Bucket:        bucket,
		CorkboardAuth: cba,
		Environment:   config.Environment,
	}, nil
}

func createIndexes(bucket *gocb.Bucket) error {
	indexQuery := gocb.NewN1qlQuery(fmt.Sprintf("SELECT `name` FROM system:indexes WHERE keyspace_id = '%s'", bucket.Name())).AdHoc(true) // nolint: gas
	createPrimaryQuery := gocb.NewN1qlQuery(fmt.Sprintf("CREATE PRIMARY INDEX idx_primary ON `%s` USING GSI", bucket.Name())).AdHoc(true) // nolint: gas
	res, err := bucket.ExecuteN1qlQuery(indexQuery, nil)
	if err != nil {
		return err
	}
	var idxs []string
	var row struct{ Name string }
	for res.Next(&row) {
		idxs = append(idxs, row.Name)
	}
	if !contains(idxs, "idx_primary") {
		_, err = bucket.ExecuteN1qlQuery(createPrimaryQuery, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

//Router returns the router containing the Corkboard endpoints
func (cb *Corkboard) Router() *httprouter.Router {
	router := httprouter.New()
	stdChain := madhatter.New(cb.defaultHeaders, cb.authToken)
	noAuthChain := madhatter.New(cb.defaultHeaders)
	noHeadersChain := madhatter.New(cb.authToken)
	environment := os.Getenv("CB_ENVIRONMENT")
	router.GET("/api/items", stdChain.Then(cb.GetItems))
	router.GET("/api/items/:id", stdChain.Then(cb.GetItemByID))
	router.POST("/api/items/new", stdChain.Then(cb.NewItem))
	router.PUT("/api/items/edit/:id", stdChain.Then(cb.EditItem))
	router.DELETE("/api/items/delete/:id", stdChain.Then(cb.DeleteItem))
	router.GET("/api/users", (stdChain.Then(cb.GetUsers)))
	router.GET("/api/category/:key", stdChain.Then(cb.GetItemsByCat))
	router.GET("/api/users/:id", stdChain.Then(cb.GetUser))
	router.PUT("/api/users/edit/:id", stdChain.Then(cb.UpdateUser))
	router.GET("/api/search/:key", stdChain.Then(cb.SearchUser))
	router.POST("/api/image/new", stdChain.Then(cb.NewImageURL))
	if environment == "dev" {
		router.POST("/api/image/post/:key", cb.MockS3)
		router.GET("/api/images/:key", noHeadersChain.Then(cb.GetImageMock))
	}

	router.DELETE("/api/users/delete/:id", stdChain.Then(cb.DeleteUser))
	router.POST("/api/users/register", noAuthChain.Then(cb.CorkboardAuth.RegisterUser()))
	router.POST("/api/users/auth", noAuthChain.Then(cb.CorkboardAuth.AuthUser()))
	return router

}

func contains(a []string, b string) bool {
	for _, item := range a {
		if item == b {
			return true
		}
	}
	return false
}
