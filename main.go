package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/couchbase/gocb.v1"
)

type Seller struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
}

type N1qlSeller struct {
	Seller Seller `json:"seller"`
}

var bucket *gocb.Bucket

func GetSellerEndpoint(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var n1qlParams []interface{}
	query := gocb.NewN1qlQuery("SELECT * FROM `default` AS seller WHERE META(seller).id = $1")
	n1qlParams = append(n1qlParams, params.ByName("id"))
	rows, _ := bucket.ExecuteN1qlQuery(query, n1qlParams)
	var row N1qlSeller
	rows.One(&row)
	json.NewEncoder(w).Encode(row.Seller)
}

func GetSellersEndpoint(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var seller []Seller
	query := gocb.NewN1qlQuery("SELECT * FROM `default` AS seller")
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row N1qlSeller
	for rows.Next(&row) {
		seller = append(seller, row.Seller)
	}
	json.NewEncoder(w).Encode(seller)
}

func CreateSellerEndpoint(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var seller Seller
	var n1qlParams []interface{}
	_ = json.NewDecoder(req.Body).Decode(&seller)
	query := gocb.NewN1qlQuery("INSERT INTO `default` (KEY, VALUE) values ($1, {'firstname': $2, 'lastname': $3, 'email': $4})")
	n1qlParams = append(n1qlParams, uuid.NewV4().String())
	n1qlParams = append(n1qlParams, seller.Firstname)
	n1qlParams = append(n1qlParams, seller.Lastname)
	n1qlParams = append(n1qlParams, seller.Email)
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(seller)
}

func UpdateSellerEndpoint(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var seller Seller
	var n1qlParams []interface{}
	_ = json.NewDecoder(req.Body).Decode(&seller)
	query := gocb.NewN1qlQuery("UPDATE `default` USE KEYS $1 SET firstname = $2, lastname = $3, email = $4")
	n1qlParams = append(n1qlParams, params.ByName("id"))
	n1qlParams = append(n1qlParams, seller.Firstname)
	n1qlParams = append(n1qlParams, seller.Lastname)
	n1qlParams = append(n1qlParams, seller.Email)
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(seller)

}

func DeleteSellerEndpoint(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var n1qlParams []interface{}
	query := gocb.NewN1qlQuery("DELETE FROM `default` AS seller WHERE META(seller).id = $1")
	n1qlParams = append(n1qlParams, params.ByName("id"))
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&Seller{})
}

func main() {
	router := httprouter.New()
	cluster, _ := gocb.Connect("couchbase://localhost") //localhost will be an ip address
	bucket, _ = cluster.OpenBucket("default", "")       //default will be your bucket name
	router.GET("/sellers", GetSellersEndpoint)
	router.GET("/seller/:id", GetSellerEndpoint)
	router.POST("/seller", CreateSellerEndpoint)
	router.POST("/seller/:id", UpdateSellerEndpoint)
	router.DELETE("/seller/:id", DeleteSellerEndpoint)
	log.Fatal(http.ListenAndServe(":8091", router))
}
