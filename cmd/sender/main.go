package main

import (
	"fmt"
	"log"

	"github.com/etclab/samba"
)

func main() {
	plaintext := []byte("Hello from the Samba-protected cloud")
	options := samba.ParseOptions("sender")

	var c samba.SambaCrypto
	if options.UseRSA {
		c = new(samba.SambaRSA)
	} else {
		c = new(samba.SambaPRE)
	}

	result, err := samba.EncryptAndSend(samba.InstanceId("http://localhost:8080"), samba.FunctionId(123), plaintext, c)
	if err != nil {
		log.Fatalf("Encryption and sending failed: %v", err)
	}

	fmt.Printf("Sent message: %s\n", plaintext)
	fmt.Printf("Uppercase version: %s\n", result)
}
