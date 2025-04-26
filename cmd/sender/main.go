package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/etclab/samba"
)

func main() {
	var proxyId samba.InstanceId = "http://localhost:8080"
	var functionId samba.FunctionId = 123

	plaintext := []byte("Hello, World!")

	// request public params from proxy
	pp := samba.FetchPublicParams(proxyId)

	// request function leader's public key from proxy
	alicePk := samba.FetchPublicKey(proxyId, functionId)

	req, err := samba.PREEncrypt(pp, alicePk, plaintext, functionId)
	if err != nil {
		log.Fatalf("Proxy re-encryption failed: %v", err)
	}

	resp, err := samba.SendMessage(req, proxyId)
	if err != nil {
		log.Fatalf("Sending to proxy failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Samba Request failed with status: %v", resp.Status)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Printf("Sent message: %s\n", plaintext)
	fmt.Printf("Uppercase version: %s\n", result)
}
