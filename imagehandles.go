package corkboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
)

//New image is a handle to deal with New Image Requests
func (cb *Corkboard) NewImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var imageRes NewImageRes
	//checksum := r.Header.Get("Content-MD5")
	var picID NewImageReq
	json.NewDecoder(r.Body).Decode(&picID)
	// id := picID.ImageKey
	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	//h := md5.New()
	for i := 0; i < len(picID.PicID); i++ {

		req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("CB_S3_BUCKET")),
			Key:    aws.String(picID.PicID[i]),
		})
		checksum := picID.Checksum[i]
		req.HTTPRequest.Header.Set("Content-MD5", checksum)
		url, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Println(err)
			return
		}
		imageRes.URL[i] = url
	}
	imageResJSON, err := json.Marshal(imageRes)
	if err != nil {
		log.Println("Could not marshal image response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(imageResJSON)
}
