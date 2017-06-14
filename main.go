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
	}
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", cork.Router()))

}
