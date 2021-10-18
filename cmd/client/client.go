package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/nwehr/chatterbox"
)

var (
	p   *messages
	inp *input
)

func main() {
	client := client{}

	{
		var ident string
		var to string

		flag.StringVar(&ident, "i", "", "Identity")
		flag.StringVar(&to, "to", "", "To")
		flag.Parse()

		client.ident = chatterbox.Identity(ident)
		client.to = []chatterbox.Identity{chatterbox.Identity(to)}
	}

	fmt.Print("connecting... ")
	if err := client.connect(); err != nil {
		printAndExit(err.Error())
	}
	fmt.Println("done")

	defer client.conn.Close()

	fmt.Print("logging in... ")
	if err := client.login(); err != nil {
		printAndExit(err.Error())
	}
	fmt.Println("done")

	go client.handleMessages()

	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v", err)
		return
	}
	defer ui.Close()

	width, height := ui.TerminalDimensions()

	p = newMessages()
	p.SetRect(0, 0, width, height-5)

	inp = newInput()
	inp.Text = ""
	inp.WrapText = true
	inp.SetRect(0, height-5, width, height)

	ui.Render(p, inp)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-c>":
				return
			case "<C-v>":
				inp.Text = "pasted data"
				inp.cursorLoc = len(inp.Text)
			case "<Enter>":
				if len(inp.Text) > 0 {
					msg := chatterbox.Send(chatterbox.Identity(client.ident), []chatterbox.Identity{client.ident, client.to[0]}, inp.Text)
					if _, err := msg.WriteTo(client.conn); err != nil {
						fmt.Println(err)
					}

					client.msg = []byte("")
				}

				inp.Text = ""
				inp.cursorLoc = 0
			case "<Left>":
				if inp.cursorLoc > 0 {
					inp.cursorLoc -= 1
				}
			case "<Right>":
				if inp.cursorLoc < len(inp.Text) {
					inp.cursorLoc += 1
				}
			case "<Backspace>":
				if len(inp.Text) > 0 {
					inp.Text = inp.Text[:inp.cursorLoc-1] + inp.Text[inp.cursorLoc:]
					inp.cursorLoc -= 1
				}
			case "<Space>":
				inp.Text = inp.Text[:inp.cursorLoc] + " " + inp.Text[inp.cursorLoc:]
				inp.cursorLoc += 1
			default:
				if len(e.ID) == 1 {
					inp.Text = inp.Text[:inp.cursorLoc] + e.ID + inp.Text[inp.cursorLoc:]
					inp.cursorLoc += 1
				}
			}

			ui.Render(p, inp)
		}
	}
}

type client struct {
	ident chatterbox.Identity
	to    []chatterbox.Identity
	conn  net.Conn
	msg   []byte
}

func (c *client) connect() error {
	_, addrs, err := net.LookupSRV("chatterbox-client", "tcp", c.ident.Host())
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addrs[0].Target, addrs[0].Port))
	if err != nil {
		return err
	}

	c.conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) login() error {
	login := chatterbox.Login(c.ident, "")
	if _, err := login.WriteTo(c.conn); err != nil {
		return err
	}

	resp := chatterbox.Message{}
	if _, err := resp.ReadFrom(c.conn); err != nil {
		return err
	}

	if resp.Type != "OK" {
		return fmt.Errorf(resp.Type)
	}

	return nil
}

func (c *client) handleMessages() {
	for {
		msg := chatterbox.Message{}
		if _, err := msg.ReadFrom(c.conn); err != nil {
			fmt.Println("read", err)
			continue
		}

		if msg.Type == "SEND" {
			p.Text += fmt.Sprintf("%s: %s\n", msg.Args["From"][0], string(msg.Data))
			ui.Render(p, inp)
		}
	}
}

func printAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
