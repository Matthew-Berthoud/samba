package samba

import (
	"log"

	bls "github.com/cloudflare/circl/ecc/bls12381"
	"github.com/etclab/pre"
)

func PREEncrypt(pp *pre.PublicParams, pk *pre.PublicKey, plaintext []byte, functionId FunctionId) (*SambaMessage, error) {
	m := pre.RandomGt()
	ct1 := pre.Encrypt(pp, m, pk)
	key := pre.KdfGtToAes256(m)
	ct := AESGCMEncrypt(key, plaintext)
	ct1s, err := SerializeCiphertext1(*ct1)
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

func PREDecrypt(pp *pre.PublicParams, sk *pre.SecretKey, m *SambaMessage) ([]byte, error) {
	var gt *bls.Gt
	if m.IsReEncrypted {
		ct2, err := DeSerializeCiphertext2(m.WrappedKey2)
		if err != nil {
			return nil, err
		}
		log.Printf("Decrypting with method 2, with re-encryption")
		gt = pre.Decrypt2(pp, &ct2, sk)
	} else {
		ct1, err := DeSerializeCiphertext1(m.WrappedKey1)
		if err != nil {
			return nil, err
		}
		log.Printf("Decrypting with method 1, without re-encryption")
		gt = pre.Decrypt1(pp, &ct1, sk)
	}
	key := pre.KdfGtToAes256(gt)
	plaintext, err := AESGCMDecrypt(key, m.Ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
