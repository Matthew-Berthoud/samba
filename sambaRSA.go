package samba

import (
	"crypto/rsa"
	"fmt"

	"github.com/etclab/pre"
)

type SambaRSA struct{}

func (s SambaRSA) Encrypt(pp *pre.PublicParams, pk any, plaintext []byte, functionId FunctionId) (*SambaMessage, error) {
	pkRSA, ok := pk.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("sk is not a RSA PublicKey")
	}
	return &SambaMessage{
		Target:        functionId,
		IsReEncrypted: false,
		Ciphertext:    RSAEncrypt(pkRSA, plaintext),
	}, nil
}

func (s SambaRSA) Decrypt(pp *pre.PublicParams, sk any, m *SambaMessage) ([]byte, error) {
	skRSA, ok := sk.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("sk is not a RSA PrivateKey")
	}
	plaintext, err := RSADecrypt(skRSA, m.Ciphertext)
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
