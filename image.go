package corkboard

//
// import (
//
// )

//NewImageRes bundles the image guid and the presigned url to be returned
//to the client so they can make a PUT request directly to the S3 instance
type NewImageRes struct {
	ImageID string `json:"imageid,omitempty"`
	URL     string `json:"url,omitempty"`
}
