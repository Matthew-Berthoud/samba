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

type SambaProxy struct {
	pp              *pre.PublicParams
	instances       []InstanceId
	keys            map[InstanceId]InstanceKeys
	functionLeaders map[FunctionId]InstanceId
}

func (s *SambaProxy) recvPublicKey(w http.ResponseWriter, req *http.Request) {
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
	s.setPublicKey(instanceId, pk)
	log.Printf("Successfully storing public key for instanceId: %s", instanceId)

	w.WriteHeader(http.StatusOK)
}

func (s *SambaProxy) setPublicKey(instanceId InstanceId, pk pre.PublicKey) {
	s.keys[instanceId] = InstanceKeys{
		PublicKey:       pk,
		ReEncryptionKey: s.keys[instanceId].ReEncryptionKey, // Preserve existing re-encryption key if resetting
	}
}

func (s *SambaProxy) sendPublicParams(w http.ResponseWriter, req *http.Request) {
	pps, err := SerializePublicParams(*s.pp)
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

func (s *SambaProxy) getReEncryptionKey(a, b InstanceId) (pre.ReEncryptionKey, error) {
	if s.keys[b].ReEncryptionKey != (pre.ReEncryptionKey{}) {
		return s.keys[b].ReEncryptionKey, nil
	}

	pks := SerializePublicKey(s.keys[b].PublicKey)

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

	instanceKeys := s.keys[rkMsg.InstanceId]
	instanceKeys.ReEncryptionKey = rk
	s.keys[rkMsg.InstanceId] = instanceKeys
	return rk, nil
}

func (s *SambaProxy) getOrSetLeader(functionId FunctionId) (InstanceId, error) {
	if functionId == 0 {
		return "", fmt.Errorf("function ID cannot be 0")
	}
	if s.functionLeaders[functionId] == "" {
		// in the real implementation there would be some better way to select a leader
		s.functionLeaders[functionId] = s.instances[0]
		log.Println("setting alice to function leader")
	}
	leaderId := s.functionLeaders[functionId]
	return leaderId, nil
}

func (s *SambaProxy) getAvailabileInstance(functionId FunctionId) InstanceId {
	//return instances[0] // ALICE
	return s.instances[1] // BOB
}

func (s *SambaProxy) reEncrypt(m1 *SambaMessage, leaderId, instanceId InstanceId) (*SambaMessage, error) {
	rkAB, err := s.getReEncryptionKey(leaderId, instanceId)
	if err != nil {
		return nil, err
	}

	ct1, err := DeSerializeCiphertext1(m1.WrappedKey1)
	if err != nil {
		return nil, err
	}

	ct2 := pre.ReEncrypt(s.pp, &rkAB, &ct1)

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

func (s *SambaProxy) recvMessage(w http.ResponseWriter, req *http.Request) {
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

	leaderId, err := s.getOrSetLeader(m.Target)
	if err != nil {
		http.Error(w, "failed to get or set leader", http.StatusInternalServerError)
		log.Printf("failed to get or set leader: %v", err)
		return
	}

	instanceId := s.getAvailabileInstance(m.Target)
	if instanceId != leaderId {
		m2, err := s.reEncrypt(&m, leaderId, instanceId)
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

func (s *SambaProxy) handlePublicKeyRequest(w http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()
	functionId, err := strconv.ParseUint(queries.Get("functionId"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing string to uint:", err)
		return
	}

	leaderId, err := s.getOrSetLeader(FunctionId(functionId))
	if err != nil {
		http.Error(w, "Could not get or set leader: %v", http.StatusInternalServerError)
		log.Printf("Could not get or set leader: %v", err)
		return
	}

	leaderKeys, exists := s.keys[leaderId]
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

func (s *SambaProxy) Boot(instanceIds []InstanceId) {
	s.pp = pre.NewPublicParams()
	s.instances = instanceIds
	s.functionLeaders = make(map[FunctionId]InstanceId)
	s.keys = make(map[InstanceId]InstanceKeys)

	http.HandleFunc("/publicParams", s.sendPublicParams)
	http.HandleFunc("/registerPublicKey", s.recvPublicKey)
	http.HandleFunc("/publicKey", s.handlePublicKeyRequest)
	http.HandleFunc("/message", s.recvMessage)
	log.Println("Proxy service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
