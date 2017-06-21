package main

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/acstech/corkboard-auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func (cb *Corkboard) authToken(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//log.Println("YOU HAVE MIDDLEWARED!!!")
		var claims corkboardauth.CustomClaims
		var parse jwt.Parser
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			os.Exit(1)
		}
		authPieces := strings.Split(authHeader, " ")
		var rawToken string
		if authPieces[0] != "Bearer:" {
			w.WriteHeader(http.StatusUnauthorized)
			os.Exit(1)
		} else if authPieces[0] == "Bearer:" {
			rawToken = authPieces[1]
		}

		//keyFunc needs to be type Keyfunc and get the public key from
		//corkboardauth somehow.

		token, error := parse.ParseWithClaims(rawToken, &claims, func(_ *jwt.Token) (interface{}, error) {

			//Need to use GetPublicPem() once it is public here. It returns pem, err
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
			next(w, r, p)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//log.Println("YOU HAVE AFTERWARED!!!")
	}
}
