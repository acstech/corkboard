package corkboard

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

//NewImageReq is used to decode NewImage body into usable data
type NewImageReq struct {
	PicID    []string `json:"picid"`
	Checksum []string `json:"checksum"`
}

//ProfilePicture bundles the ImageID, userID, and image into one object to simplify
//posting/getting profile photos
type ProfilePicture struct {
	ImageID string `json:"imageid"`
	UserID  string `json:"userid"`
	Image   []byte `json:"image,omitempty"`
}

//NewImageRes bundles the image guid and the presigned url to be returned
//to the client so they can make a PUT request directly to the S3 instance
type NewImageRes struct {
	ImageID []string `json:"picid"`
	URL     []string `json:"url"`
}

// func (cb *Corkboard) getItemImages() ([]ItemImage, error) {
//
// }
