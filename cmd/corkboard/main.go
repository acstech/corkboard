package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/acstech/corkboard"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

/* Need to eventually set os env varibales forthe CB Connection.
For now, the connection varibales are hard coded*/
//main ...
func main() {

	cork, err := corkboard.NewCorkboard(&corkboard.CBConfig{
		Connection:  os.Getenv("CB_CONNECTION"),
		BucketName:  os.Getenv("CB_BUCKET"),
		BucketPass:  os.Getenv("CB_BUCKET_PASS"),
		PrivateRSA:  os.Getenv("CB_PRIVATE_RSA"),
		Environment: os.Getenv("CB_ENVIRONMENT"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Type", "Date", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           3000,
	})
	port := fmt.Sprintf(":%s", os.Getenv("CB_PORT"))
	log.Println("Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, c.Handler(cork.Router())))

}
