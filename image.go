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

// type Image struct {
// 	Image []byte `json:"image"`
// }

//NewImageReq is used to decode NewImage body into usable data.
//Used with NewImage and MockURL handles.
type NewImageReq struct {
	Checksum  string `json:"checksum"`
	Extension string `json:"extension"`
}

//ProfilePicture bundles the ImageID, userID, and image into one object to simplify
//posting/getting profile photos
// type ProfilePicture struct {
// 	ImageID string `json:"imageid"`
// 	UserID  string `json:"userid"`
// 	Image   []byte `json:"image,omitempty"`
// }

//NewImageRes bundles the image guid and the presigned url to be returned
//to the client so they can make a PUT request directly to the S3 instance. Used with
//NewImage handle
type NewImageRes struct {
	ImageKey string `json:"picid"`
	URL      string `json:"url"`
}
