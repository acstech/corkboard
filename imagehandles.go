package corkboard

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//This blank import is to ensure proper use of autoload
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

const path = "./s3images"
const jpeg = "jpeg"

func mockURL(w http.ResponseWriter, r *http.Request) {
	picID := new(NewImageReq)
	var imageRes NewImageRes
	err := json.NewDecoder(r.Body).Decode(&picID)
	if err != nil {
		log.Println(err)
		return
	}
	// valid image extensions
	var imageExtension string
	if picID.Extension == "jpg" {
		imageExtension = jpeg
	} else if picID.Extension == "png" || picID.Extension == jpeg {
		imageExtension = picID.Extension
	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid extension")
		return
	}

	imageGUID := uuid.NewV4()
	key := fmt.Sprintf("%s.%s", imageGUID, imageExtension)
	imageRes.ImageKey = key
	url := fmt.Sprintf("http://localhost:%s/api/image/post/%s", os.Getenv("CB_PORT"), key)

	imageRes.URL = url
	imageResJSON, err := json.Marshal(imageRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(imageResJSON)
	if err != nil {
		log.Println(err)
		return
	}
}

func imageURL(w http.ResponseWriter, r *http.Request) {
	picID := new(NewImageReq)
	var imageRes NewImageRes
	err := json.NewDecoder(r.Body).Decode(&picID)
	if err != nil {
		log.Println(err)
		return
	}
	var imageExtension string
	if picID.Extension == "jpg" {
		imageExtension = jpeg
	} else if picID.Extension == "png" || picID.Extension == jpeg {
		imageExtension = picID.Extension
	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid extension")
		return
	}
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
	//checksum := picID.Checksum
	// req.HTTPRequest.Header.Set("Content-MD5", checksum)
	// req.HTTPRequest.Header.Set("Content-Type", fmt.Sprintf("image/%s", imageExtension))
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(imageResJSON)
	if err != nil {
		log.Println(err)
		return
	}
}

//NewImageURL is a handle to deal with New Image Requests
func (cb *Corkboard) NewImageURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if cb.Environment == envDev {
		mockURL(w, r)
	} else {
		imageURL(w, r)
	}
}

//DeleteImage does a simple object removal from database
func (cb *Corkboard) DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	key := ps.ByName("key")
	err := cb.deleteImageID(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}

//MockS3 checks for directory where files will be stored. If they don't, create it for them
// the "presigned url's" that direct to this endpoint will have to be mocked by a fake "dev env"
//endpoint. This endpoint should only be used for development purposes as well.
//Still want to use the image GUID.tag as the key
func (cb *Corkboard) MockS3(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	key := ps.ByName("key")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777) //nolint: gas, errcheck
	}
	defer r.Body.Close() //nolint: errcheck

	file, err := os.Create(fmt.Sprintf("%s/%s", path, key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	_, err = io.Copy(file, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

//GetImageMock retrieves image from mocks3 storage during development
func (cb *Corkboard) GetImageMock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	ext := strings.Split(key, ".")
	var extension string
	if len(ext) == 2 {
		extension = ext[1]
	}
	image := fmt.Sprintf("%s/%s", path, key)
	// var imageByte []byte

	pic, err := ioutil.ReadFile(image)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", extension))

	_, err = w.Write(pic)
	if err != nil {
		log.Println(err)
		return
	}
}
