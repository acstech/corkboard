package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

/* Need to eventually set os env varibales forthe CB Connection.
For now, the connection varibales are hard coded*/
//main ...
func main() {
	cork, err := NewCorkboard(&CBConfig{
		Connection: os.Getenv("CB_CONNECTION"),
		BucketName: os.Getenv("CB_BUCKET"),
		BucketPass: os.Getenv("CB_BUCKET_PASS"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", cork.Router()))

}
