package corkboard

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//This blank import is to ensure proper use of autoload
	_ "github.com/joho/godotenv/autoload"
)

//
// import (
//
// )

//ItemImage bundles the ImageID, ItemID, and the image into one object to simplify
//posting/getting item photos
type ItemImage struct {
	ImageID string `json:"picid"`
	ItemID  string `json:"itemid"`
}

//NewImageReq is used to decode NewImage body into usable data.
//Used with NewImage and MockURL handles.
type NewImageReq struct {
	Checksum  string `json:"checksum"`
	Extension string `json:"extension"`
}

//NewImageRes bundles the image guid and the presigned url to be returned
//to the client so they can make a PUT request directly to the S3 instance. Used with
//NewImage handle
type NewImageRes struct {
	ImageKey string `json:"picid"`
	URL      string `json:"url"`
}

func (cb *Corkboard) getImageURL(key string) string {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Println(err)
		return ""
	}
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("CB_S3_BUCKET")),
		Key:    aws.String(key),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println(err)
		return ""
	}
	return url
}

func (cb *Corkboard) deleteImageID(key string) error {
	if cb.Environment == envDev {
		// delete an image: get current image
		filepath := fmt.Sprintf("%s/%s", path, key)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			return err
		}
		err := os.Remove(filepath)
		if err != nil {
			log.Println(err)
		}
		return nil
	} else {
		sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
		if err != nil {
			log.Println(err)
			return err
		}
		svc := s3.New(sess)
		_, err2 := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("CB_S3_BUCKET")),
			Key:    aws.String(key),
		})
		if err2 != nil {
			log.Println(err2)
			return err2
		}
		return nil
	}

}
