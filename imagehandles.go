package corkboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

//New image is a handle to deal with New Image Requests
func (cb *Corkboard) NewImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var imageRes NewImageRes
	newID := uuid.NewV4()
	imageID := newID.String()
	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	//h := md5.New()
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("CB_S3_BUCKET")),
		Key:    aws.String(os.Getenv("CB_S3_BUCKET_KEY")),
		Body:   strings.NewReader("EXPECTED CONTENTS"),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println("The URL is: ", url)
	// md5s := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// req.HTTPRequest.Header.Set("Content-MD5", md5s)

	//This is supposedly where the presigned URL is created, but I can't find
	//where I can use the image guid and the MD5 checksum in the url-gen

	// fmt.Println("URL ", url)
	imageRes.ImageID = imageID
	imageRes.URL = url
	imageResJSON, err := json.Marshal(imageRes)
	if err != nil {
		log.Println("Could not marshal image response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(imageResJSON)
}
