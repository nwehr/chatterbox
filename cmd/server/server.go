package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/nwehr/chatterbox/pkg/identity"
)

func main() {
	var address, port string

	flag.StringVar(&address, "a", "0.0.0.0", "Bind address, default 0.0.0.0")
	flag.StringVar(&port, "p", "6847", "Listen port, default 6847")

	flag.Parse()

	listen := fmt.Sprintf("%s:%s", address, port)

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			ident := identity.Identity("")

			{
				buff := make([]byte, 256)
				_, err := conn.Read(buff)
				if err != nil {

				}

				ident = identity.Identity(buff)
			}

			log.Printf("new session for %s\n", string(ident))

			for {
				buff := make([]byte, 2048)
				conn.SetReadDeadline(time.Now().Add(time.Second * 30))
				_, err := conn.Read(buff)

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

				fmt.Println(string(buff))

				conn.Write([]byte("Ok"))
			}
		}(conn)
	}
}
