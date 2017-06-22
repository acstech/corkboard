package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"strings"

	"github.com/acstech/corkboard-auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

//ReqCtxKeys is a type to hold all context keys
type ReqCtxKeys string

var (
	//ReqCtxClaims is a key for the custom claims in the context
	ReqCtxClaims ReqCtxKeys = "claims"
)

func (cb *Corkboard) authToken(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var claims corkboardauth.CustomClaims
		var parse jwt.Parser
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authPieces := strings.Split(authHeader, " ")
		var rawToken string
		if authPieces[0] != "Bearer:" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if authPieces[0] == "Bearer" {
			rawToken = authPieces[1]
		}

		token, error := parse.ParseWithClaims(rawToken, &claims, func(_ *jwt.Token) (interface{}, error) {

			pubPem, err := cb.CorkboardAuth.GetPublicPem()
			if err != nil {
				return nil, err
			}
			pubBlock, _ := pem.Decode(pubPem)
			return x509.ParsePKIXPublicKey(pubBlock.Bytes)
		})

		if error != nil {
			log.Println(error)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if token.Valid {
			r = r.WithContext(context.WithValue(r.Context(), ReqCtxClaims, claims))
			next(w, r, p)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

func (cb *Corkboard) defaultHeaders(httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
	}
}
