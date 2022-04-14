package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

type JwtAuth struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJwtAuth() *JwtAuth {
	privateKey := readPEMKey()
	publicKey := &privateKey.PublicKey

	return &JwtAuth{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func GenerateKey() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)

	savePEMKey("privateKey.pem", key)
}

func readPEMKey() *rsa.PrivateKey {
	raw, err := ioutil.ReadFile("privateKey.pem")
	checkError(err)

	block, _ := pem.Decode(raw)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	checkError(err)

	return key
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
