package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/ssh"
	"net/http"
	"time"
)

type SigningRequest struct {
	Type          uint32   `json:"type"`
	Principals    []string `json:"principals"`
	Key           string   `json:"pubkey"`
	PublicKey     ssh.PublicKey
	SourceAddress *[]string `json:"source-address,omitempty"`
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/host", func(w http.ResponseWriter, r *http.Request) {
		Get(ssh.HostCert, w, r)
	})
	router.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		Get(ssh.UserCert, w, r)
	})
	router.Post("/sign", Post)

	http.ListenAndServe("localhost:8080", router)
	//http.ListenAndServeTLS("localhost:8080", "ssl/ssl_cert.pem", "ssl/ssl_key.pem", router)
}

func Get(certType int, w http.ResponseWriter, r *http.Request) {
	content, err := PublicKey(certType)
	if err != nil {
		write(w, Response{Error: &Error{Code: 500, Message: err.Error()}})
		return
	}
	write(w, Response{Success: true, Certificate: content})
}

func Post(w http.ResponseWriter, r *http.Request) {
	var sr SigningRequest
	if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
		write(w, Response{Error: &Error{Code: 400, Message: err.Error()}})
		return
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sr.Key))
	if err != nil {
		write(w, Response{Error: &Error{Code: 400, Message: err.Error()}})
		return
	}

	sr.PublicKey = pubKey
	if valid := validateKey(sr.PublicKey); !valid {
		write(w, Response{Error: &Error{Code: 400, Message: "Invalid Public Key"}})
		return
	}

	cert, err := NewSigner(sr.Type).Sign(sr)
	if err != nil {
		write(w, Response{Error: &Error{Code: 400, Message: err.Error()}})
		return
	}
	write(w, Response{Success: true, Certificate: cert})

}

func write(w http.ResponseWriter, r Response) {
	r.Write(w)
}
