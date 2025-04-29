package samba

import (
	"os"
	"testing"

	"github.com/etclab/pre"
)

const PLAINTEXT = "Hello from the Samba-protected cloud"

var blackhole any
var pp *pre.PublicParams
var alice *SambaInstance
var bob *SambaInstance

func TestMain(m *testing.M) {
	pp = pre.NewPublicParams()

	alice = &SambaInstance{
		KeyPair: pre.KeyGen(pp),
		Id:      "alice",
	}

	bob = &SambaInstance{
		KeyPair: pre.KeyGen(pp),
		Id:      "bob",
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func BenchmarkEncrypt(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		m, err := Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
		if err != nil {
			b.Fatal(err)
		}
		blackhole = m
	}
}

func BenchmarkDecrypt1(b *testing.B) {
	m, err := Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		plaintext, err := Decrypt(pp, alice.KeyPair.SK, m)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = plaintext
	}
}

func BenchmarkPREDecrypt2(b *testing.B) {
	m, err := Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
	if err != nil {
		b.Fatal(err)
	}

	// Build RK request
	pks := new(PublicKeySerialized)
	pks.Serialize(bob.KeyPair.PK)
	rkReq := ReEncryptionKeyRequest{
		InstanceId:         alice.Id,
		PublicKeySerialzed: *pks,
	}

	// Get RK
	rkMsg, err := GenReEncryptionKey(pp, alice.KeyPair.SK, &rkReq)
	if err != nil {
		b.Fatal(err)
	}
	rkAB, err := rkMsg.ReEncryptionKeySerialized.DeSerialize()
	if err != nil {
		b.Fatal(err)

	}

	// ReEncrypt
	m2, err := ReEncrypt(pp, rkAB, m)

	b.ResetTimer()
	for b.Loop() {
		plaintext, err := Decrypt(pp, bob.KeyPair.SK, m2)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = plaintext
	}
}
