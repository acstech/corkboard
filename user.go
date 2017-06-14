package corkboard

type User struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	//Profilepic
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	//Itemlist
}
