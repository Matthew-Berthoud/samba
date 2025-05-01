package main

import (
	"github.com/etclab/samba"
)

func main() {
	var proxyId samba.InstanceId = "http://localhost:8080"
	var aliceId samba.InstanceId = "http://localhost:8081"

	options := samba.ParseOptions("alice")

	var c samba.SambaCrypto
	if options.UseRSA {
		c = new(samba.SambaRSA)
	} else {
		c = new(samba.SambaPRE)
	}

	s := samba.SambaInstance{}
	s.Boot(aliceId, proxyId, c)
}
