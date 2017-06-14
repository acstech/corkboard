package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

/* Need to eventually set os env varibales forthe CB Connection.
For now, the connection varibales a re hard coded*/
//main ...
func main() {
	cork, err := NewCorkboard(&CBConfig{
		Connection: "couchbase://localhost",
		BucketName: "default",
		BucketPass: "",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

		// // cluster, _ := gocb.Connect("couchbase://localhost") //localhost will be an ip address
		// // bucket, _ = cluster.OpenBucket("default", "")       //default will be your bucket name
		// router.GET("/sellers", GetSellersEndpoint)
		// router.GET("/seller/:id", GetSellerEndpoint)
		// router.POST("/seller", CreateSellerEndpoint)
		// router.POST("/seller/:id", UpdateSellerEndpoint)
		// router.DELETE("/seller/:id", DeleteSellerEndpoint)
		log.Println("Listening on port 8080")
		log.Fatal(http.ListenAndServe(":8080", cork.Router()))
	}
}
