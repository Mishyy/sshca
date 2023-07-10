package main

import (
	"golang.org/x/crypto/ssh"
)

func validateKey(key ssh.PublicKey) bool {
	switch key.Type() {
	case ssh.KeyAlgoRSA:
		fallthrough
	case ssh.KeyAlgoED25519:
		return true
	}
	return false
}
