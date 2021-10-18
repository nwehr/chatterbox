package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/nwehr/chatterbox"
)

func main() {
	var domain, address, port string

	{
		flag.StringVar(&address, "b", "0.0.0.0", "Bind address, default 0.0.0.0")
		flag.StringVar(&port, "p", "6847", "Listen port, default 6847")
		flag.StringVar(&domain, "d", "home.lan", "Domain for this server, default home.lan")
		flag.Parse()
	}

	serv := &server{
		domain:   domain,
		sessions: map[chatterbox.Identity][]net.Conn{},
	}

	fmt.Println(serv.listen(fmt.Sprintf("%s:%s", address, port)))
}

type server struct {
	sessions map[chatterbox.Identity][]net.Conn
	domain   string
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

	// first message by a client should be LOGIN
	{
		msg := chatterbox.Message{}
		if err := msg.Read(conn); err != nil {
			fmt.Println("could not read login request", err)
			return

		}

		if msg.Type != "LOGIN" {
			fmt.Printf("expected LOGIN; got %s\n", msg.Type)
			return
		}

		if err := chatterbox.Ok().Write(conn); err != nil {
			fmt.Println(err)
		}

		ident = chatterbox.Identity(msg.Args["Identity"][0])
	}

	s.sessions[ident] = append(s.sessions[ident], conn)

	fmt.Printf("new session for %s; %d total\n", string(ident), len(s.sessions[ident]))

	for {
		msg := chatterbox.Message{}
		if err := msg.Read(conn); err != nil {
			if err == os.ErrDeadlineExceeded {
				continue
			}

			if err == io.EOF {
				fmt.Println("connection closed")
				conn.Close()
				break
			}

			if err == net.ErrClosed {
				fmt.Println("connection closed")
				break
			}

			if err != nil {
				fmt.Printf("could not read: %s\n", err.Error())
				break
			}
		}

		s.handleMessage(msg)
	}
}

func (s *server) handleMessage(msg chatterbox.Message) {
	switch msg.Type {
	case "SEND":
		for _, to := range msg.Args["To"] {
			ident := chatterbox.Identity(to)

			if ident.Host() == s.domain {
				for _, conn := range s.sessions[ident] {
					go msg.Write(conn)
				}
			}
		}
	}
}
