package main

import (
	"github.com/etclab/samba"
)

func main() {
	var bobId samba.InstanceId = "http://localhost:8082"
	var proxyId samba.InstanceId = "http://localhost:8080"

	options := samba.ParseOptions("bob")

	var c samba.SambaCrypto
	if options.UseRSA {
		c = new(samba.SambaRSA)
	} else {
		c = new(samba.SambaPRE)
	}

	s := samba.SambaInstance{}
	s.Boot(bobId, proxyId, c)
}
