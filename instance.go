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
	Id           InstanceId
	KeyPair      *pre.KeyPair
	PublicParams *pre.PublicParams
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

	m, err := GenReEncryptionKey(s.PublicParams, s.KeyPair.SK, &rkReq)
	if err != nil {
		http.Error(w, "Failed to generate re-encryption key", http.StatusInternalServerError)
		log.Printf("Failed to generate re-encryption key: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *SambaInstance) handleMessage(w http.ResponseWriter, req *http.Request) {
	HandleMessage(w, req, s.KeyPair, s.PublicParams)
}

func (s *SambaInstance) port() string {
	re := regexp.MustCompile(`:\d+`)
	return re.FindString(string(s.Id))
}

func (s *SambaInstance) Boot(selfId, proxyId InstanceId) {
	s.Id = selfId
	s.PublicParams = FetchPublicParams(proxyId)
	s.KeyPair = pre.KeyGen(s.PublicParams)

	RegisterPublicKey(proxyId, selfId, s.KeyPair.PK)

	http.HandleFunc("/requestReEncryptionKey", s.genReEncryptionKey)
	http.HandleFunc("/message", s.handleMessage)

	port := s.port()
	log.Println("Alice service running on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
