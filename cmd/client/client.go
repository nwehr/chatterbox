package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/nwehr/chatterbox"
)

type options struct {
	Identity chatterbox.Identity
}

func main() {
	var options options

	{
		var ident string
		var message string
		var to string

		flag.StringVar(&ident, "i", "", "Identity")
		flag.StringVar(&message, "m", "", "Message")
		flag.StringVar(&to, "to", "", "To")
		flag.Parse()

		options.Identity = chatterbox.Identity(ident)
	}

	_, addrs, err := net.LookupSRV("chatterbox-client", "tcp", options.Identity.Host())
	if err != nil {
		printAndExit(err.Error())
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addrs[0].Target, addrs[0].Port))
	if err != nil {
		printAndExit(err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		printAndExit(err.Error())
	}

	defer conn.Close()

	login := chatterbox.LoginRequest(options.Identity, "")

	if err = login.Write(conn); err != nil {
		printAndExit(err.Error())
	}

	resp := chatterbox.Response{}
	resp.Read(conn)

	fmt.Println(resp.Type)
}

func printAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
