package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/etclab/samba"
)

func main() {
	expectedPlaintext := []byte("Hello, World!")

	aliceSK, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	bobSK, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	alicePK := &aliceSK.PublicKey
	bobPK := &bobSK.PublicKey

	c := new(samba.SambaRSA)

	ctAP, err := c.Encrypt(nil, alicePK, expectedPlaintext, samba.FunctionId(123))
	if err != nil {
		panic(err)
	}

	proxyPlaintext, err := c.Decrypt(nil, aliceSK, ctAP)
	if err != nil {
		panic(err)
	}

	ctPB, err := c.Encrypt(nil, bobPK, proxyPlaintext, samba.FunctionId(123))
	if err != nil {
		panic(err)
	}

	bobPlaintext, err := c.Decrypt(nil, bobSK, ctPB)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Expected plaintext: %s\n", expectedPlaintext)
	fmt.Printf("Plaintext at proxy: %s\n", proxyPlaintext)
	fmt.Printf("Plaintext at Bob: %s\n", bobPlaintext)
}
