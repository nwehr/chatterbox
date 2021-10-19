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

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SetManager(messageLogManager{}, inputManager{})

	go client.handleMessages(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, client.handleEnter); err != nil {
		log.Panicln(err)
	}

	// if messages, err := g.View("messages"); err != nil {
	// 	fmt.Fprintf(messages, "Welcome to chatterbox v0.0.1\n")
	// 	fmt.Fprintf(messages, "You are logged in as %s\n", client.ident)
	// }

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type messageLogManager struct {
}

func (m messageLogManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("messages", 0, 0, maxX, maxY-2); err != nil {
		v.Frame = false
		v.Autoscroll = true

		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

type inputManager struct {
}

func (m inputManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("input", 0, maxY-3, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}

		v.Editable = true
		v.Autoscroll = true
		v.Frame = false
		v.Wrap = true
	}
	return nil
}

func (c client) handleEnter(g *gocui.Gui, input *gocui.View) error {
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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

func (c *client) handleMessages(g *gocui.Gui) {
	for {
		msg := chatterbox.Message{}
		if _, err := msg.ReadFrom(c.conn); err != nil {
			fmt.Println("read", err)
			continue
		}

		if msg.Type == "SEND" {
			g.Update(func(g *gocui.Gui) error {
				messages, err := g.View("messages")
				if err != nil {
					return err
				}

				fmt.Fprintf(messages, "\033[1m%s\033[0m %s\n", msg.Args["From"][0], string(msg.Data))
				return nil
			})
		}
	}
}

func printAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

// func main() {

// 	if err := ui.Init(); err != nil {
// 		fmt.Printf("failed to initialize termui: %v", err)
// 		return
// 	}
// 	defer ui.Close()

// 	width, height := ui.TerminalDimensions()

// 	p = newMessages()
// 	p.SetRect(0, 0, width, height-5)

// 	inp = newInput()
// 	inp.Text = ""
// 	inp.WrapText = true
// 	inp.SetRect(0, height-5, width, height)

// 	ui.Render(p, inp)

// 	for e := range ui.PollEvents() {
// 		if e.Type == ui.KeyboardEvent {
// 			switch e.ID {
// 			case "<C-c>":
// 				return
// 			case "<C-v>":
// 				inp.Text = "pasted data"
// 				inp.cursorLoc = len(inp.Text)
// 			case "<Enter>":
// 				if len(inp.Text) > 0 {
// 					msg := chatterbox.Send(chatterbox.Identity(client.ident), []chatterbox.Identity{client.ident, client.to[0]}, inp.Text)
// 					if _, err := msg.WriteTo(client.conn); err != nil {
// 						fmt.Println(err)
// 					}

// 					client.msg = []byte("")
// 				}

// 				inp.Text = ""
// 				inp.cursorLoc = 0
// 			case "<Left>":
// 				if inp.cursorLoc > 0 {
// 					inp.cursorLoc -= 1
// 				}
// 			case "<Right>":
// 				if inp.cursorLoc < len(inp.Text) {
// 					inp.cursorLoc += 1
// 				}
// 			case "<Backspace>":
// 				if len(inp.Text) > 0 {
// 					inp.Text = inp.Text[:inp.cursorLoc-1] + inp.Text[inp.cursorLoc:]
// 					inp.cursorLoc -= 1
// 				}
// 			case "<Space>":
// 				inp.Text = inp.Text[:inp.cursorLoc] + " " + inp.Text[inp.cursorLoc:]
// 				inp.cursorLoc += 1
// 			default:
// 				if len(e.ID) == 1 {
// 					inp.Text = inp.Text[:inp.cursorLoc] + e.ID + inp.Text[inp.cursorLoc:]
// 					inp.cursorLoc += 1
// 				}
// 			}

// 			ui.Render(p, inp)
// 		}
// 	}
// }
