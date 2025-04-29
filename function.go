package samba

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/etclab/pre"
)

type SambaInstance struct {
	id InstanceId
	kp *pre.KeyPair
	pp *pre.PublicParams
}

func (s *SambaInstance) genReEncryptionKey(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Failed to read request body: %v", err)
		return
	}

	var rkReq ReEncryptionKeyRequest
	if err := json.Unmarshal(body, &rkReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		log.Printf("Invalid request format: %v", err)
		return
	}

	pk, err := DeSerializePublicKey(rkReq.PublicKeySerialzed)
	if err != nil {
		http.Error(w, "Failed to deserialize public key", http.StatusBadRequest)
		log.Printf("Failed to deserialize public key: %v", err)
		return
	}

	rkAB := pre.ReEncryptionKeyGen(s.pp, s.kp.SK, &pk)
	rks := SerializeReEncryptionKey(*rkAB)
	response := ReEncryptionKeyMessage{
		InstanceId:                rkReq.InstanceId,
		ReEncryptionKeySerialized: rks,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *SambaInstance) handleMessage(w http.ResponseWriter, req *http.Request) {
	HandleMessage(w, req, s.kp, s.pp)
}

func (s *SambaInstance) port() string {
	re := regexp.MustCompile(`:\d+`)
	return re.FindString(string(s.id))
}

func (s *SambaInstance) Boot(selfId, proxyId InstanceId) {
	s.id = selfId
	s.pp = FetchPublicParams(proxyId)
	s.kp = pre.KeyGen(s.pp)

	RegisterPublicKey(proxyId, selfId, s.kp.PK)

	http.HandleFunc("/requestReEncryptionKey", s.genReEncryptionKey)
	http.HandleFunc("/message", s.handleMessage)

	port := s.port()
	log.Println("Alice service running on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
