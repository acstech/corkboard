package corkboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	picID := newID.String()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Println(err)
		return
	}
	imageRes.ImageID = picID
	svc := s3.New(sess)

	input := &s3.PutObjectInput{
		Bucket: aws.String("CB_S3_BUCKET"),
		Key:    aws.String(picID),
	}
	req, _ := svc.PutObjectRequest(input)

	var checksum string
	checksum = r.Header.Get("Content-MD5")
	log.Println("Checksum: ", checksum)
	//This is supposedly where the presigned URL is created, but I can't find
	//where I can use the image guid and the MD5 checksum in the url-gen
	url, err := req.Presign(time.Minute * 15)
	if err != nil {
		log.Println("Error Presigning Request")
		log.Println(err)
		return
	}
	fmt.Println("URL ", url)
	imageRes.URL = url
	imageResJSON, err := json.Marshal(imageRes)
	if err != nil {
		log.Println("Could not marshal image response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(imageResJSON)
}
