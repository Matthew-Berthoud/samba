package main

import (
	"github.com/etclab/samba"
)

func main() {
	var aliceId samba.InstanceId = "http://localhost:8081"
	var bobId samba.InstanceId = "http://localhost:8082"

	options := samba.ParseOptions("proxy")

	var c samba.SambaCrypto
	if options.UseRSA {
		c = new(samba.SambaRSA)
	} else {
		c = new(samba.SambaPRE)
	}

	s := samba.SambaProxy{}
	s.Boot([]samba.InstanceId{aliceId, bobId}, c)
}
