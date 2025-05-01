package samba

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/etclab/pre"
)

type SambaRSA struct{}

type SambaRSAPlaintext struct {
	AesKey        []byte
	AesCiphertext []byte
}

func (s SambaRSA) Encrypt(pp *pre.PublicParams, pk any, plaintext []byte, functionId FunctionId) (*SambaMessage, error) {
	pkRSA, ok := pk.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("sk is not a RSA PublicKey")
	}

	aesKey := make([]byte, 32)
	aesCiphertext := AESGCMEncrypt(aesKey, plaintext)

	pt := &SambaRSAPlaintext{
		AesKey:        aesKey,
		AesCiphertext: aesCiphertext,
	}

	ptEncoded, err := json.Marshal(pt)
	if err != nil {
		return nil, err
	}

	return &SambaMessage{
		Target:        functionId,
		IsReEncrypted: false,
		Ciphertext:    RSAEncrypt(pkRSA, ptEncoded),
	}, nil
}

func (s SambaRSA) Decrypt(pp *pre.PublicParams, sk any, m *SambaMessage) ([]byte, error) {
	skRSA, ok := sk.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("sk is not a RSA PrivateKey")
	}
	ptEncoded, err := RSADecrypt(skRSA, m.Ciphertext)
	if err != nil {
		return nil, err
	}

	var pt SambaRSAPlaintext
	err = json.Unmarshal(ptEncoded, &pt)
	if err != nil {
		return nil, err
	}

	plaintext, err := AESGCMDecrypt(pt.AesKey, pt.AesCiphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (s SambaRSA) ReEncrypt(pp *pre.PublicParams, rk *pre.ReEncryptionKey, m *SambaMessage) (*SambaMessage, error) {
	return nil, fmt.Errorf("ReEncrypt not implemented for SambaRSA")
}

func (s SambaRSA) GenReEncryptionKey(pp *pre.PublicParams, sk *pre.SecretKey, req *ReEncryptionKeyRequest) (*ReEncryptionKeyMessage, error) {
	return nil, fmt.Errorf("GenReEncryptionKey not implemented for SambaRSA")
}
