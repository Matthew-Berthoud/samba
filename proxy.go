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

const FUNCTION_ID FunctionId = 123

type SambaProxy struct {
	pp                *pre.PublicParams
	functionInstances map[FunctionId][]InstanceId
	instanceKeys      map[InstanceId]InstanceKeys
	functionLeaders   map[FunctionId]InstanceId
}

func (s *SambaProxy) recvPublicKey(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	pks := new(PublicKeySerialized)
	err := json.NewDecoder(req.Body).Decode(&pks)
	if err != nil {
		log.Printf("Failed to decode public key: %v", err)
		http.Error(w, "Failed to decode public key", http.StatusBadRequest)
		return
	}

	pk, err := pks.DeSerialize()
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

func (s *SambaProxy) setPublicKey(instanceId InstanceId, pk *pre.PublicKey) {
	s.instanceKeys[instanceId] = InstanceKeys{
		PublicKey:       pk,
		ReEncryptionKey: s.instanceKeys[instanceId].ReEncryptionKey, // Preserve existing re-encryption key if resetting
	}
}

func (s *SambaProxy) sendPublicParams(w http.ResponseWriter, req *http.Request) {
	pps := new(PublicParamsSerialized)
	err := pps.Serialize(s.pp)
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

func (s *SambaProxy) requestReEncryptionKey(a, b InstanceId) (*pre.ReEncryptionKey, error) {
	pks := new(PublicKeySerialized)
	pks.Serialize(s.instanceKeys[b].PublicKey)
	req := ReEncryptionKeyRequest{
		InstanceId:         b,
		PublicKeySerialzed: *pks,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(string(a)+"/requestReEncryptionKey", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var rkMsg ReEncryptionKeyMessage
	if err := json.NewDecoder(resp.Body).Decode(&rkMsg); err != nil {
		return nil, err
	}

	rk, err := rkMsg.ReEncryptionKeySerialized.DeSerialize()
	if err != nil {
		return nil, err
	}

	instanceKeys := s.instanceKeys[rkMsg.InstanceId]
	instanceKeys.ReEncryptionKey = rk
	s.instanceKeys[rkMsg.InstanceId] = instanceKeys
	return rk, nil
}

func (s *SambaProxy) getOrSetLeader(functionId FunctionId) (InstanceId, error) {
	if functionId == 0 {
		return "", fmt.Errorf("function ID cannot be 0")
	}
	if s.functionLeaders[functionId] == "" {
		// in the real implementation there would be some better way to select a leader
		s.functionLeaders[functionId] = s.functionInstances[FUNCTION_ID][0]
		log.Println("setting alice to function leader")
	}
	leaderId := s.functionLeaders[functionId]
	return leaderId, nil
}

func (s *SambaProxy) getAvailabileInstance(functionId FunctionId) InstanceId {
	// return s.functionInstances[functionId][0] // ALICE
	return s.functionInstances[functionId][1] // BOB
}

func (s *SambaProxy) reEncrypt(m1 *SambaMessage, leaderId, instanceId InstanceId) (*SambaMessage, error) {
	rk := s.instanceKeys[instanceId].ReEncryptionKey
	var err error
	if rk == nil {
		rk, err = s.requestReEncryptionKey(leaderId, instanceId)
		if err != nil {
			return nil, err
		}
	}

	m2, err := ReEncrypt(s.pp, rk, m1)
	if err != nil {
		return nil, err
	}

	return m2, nil
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

	leaderKeys, exists := s.instanceKeys[leaderId]
	if !exists {
		http.Error(w, "Function leader has no public key", http.StatusInternalServerError)
		log.Printf("Function leader has no public key for leaderId %s", leaderId)
		return
	}

	pks := new(PublicKeySerialized)
	pks.Serialize(leaderKeys.PublicKey)
	jsonData, err := json.Marshal(pks)
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
	s.functionInstances = make(map[FunctionId][]InstanceId)
	s.functionInstances[FUNCTION_ID] = instanceIds
	s.functionLeaders = make(map[FunctionId]InstanceId)
	s.instanceKeys = make(map[InstanceId]InstanceKeys)

	http.HandleFunc("/publicParams", s.sendPublicParams)
	http.HandleFunc("/registerPublicKey", s.recvPublicKey)
	http.HandleFunc("/publicKey", s.handlePublicKeyRequest)
	http.HandleFunc("/message", s.recvMessage)
	log.Println("Proxy service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
