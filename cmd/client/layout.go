package main

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type statusBarManager struct{}

func (m statusBarManager) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	if v, err := g.SetView("status", -1, 0, maxX, 2); err != nil {
		v.Frame = false
		v.Autoscroll = true
		v.BgColor = gocui.ColorCyan
		v.FgColor = gocui.ColorBlack

		toStrs := []string{}

		for _, t := range c.to {
			toStrs = append(toStrs, string(t))
		}

		fmt.Fprintf(v, "chatterbox v0.0.1 %s -> %s", c.ident, strings.Join(toStrs, ";"))

		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

type outputManager struct{}

func (m outputManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("messages", 0, 2, maxX, maxY-2); err != nil {
		v.Frame = false
		v.Autoscroll = true

		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

type inputManager struct{}

func (m inputManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("input", -1, maxY-3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}

		v.Editable = true
		v.Autoscroll = true
		v.BgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorBlack
		v.Frame = false
		v.Wrap = true
	}
	return nil
}
