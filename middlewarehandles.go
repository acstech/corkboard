package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func testMiddleware2(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Println("YOU HAVE MIDDLEWARED!!!")
		//auth := cba.token.Valid
		auth := true
		if auth {
			next(w, r, p)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
		log.Println("YOU HAVE AFTERWARED!!!")
	}
}
