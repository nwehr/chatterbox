package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/nwehr/chatterbox"
	"github.com/nwehr/chatterbox/pkg/session"
)

var sessions []session.Session

func main() {
	var address, port string

	{
		flag.StringVar(&address, "b", "0.0.0.0", "Bind address, default 0.0.0.0")
		flag.StringVar(&port, "p", "6847", "Listen port, default 6847")
		flag.Parse()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleNewConnection(conn)
	}
}

func handleNewConnection(conn net.Conn) {
	var ident chatterbox.Identity

	{
		req := chatterbox.Request{}
		err := req.Read(conn)
		if err != nil {

		}

		ident = chatterbox.Identity(req.Args["Identity"])

		chatterbox.OKResponse().Write(conn)
	}

	s := session.NewSession(ident, conn)

	defer session.QuitSession(s.ID)

	log.Printf("new session for %s\n", string(ident))

	for {
		req := chatterbox.Request{}
		err := req.Read(conn)

		if err == os.ErrDeadlineExceeded {
			continue
		}

		if err == io.EOF {
			log.Println("connection closed")
			conn.Close()
			break
		}

		if err == net.ErrClosed {
			log.Println("connection closed")
			break
		}

		if err != nil {
			log.Printf("could not read: %s\n", err.Error())
			break
		}

		fmt.Println(req.Type)
	}
}
