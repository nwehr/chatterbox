package main

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type messages struct {
	widgets.Paragraph
}

func newMessages() *messages {
	return &messages{
		Paragraph: *widgets.NewParagraph(),
	}
}

func (self *messages) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	cells := ui.ParseStyles(self.Text, self.TextStyle)
	if self.WrapText {
		cells = ui.WrapCells(cells, uint(self.Inner.Dx()))
	}

	rows := ui.SplitCells(cells, '\n')

	if len(rows) >= self.Inner.Max.Y {
		rows = rows[len(rows)-self.Inner.Max.Y+1:]
	}

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
}
