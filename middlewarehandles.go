package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func testMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	}
}
