package main

import (
	"image"
	"math"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	windowHeight = 20
	windowWidth  = 40
)

type input struct {
	widgets.Paragraph
	cursorLoc int
}

func (self *input) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	cells := ui.ParseStyles(self.Text, self.TextStyle)
	if self.WrapText {
		cells = ui.WrapCells(cells, uint(self.Inner.Dx()))
	}

	rows := ui.SplitCells(cells, '\n')

	for y, row := range rows {
		if y+self.Inner.Min.Y >= self.Inner.Max.Y {
			break
		}
		row = ui.TrimCells(row, self.Inner.Dx())
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(self.Inner.Min))
		}
	}

	{
		ch := ' '

		if self.cursorLoc < len(self.Text) && len(self.Text) > 0 {
			ch = rune(self.Text[self.cursorLoc])
		}

		cursor := ui.Cell{
			Rune:  ch,
			Style: ui.NewStyle(ui.ColorBlack, ui.ColorWhite),
		}

		x := self.cursorLoc % self.Inner.Max.X
		y := int(math.Floor(float64(self.cursorLoc) / float64(self.Inner.Max.X)))

		buf.SetCell(cursor, image.Pt(x, y).Add(self.Inner.Min))
	}
}

func newInput() *input {
	return &input{
		Paragraph: *widgets.NewParagraph(),
	}
}

// if err := ui.Init(); err != nil {
// 	log.Fatalf("failed to initialize termui: %v", err)
// }
// defer ui.Close()

// width, height := ui.TerminalDimensions()

// inp := newInput()
// inp.Text = ""
// inp.WrapText = true
// inp.SetRect(0, height-5, width, height)

// ui.Render(inp)

// for e := range ui.PollEvents() {
// 	if e.Type == ui.KeyboardEvent {
// 		switch e.ID {
// 		case "<C-c>":
// 			return
// 		case "<C-v>":
// 			inp.Text = "pasted data"
// 			inp.cursorLoc = len(inp.Text)
// 		case "<Enter>":
// 			inp.Text = ""
// 			inp.cursorLoc = 0
// 		case "<Left>":
// 			if inp.cursorLoc > 0 {
// 				inp.cursorLoc -= 1
// 			}
// 		case "<Right>":
// 			if inp.cursorLoc < len(inp.Text) {
// 				inp.cursorLoc += 1
// 			}
// 		case "<Backspace>":
// 			if len(inp.Text) > 0 {
// 				inp.Text = inp.Text[:inp.cursorLoc-1] + inp.Text[inp.cursorLoc:]
// 				inp.cursorLoc -= 1
// 			}
// 		case "<Space>":
// 			inp.Text = inp.Text[:inp.cursorLoc] + " " + inp.Text[inp.cursorLoc:]
// 			inp.cursorLoc += 1
// 		default:
// 			if len(e.ID) == 1 {
// 				inp.Text = inp.Text[:inp.cursorLoc] + e.ID + inp.Text[inp.cursorLoc:]
// 				inp.cursorLoc += 1
// 			}
// 		}

// 		ui.Render(inp)
// 	}
// }
// conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5051})
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println("connected")

// go func(conn net.Conn) {
// 	for {
// 		buff := make([]byte, 2048)
// 		_, err = conn.Read(buff)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		fmt.Println(string(buff))
// 	}

// }(conn)

// reader := bufio.NewReader(os.Stdin)

// for {
// 	fmt.Print("> ")
// 	text, _ := reader.ReadString('\n')

// 	// convert CRLF to LF
// 	text = strings.Replace(text, "\n", "", -1)

// 	_, err := conn.Write([]byte(text))
// 	if err != nil {
// 		log.Println("could not write: %s", err.Error())
// 	}
// }
