package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/nwehr/chatterbox"
)

var c client

func main() {
	c = client{
		statusBarManager: statusBarManager{},
		outputManager:    outputManager{},
		inputManager:     inputManager{},
	}

	{
		var ident string
		var to string

		flag.StringVar(&ident, "i", "", "Identity")
		flag.StringVar(&to, "to", "", "To")
		flag.Parse()

		c.ident = chatterbox.Identity(ident)
		c.to = []chatterbox.Identity{chatterbox.Identity(to)}
	}

	if err := c.connect(); err != nil {
		printAndExit(err.Error())
	}

	defer c.conn.Close()

	if err := c.login(); err != nil {
		printAndExit(err.Error())
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SetManager(c.statusBarManager, c.outputManager, c.inputManager)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, c.handleInput); err != nil {
		log.Panicln(err)
	}

	go c.handleMessages(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

type client struct {
	ident            chatterbox.Identity
	to               []chatterbox.Identity
	conn             net.Conn
	statusBarManager gocui.Manager
	inputManager     gocui.Manager
	outputManager    gocui.Manager
}

func (c *client) connect() error {
	fmt.Printf("connecting... ")

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

	fmt.Println("done")

	return nil
}

func (c *client) login() error {
	fmt.Printf("logging in... ")

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

	fmt.Println("done")

	return nil
}

func (c *client) handleMessages(g *gocui.Gui) {
	for {
		msg := chatterbox.Message{}
		if _, err := msg.ReadFrom(c.conn); err != nil {
			break
		}

		if msg.Type == "SEND" {
			g.Update(func(g *gocui.Gui) error {
				messages, err := g.View("messages")
				if err != nil {
					return err
				}

				if msg.Args["From"][0] == c.ident.String() {
					fmt.Fprintf(messages, "\033[36;1;1m%s\033[0m %s\n", msg.Args["From"][0], string(msg.Data))
				} else {
					fmt.Fprintf(messages, "\033[35;1;1m%s\033[0m %s\n", msg.Args["From"][0], string(msg.Data))
				}
				return nil
			})
		}
	}
}

func (c *client) handleInput(g *gocui.Gui, input *gocui.View) error {
	data := strings.TrimSpace(input.Buffer())

	if len(data) > 0 {
		msg := chatterbox.Send(chatterbox.Identity(c.ident), []chatterbox.Identity{c.ident, c.to[0]}, data)
		if _, err := msg.WriteTo(c.conn); err != nil {
			fmt.Println(err)
		}
	}

	input.Clear()
	input.SetCursor(0, 0)

	return nil
}

func printAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
