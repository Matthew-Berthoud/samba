package samba

import (
	"fmt"

	bls "github.com/cloudflare/circl/ecc/bls12381"
	"github.com/etclab/pre"
)

type SambaPRE struct{}

func (s SambaPRE) Encrypt(pp *pre.PublicParams, pk any, plaintext []byte, functionId FunctionId) (*SambaMessage, error) {
	pkPRE, ok := pk.(*pre.PublicKey)
	if !ok {
		return nil, fmt.Errorf("pk is not a proxy re-encryption PublicKey")
	}
	m := pre.RandomGt()
	ct1 := pre.Encrypt(pp, m, pkPRE)
	key := pre.KdfGtToAes256(m)
	ct := AESGCMEncrypt(key, plaintext)

	var ct1s Ciphertext1Serialized
	err := ct1s.Serialize(ct1)
	if err != nil {
		return nil, err
	}

	return &SambaMessage{
		Target:        functionId,
		IsReEncrypted: false,
		WrappedKey1:   ct1s,
		Ciphertext:    ct,
	}, nil
}

func (s SambaPRE) Decrypt(pp *pre.PublicParams, sk any, m *SambaMessage) ([]byte, error) {
	skPRE, ok := sk.(*pre.SecretKey)
	if !ok {
		return nil, fmt.Errorf("pk is not a proxy re-encryption SecretKey")
	}

	var gt *bls.Gt

	if m.IsReEncrypted {
		ct2, err := m.WrappedKey2.DeSerialize()
		if err != nil {
			return nil, err
		}
		gt = pre.Decrypt2(pp, ct2, skPRE)
	} else {
		ct1, err := m.WrappedKey1.DeSerialize()
		if err != nil {
			return nil, err
		}
		gt = pre.Decrypt1(pp, ct1, skPRE)
	}

	key := pre.KdfGtToAes256(gt)
	plaintext, err := AESGCMDecrypt(key, m.Ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (s SambaPRE) ReEncrypt(pp *pre.PublicParams, rk *pre.ReEncryptionKey, m *SambaMessage) (*SambaMessage, error) {
	ct1, err := m.WrappedKey1.DeSerialize()
	if err != nil {
		return nil, err
	}

	ct2 := pre.ReEncrypt(pp, rk, ct1)

	var wk2 Ciphertext2Serialized
	err = wk2.Serialize(ct2)
	if err != nil {
		return nil, err
	}

	return &SambaMessage{
		Target:        m.Target,
		IsReEncrypted: true,
		WrappedKey2:   wk2,
		Ciphertext:    m.Ciphertext,
	}, nil
}

func (s SambaPRE) GenReEncryptionKey(pp *pre.PublicParams, sk *pre.SecretKey, req *ReEncryptionKeyRequest) (*ReEncryptionKeyMessage, error) {
	pk, err := req.PublicKeySerialzed.DeSerialize()
	if err != nil {
		return nil, err
	}

	rkAB := pre.ReEncryptionKeyGen(pp, sk, pk)
	var rks ReEncryptionKeySerialized
	rks.Serialize(rkAB)

	return &ReEncryptionKeyMessage{
		InstanceId:                req.InstanceId,
		ReEncryptionKeySerialized: rks,
	}, nil
}
