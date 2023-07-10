package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	r, err := http.Get("http://localhost:8080/user")
	if err != nil {
		fmt.Printf("[sshca] error: %s\n", err)
		return
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("[sshca] invalid response: %s\n", err)
		return
	}

	if !response.Success {
		fmt.Printf("[sshca] failed: %d\n\t%s\n", response.Error.Code, response.Error.Message)
		return
	}

	os.WriteFile("ssh_user_ca-cert.pub", []byte(*response.Certificate), 0400)
}
