package samba

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const USAGE = `Usage: %s [options]

Process text with Samba cryptographic protections.

Options:
  -rsa  Use SambaRSA encryption instead of default PRE

Examples:
  $ ./%s        # Uses PRE by default
  $ ./%s -rsa   # Uses RSA encryption
`

type Options struct {
	UseRSA bool
}

func ParseOptions(name string) *Options {
	options := &Options{}
	usage := fmt.Sprintf(USAGE, name, name, name)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	flag.BoolVar(&options.UseRSA, "rsa", false, "Use SambaRSA encryption")
	flag.Parse()

	if flag.NArg() != 0 {
		log.Fatalf("Unexpected positional arguments: %v", flag.Args())
	}

	return options
}
