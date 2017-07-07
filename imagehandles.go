package corkboard

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//This blank import is to ensure proper use of autoload
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

//NewImageURL is a handle to deal with New Image Requests
func (cb *Corkboard) NewImageURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if os.Getenv("CB_ENVIRONMENT") == "dev" {
		log.Println("MockURL is being called")
		picID := new(NewImageReq)
		var imageRes NewImageRes
		err := json.NewDecoder(r.Body).Decode(&picID)
		if err != nil {
			log.Println(err)
			return
		}
		imageExtension := picID.Extension
		imageGUID := uuid.NewV4()
		key := fmt.Sprintf("%s.%s", imageGUID, imageExtension)
		imageRes.ImageKey = key
		url := fmt.Sprintf("localhost:%s/api/image/post/%s", os.Getenv("CB_PORT"), key)

		imageRes.URL = url
		imageResJSON, err := json.Marshal(imageRes)
		if err != nil {
			log.Println("Could not marshal image response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(imageResJSON)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		picID := new(NewImageReq)
		var imageRes NewImageRes
		err := json.NewDecoder(r.Body).Decode(&picID)
		if err != nil {
			log.Println(err)
			return
		}
		imageExtension := picID.Extension
		imageGUID := uuid.NewV4().String()
		key := fmt.Sprintf("%s.%s", imageGUID, imageExtension)
		sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
		if err != nil {
			log.Println(err)
			return
		}
		svc := s3.New(sess)
		req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("CB_S3_BUCKET")),
			Key:    aws.String(key),
		})
		checksum := picID.Checksum
		req.HTTPRequest.Header.Set("Content-MD5", checksum)
		url, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Println(err)
			return
		}
		imageRes.ImageKey = key
		imageRes.URL = url
		//}
		imageResJSON, err := json.Marshal(imageRes)
		if err != nil {
			log.Println("Could not marshal image response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(imageResJSON)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

//DeleteImageURL is a function
func (cb *Corkboard) DeleteImageURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	key := ps.ByName("key")
	if os.Getenv("CB_ENVIRONMENT") == "dev" {
		// delete an image: get current image
		path := "./s3images"
		filepath := fmt.Sprintf("%s/%s", path, key)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			log.Println("image not exist")
		}
		os.Remove(filepath)

	} else {
		sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
		if err != nil {
			log.Println(err)
			return
		}
		svc := s3.New(sess)
		_, err2 := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("CB_S3_BUCKET")),
			Key:    aws.String(key),
		})
		if err2 != nil {
			log.Println(err2)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)

}

//MockS3 checks for directory where files will be stored. If they don't, create it for them
// the "presigned url's" that direct to this endpoint will have to be mocked by a fake "dev env"
//endpoint. This endpoint should only be used for development purposes as well.
//Still want to use the image GUID.tag as the key
func (cb *Corkboard) MockS3(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := "./s3images"

	key := ps.ByName("key")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777) //nolint: gas, errcheck
		log.Println("Created directory 's3images' inside current directory")
	}
	defer r.Body.Close() //nolint: errcheck

	file, err := os.Create(fmt.Sprintf("%s/%s", path, key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println(file.Name())
	//It works up to here, file is created in the correct dir with the correct name.
	//Now we just need to be able to read the data from the form and copy it into the
	//file we have just created somehow.
	image, _, err1 := r.FormFile("Image")
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusInternalServerError)
	}
	//log.Println(image)
	_, err = io.Copy(file, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
