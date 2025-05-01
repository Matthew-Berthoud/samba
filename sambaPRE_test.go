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
var sambaPRE *SambaPRE

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

	sambaPRE = new(SambaPRE)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func BenchmarkEncrypt(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		m, err := sambaPRE.Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
		if err != nil {
			b.Fatal(err)
		}
		blackhole = m
	}
}

func BenchmarkDecrypt1(b *testing.B) {
	m, err := sambaPRE.Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		plaintext, err := sambaPRE.Decrypt(pp, alice.KeyPair.SK, m)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = plaintext
	}
}

func BenchmarkGenReEncryptionKey(b *testing.B) {
	// Build RK request
	pks := new(PublicKeySerialized)
	pks.Serialize(bob.KeyPair.PK)
	rkReq := ReEncryptionKeyRequest{
		InstanceId:         alice.Id,
		PublicKeySerialzed: *pks,
	}

	b.ResetTimer()
	for b.Loop() {
		rkMsg, err := sambaPRE.GenReEncryptionKey(pp, alice.KeyPair.SK, &rkReq)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = rkMsg
	}
}

func BenchmarkReEncrypt(b *testing.B) {
	m, err := sambaPRE.Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
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
	rkMsg, err := sambaPRE.GenReEncryptionKey(pp, alice.KeyPair.SK, &rkReq)
	if err != nil {
		b.Fatal(err)
	}
	rkAB, err := rkMsg.ReEncryptionKeySerialized.DeSerialize()
	if err != nil {
		b.Fatal(err)

	}

	b.ResetTimer()
	for b.Loop() {
		m2, err := sambaPRE.ReEncrypt(pp, rkAB, m)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = m2
	}
}

func BenchmarkDecrypt2(b *testing.B) {
	m, err := sambaPRE.Encrypt(pp, alice.KeyPair.PK, []byte(PLAINTEXT), FunctionId(1))
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
	rkMsg, err := sambaPRE.GenReEncryptionKey(pp, alice.KeyPair.SK, &rkReq)
	if err != nil {
		b.Fatal(err)
	}
	rkAB, err := rkMsg.ReEncryptionKeySerialized.DeSerialize()
	if err != nil {
		b.Fatal(err)

	}

	// ReEncrypt
	m2, err := sambaPRE.ReEncrypt(pp, rkAB, m)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		plaintext, err := sambaPRE.Decrypt(pp, bob.KeyPair.SK, m2)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = plaintext
	}
}
