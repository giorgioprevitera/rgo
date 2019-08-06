package main

import (
	"github.com/jroimartin/gocui"
)

func setKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, openPost); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'r', gocui.ModNone, populateView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, openPost); err != nil {
		return err
	}

	return nil
}
