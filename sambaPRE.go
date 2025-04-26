package samba

import (
	"log"

	"github.com/etclab/pre"
)

func PREEncrypt(pp *pre.PublicParams, pk *pre.PublicKey, plaintext []byte, functionId FunctionId) SambaMessage {
	m := pre.RandomGt()
	ct1 := pre.Encrypt(pp, m, pk)
	key := pre.KdfGtToAes256(m)
	ct := AESGCMEncrypt(key, plaintext)
	ct1s, err := SerializeCiphertext1(*ct1)
	if err != nil {
		log.Fatalf("Failed to serialize: %v", err)
	}
	return SambaMessage{
		Target:        functionId,
		IsReEncrypted: false,
		WrappedKey1:   ct1s,
		Ciphertext:    ct,
	}
}
