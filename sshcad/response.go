package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success     bool    `json:"success"`
	Certificate *string `json:"certificate,omitempty"`
	Error       *Error  `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r Response) Write(w http.ResponseWriter) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(r); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	code := 200
	if r.Error != nil {
		code = r.Error.Code
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprint(w, buffer.String())
}
