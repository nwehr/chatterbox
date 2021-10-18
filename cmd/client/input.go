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
