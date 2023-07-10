package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
	"time"
)

type Signer interface {
	Sign(sr SigningRequest) (cert *string, err error)
	FileName() string
}

type HostSigner struct {
	Signer
}

type UserSigner struct {
	Signer
}

func PublicKey(certType int) (pubKey *string, err error) {
	var fileName string
	switch certType {
	case ssh.HostCert:
		fileName = "ssh_host_ca.pub"
	case ssh.UserCert:
		fallthrough
	default:
		fileName = "ssh_user_ca.pub"
	}

	key, err := publicKey(fileName)
	if err != nil {
		return nil, err
	}

	key = key[:len(key)-1]
	pKey := string(key)
	return &pKey, nil
}

func NewSigner(certType uint32) Signer {
	switch certType {
	case ssh.HostCert:
		return &HostSigner{}
	case ssh.UserCert:
		fallthrough
	default:
		return &UserSigner{}
	}
}

func (hs *HostSigner) Sign(sr SigningRequest) (cert *string, err error) {
	signer, err := signer(hs)
	if err != nil {
		return nil, err
	}
	return signAndMarshal(prepareCertificate(sr), signer)
}

func (hs *HostSigner) FileName() string {
	return "ssh_host_ca"
}

func (us *UserSigner) Sign(sr SigningRequest) (cert *string, err error) {
	signer, err := signer(us)
	if err != nil {
		return nil, err
	}

	certificate := prepareCertificate(sr)
	if sr.SourceAddress != nil {
		certificate.Permissions = ssh.Permissions{CriticalOptions: map[string]string{
			"source-address": strings.Join(*sr.SourceAddress, ","),
		},
		}
	}
	return signAndMarshal(certificate, signer)
}

func (us *UserSigner) FileName() (fileName string) {
	return "ssh_user_ca"
}

func publicKey(fileName string) (publicKeyBytes []byte, err error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func signer(s Signer) (signer ssh.Signer, err error) {
	content, err := os.ReadFile(s.FileName())
	if err != nil {
		fmt.Printf("could not read certificate: %s\n", err)
		return nil, err
	}

	signer, err = ssh.ParsePrivateKey(content)
	if err != nil {
		fmt.Printf("could not parse certificate: %s\n", err)
		return nil, err
	}
	return signer, nil
}

func prepareCertificate(sr SigningRequest) (certificate ssh.Certificate) {
	now := time.Now()
	certificate = ssh.Certificate{
		Key:             sr.PublicKey,
		Serial:          uint64(now.Unix()),
		CertType:        sr.Type,
		ValidPrincipals: sr.Principals,
		ValidAfter:      uint64(now.Unix()),
		ValidBefore:     uint64(now.Add(time.Minute * 5).Unix()),
	}
	return certificate
}

func signAndMarshal(certificate ssh.Certificate, signer ssh.Signer) (cert *string, err error) {
	if err := certificate.SignCert(rand.Reader, signer); err != nil {
		fmt.Printf("could not sign key: %s\n", err)
		return nil, err
	}

	buffer := &bytes.Buffer{}
	buffer.WriteString(certificate.Type())
	buffer.WriteByte(' ')
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	encoder.Write(certificate.Marshal())
	encoder.Close()
	c := buffer.String()
	return &c, nil
}
