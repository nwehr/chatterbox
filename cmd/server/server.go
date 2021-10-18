package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/nwehr/chatterbox"
)

// var sessions []chatterbox.Session

type server struct {
	sessions map[string][]net.Conn
}

func (s *server) listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Println("listening on ", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go s.handleNewConnection(conn)
	}
}

func (s *server) handleNewConnection(conn net.Conn) {
	ident := chatterbox.Identity("")

	// first request by client should be LOGIN
	{
		login := chatterbox.Request{}
		if err := login.Read(conn); err != nil {
			fmt.Println("could not read login request", err)
			return

		}

		if login.Type != "LOGIN" {
			fmt.Printf("expected LOGIN; got %s\n", login.Type)
			return
		}

		chatterbox.OKResponse().Write(conn)

		ident = chatterbox.Identity(login.Args["Identity"][0])
	}

	s.sessions[string(ident)] = append(s.sessions[string(ident)], conn)

	fmt.Println(len(s.sessions[string(ident)]), " sessions for ", string(ident))

	for {
		req := chatterbox.Request{}
		if err := req.Read(conn); err != nil {
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
		}

		s.handleRequest(req)
	}
}

func (s *server) handleRequest(req chatterbox.Request) {
	fmt.Println("sessions ", s.sessions)

	switch req.Type {
	case "SEND":
		for _, to := range req.Args["To"] {
			fmt.Println("forwarding to ", string(to))
			fmt.Println(len(s.sessions[to]), " sessions")
			for _, conn := range s.sessions[to] {
				fmt.Println("writing to ", string(to))
				go req.Write(conn)
			}
		}
	}
}

func main() {
	var address, port string

	{
		flag.StringVar(&address, "b", "0.0.0.0", "Bind address, default 0.0.0.0")
		flag.StringVar(&port, "p", "6847", "Listen port, default 6847")
		flag.Parse()
	}

	serv := &server{}
	serv.sessions = map[string][]net.Conn{}

	log.Fatal(serv.listen(fmt.Sprintf("%s:%s", address, port)))

}
