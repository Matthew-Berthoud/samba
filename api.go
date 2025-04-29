package samba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/etclab/pre"
)

func fetch[T any](fullUrl string) T {
	u, err := url.Parse(fullUrl)
	if err != nil {
		panic(fmt.Sprintf("Invalid URL: %v", err))
	}
	resp, err := http.Get(u.String())
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("Fetching returned status %d, body: %s",
			resp.StatusCode, body))
	}

	var t T
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		panic(fmt.Sprintf("Failed to decode: %v", err))
	}

	return t
}

func FetchPublicParams(proxyId InstanceId) *pre.PublicParams {
	fullUrl := fmt.Sprintf("%s/publicParams", proxyId)
	m := fetch[PublicParamsSerialized](fullUrl)
	pp, err := m.DeSerialize()
	if err != nil {
		panic(fmt.Sprintf("Failed to deserialize public params: %v", err))
	}
	return pp
}

func FetchPublicKey(proxyId InstanceId, functionId FunctionId) *pre.PublicKey {
	fullUrl := fmt.Sprintf("%s/publicKey?functionId=%d", proxyId, functionId)
	m := fetch[PublicKeySerialized](fullUrl)
	pk, err := m.DeSerialize()
	if err != nil {
		panic(fmt.Sprintf("Failed to deserialize public key: %v", err))
	}
	return pk
}

func RegisterPublicKey(proxyId, instanceId InstanceId, pk *pre.PublicKey) {
	fullUrl := fmt.Sprintf("%s/registerPublicKey?instanceId=%s", proxyId, instanceId)
	pks := new(PublicKeySerialized)
	pks.Serialize(pk)
	body, err := json.Marshal(pks)
	if err != nil {
		log.Fatalf("Failed to marshal serialized public key: %v", err)
	}

	resp, err := http.Post(fullUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to post public key: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("RegisterPublicKey returned non-OK status: %d", resp.StatusCode)
	}
}

func SendMessage(m *SambaMessage, instanceId InstanceId) (response *http.Response, err error) {
	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(string(instanceId)+"/message", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func HandleMessage(w http.ResponseWriter, req *http.Request, kp *pre.KeyPair, pp *pre.PublicParams) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Failed to read request body: %v", err)
		return
	}

	var m SambaMessage
	if err := json.Unmarshal(body, &m); err != nil {
		http.Error(w, "Invalid message format", http.StatusBadRequest)
		log.Printf("Invalid message format: %v", err)
		return
	}

	plaintext, err := Decrypt(pp, kp.SK, &m)
	if err != nil {
		log.Printf("Failed to decrypt message: %v", err)
		http.Error(w, "Failed to decrypt message", http.StatusInternalServerError)
	}

	result := strings.ToUpper(string(plaintext))
	w.Write([]byte(result))
}
