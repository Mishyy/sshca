package main

type Response struct {
	Success     bool    `json:"success"`
	Certificate *string `json:"certificate,omitempty"`
	Error       *Error  `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
