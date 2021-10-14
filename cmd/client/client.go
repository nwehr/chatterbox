package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/nwehr/chatterbox/pkg/identity"
)

func main() {
	var ident string
	var message string
	var to string

	flag.StringVar(&ident, "i", "", "Identity")
	flag.StringVar(&message, "m", "", "Message")
	flag.StringVar(&to, "to", "", "To")
	flag.Parse()

	if ident == "" {
		printAndExit("no identity provided")
	}

	_, addrs, err := net.LookupSRV("chatterbox-client", "tcp", identity.Identity(ident).Host())
	if err != nil {
		printAndExit(err.Error())
	}

	fmt.Printf("connecting to %s:%d\n", addrs[0].Target, addrs[0].Port)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addrs[0].Target, addrs[0].Port))
	if err != nil {
		printAndExit(err.Error())
	}

	fmt.Printf("connecting to %s\n", addr.String())

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		printAndExit(err.Error())
	}

	defer conn.Close()

	_, err = conn.Write([]byte(ident))
	if err != nil {
		printAndExit(err.Error())
	}
}

func printAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
