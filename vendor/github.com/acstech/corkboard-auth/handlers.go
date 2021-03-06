package corkboardauth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

//ErrorRes is when something in the API goes wrong
type ErrorRes struct {
	Message string `json:"message"`
}

//ErrorsRes is for when many errors can be returned
type ErrorsRes struct {
	Errors []ErrorRes `json:"errors,omitempty"`
}

//RegisterUserReq is a request from a client to create a new user
type RegisterUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
	SiteID   string `json:"siteId"`
}

//AuthReq is a request to authenticate a user
type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	SiteID   string `json:"siteId"`
}

//AuthRes is a response from the API for authenticating a user
type AuthRes struct {
	Token string `json:"token"`
}

//RegisterUser is an HTTP Router Handle for registering new users
// nolint: gocyclo
func (cba *CorkboardAuth) RegisterUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req RegisterUserReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println(err)
			writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Request must be in JSON format"}}})
			return
		}
		//TODO: Check that site exist
		var errs []ErrorRes

		email, _ := regexp.MatchString("(.+@.+\\...+)", req.Email)
		if email == false {
			errs = append(errs, ErrorRes{Message: "Email must be in valid format"})
		}

		if req.Email == "" {
			errs = append(errs, ErrorRes{Message: "Must include an email address"})
		}

		if len(req.Email) > 150 {
			errs = append(errs, ErrorRes{Message: "Email cannot be more than 150 characters"})
		}

		password, _ := regexp.MatchString("(.{8})", req.Password)
		if password == false {
			errs = append(errs, ErrorRes{Message: "Password must be at least 8 characters"})
		}

		if len(req.Password) > 20 {
			errs = append(errs, ErrorRes{Message: "Password cannot be longer than 20 characters"})
		}

		if req.SiteID == "" {
			errs = append(errs, ErrorRes{Message: "Must include a siteId"})
		} else {
			var id uuid.UUID
			id, err = uuid.FromString(req.SiteID)
			if err != nil {
				errs = append(errs, ErrorRes{Message: "siteId is not a proper ID"})
			} else {
				req.SiteID = id.String() //this is to force a certain format
			}
		}
		if req.Password == "" {
			errs = append(errs, ErrorRes{Message: "Must supply a password"})
		}
		if req.Password != req.Confirm {
			errs = append(errs, ErrorRes{Message: "password and confirm must match"})
		}
		if len(errs) > 0 {
			writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: errs})
			return
		}
		_, err = cba.findUser(req.Email)
		if err != nil {
			if err == errNoUserFound {
				var cryptPass []byte
				cryptPass, err = bcrypt.GenerateFromPassword([]byte(req.Password), 10)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
					return
				}
				err = cba.addUser(&User{
					Email:    req.Email,
					Password: base64.StdEncoding.EncodeToString(cryptPass),
					Sites:    []string{req.SiteID},
				})
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
					return
				}
				w.WriteHeader(http.StatusCreated)
				return
			}
			writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
			return
		}
		writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Email is already registered"}}})
	}
}

//AuthUser is an HTTP Router Handle for Authentication new users and return tokens
// nolint: gocyclo
func (cba *CorkboardAuth) AuthUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req AuthReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println(err)
			writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Request must be in JSON format"}}})
			return
		}
		var errs []ErrorRes
		if req.Email == "" {
			errs = append(errs, ErrorRes{Message: "Must include an email address"})
		}
		if req.Password == "" {
			errs = append(errs, ErrorRes{Message: "Must supply a password"})
		}
		if req.SiteID == "" {
			errs = append(errs, ErrorRes{Message: "Must supply a siteId"})
		} else {
			var id uuid.UUID
			id, err = uuid.FromString(req.SiteID)
			if err != nil {
				errs = append(errs, ErrorRes{Message: "siteId is not a proper ID"})
			} else {
				req.SiteID = id.String() //this is to force a certain format
			}
		}
		if len(errs) > 0 {
			writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: errs})
			return
		}
		id, _ := uuid.FromString(req.SiteID) //Already checked above
		user, err := cba.findUserFromSite(req.Email, id)
		if err != nil {
			if err == errNoUserFound {
				writeResponse(w, http.StatusUnauthorized, nil)
				return
			}
			writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
			return
		}
		cryptPass, err := base64.StdEncoding.DecodeString(user.Password)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
			return
		}
		//TODO: find a better way to check as this is taking about 2 seconds to check the passwords
		err = bcrypt.CompareHashAndPassword(cryptPass, []byte(req.Password))
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, nil)
			return
		}
		token, err := cba.generateUserToken(user)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
			return
		}
		writeResponse(w, http.StatusOK, &AuthRes{
			Token: token,
		})
	}
}

//PublicKey is a way to get the public key to verify tokens
func (cba *CorkboardAuth) PublicKey() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		pem, err := cba.GetPublicPem()
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
			return
		}
		w.Header().Set("Content-Type", "application/x-pem-file")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(pem)
		if err != nil {
			fmt.Println("Could not write to respone: ", err)
		}
	}
}

func writeResponse(w http.ResponseWriter, status int, body interface{}) {
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Could not encode the body into json"}}})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = io.Copy(w, buff)
	if err != nil {
		fmt.Println("Could not write to response: ", err)
	}
}
