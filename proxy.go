package samba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"log"
	"net/http"

	"github.com/etclab/pre"
)

var pp *pre.PublicParams = pre.NewPublicParams()
var functionLeaders = make(map[FunctionId]InstanceId)
var keys = make(map[InstanceId]InstanceKeys)
var instances []InstanceId

func recvPublicKey(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var pks PublicKeySerialized
	err := json.NewDecoder(req.Body).Decode(&pks)
	if err != nil {
		log.Printf("Failed to decode public key: %v", err)
		http.Error(w, "Failed to decode public key", http.StatusBadRequest)
		return
	}

	pk, err := DeSerializePublicKey(pks)
	if err != nil {
		log.Printf("Failed to deserialize public key: %v", err)
		http.Error(w, "Failed to deserialize public key", http.StatusBadRequest)
		return
	}

	queries := req.URL.Query()
	instanceId := InstanceId(queries.Get("instanceId"))
	setPublicKey(instanceId, pk)
	log.Printf("Successfully storing public key for instanceId: %s", instanceId)

	w.WriteHeader(http.StatusOK)
}

func setPublicKey(instanceId InstanceId, pk pre.PublicKey) {
	keys[instanceId] = InstanceKeys{
		PublicKey:       pk,
		ReEncryptionKey: keys[instanceId].ReEncryptionKey, // Preserve existing re-encryption key if resetting
	}
}

func sendPublicParams(w http.ResponseWriter, req *http.Request) {
	pps, err := SerializePublicParams(*pp)
	if err != nil {
		http.Error(w, "Failed to serialize fields in public parameters", http.StatusInternalServerError)
		log.Printf("Failed to serialize fields in public parameters")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pps)
	if err != nil {
		http.Error(w, "Failed to encode and respond with public parameters", http.StatusInternalServerError)
		log.Printf("Failed to encode and respond with public parameters")
		return
	}
}

func getReEncryptionKey(a, b InstanceId) (pre.ReEncryptionKey, error) {
	if keys[b].ReEncryptionKey != (pre.ReEncryptionKey{}) {
		return keys[b].ReEncryptionKey, nil
	}

	pks := SerializePublicKey(keys[b].PublicKey)

	req := ReEncryptionKeyRequest{
		InstanceId:         b,
		PublicKeySerialzed: pks,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return pre.ReEncryptionKey{}, err
	}

	resp, err := http.Post(string(a)+"/requestReEncryptionKey", "application/json", bytes.NewReader(body))
	if err != nil {
		return pre.ReEncryptionKey{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pre.ReEncryptionKey{}, fmt.Errorf("requestReEncryptionKey failed with status %d", resp.StatusCode)
	}

	var rkMsg ReEncryptionKeyMessage
	if err := json.NewDecoder(resp.Body).Decode(&rkMsg); err != nil {
		return pre.ReEncryptionKey{}, err
	}

	rk, err := DeSerializeReEncryptionKey(rkMsg.ReEncryptionKeySerialized)
	if err != nil {
		return pre.ReEncryptionKey{}, err
	}

	instanceKeys := keys[rkMsg.InstanceId]
	instanceKeys.ReEncryptionKey = rk
	keys[rkMsg.InstanceId] = instanceKeys
	return rk, nil
}

func getOrSetLeader(functionId FunctionId) (InstanceId, error) {
	if functionId == 0 {
		return "", fmt.Errorf("function ID cannot be 0")
	}
	if functionLeaders[functionId] == "" {
		// in the real implementation there would be some better way to select a leader
		functionLeaders[functionId] = instances[0]
		log.Println("setting alice to function leader")
	}
	leaderId := functionLeaders[functionId]
	return leaderId, nil
}

func getAvailabileInstance(functionId FunctionId) InstanceId {
	//return instances[0] // ALICE
	return instances[1] // BOB
}

func reEncrypt(m1 *SambaMessage, leaderId, instanceId InstanceId) (*SambaMessage, error) {
	rkAB, err := getReEncryptionKey(leaderId, instanceId)
	if err != nil {
		return nil, err
	}

	ct1, err := DeSerializeCiphertext1(m1.WrappedKey1)
	if err != nil {
		return nil, err
	}

	ct2 := pre.ReEncrypt(pp, &rkAB, &ct1)

	wk2, err := SerializeCiphertext2(*ct2)
	if err != nil {
		return nil, err
	}

	m2 := SambaMessage{
		Target:        m1.Target,
		IsReEncrypted: true,
		WrappedKey2:   wk2,
		Ciphertext:    m1.Ciphertext,
	}

	return &m2, nil
}

func recvMessage(w http.ResponseWriter, req *http.Request) {
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

	leaderId, err := getOrSetLeader(m.Target)
	if err != nil {
		http.Error(w, "failed to get or set leader", http.StatusInternalServerError)
		log.Printf("failed to get or set leader: %v", err)
		return
	}

	instanceId := getAvailabileInstance(m.Target)
	if instanceId != leaderId {
		m2, err := reEncrypt(&m, leaderId, instanceId)
		if err != nil {
			http.Error(w, "reEncryption failed", http.StatusInternalServerError)
			log.Printf("reEncryption failed: %v", err)
			return
		}
		m = *m2
	}

	resp, err := SendMessage(&m, instanceId)
	if err != nil {
		http.Error(w, "Message forwarding failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Message forwarding failed: %v", err)
		return
	}

	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Failed to write response body: %v", err)
	}
}

func handlePublicKeyRequest(w http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()
	functionId, err := strconv.ParseUint(queries.Get("functionId"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing string to uint:", err)
		return
	}

	leaderId, err := getOrSetLeader(FunctionId(functionId))
	if err != nil {
		http.Error(w, "Could not get or set leader: %v", http.StatusInternalServerError)
		log.Printf("Could not get or set leader: %v", err)
		return
	}

	leaderKeys, exists := keys[leaderId]
	if !exists {
		http.Error(w, "Function leader has no public key", http.StatusInternalServerError)
		log.Printf("Function leader has no public key for leaderId %s", leaderId)
		return
	}

	msg := SerializePublicKey(leaderKeys.PublicKey)
	jsonData, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "Failed to encode public key", http.StatusInternalServerError)
		log.Printf("Error marshaling public key message: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func BootProxy(instanceIds []InstanceId) {
	instances = instanceIds
	http.HandleFunc("/publicParams", sendPublicParams)
	http.HandleFunc("/registerPublicKey", recvPublicKey)
	http.HandleFunc("/publicKey", handlePublicKeyRequest)
	http.HandleFunc("/message", recvMessage)
	log.Println("Proxy service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
