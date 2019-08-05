package main

import (
	"fmt"
	"log"

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

func openPost(g *gocui.Gui, v *gocui.View) error {
	log.Println("opening post")
	myView, _ := g.View("post")
	myView.Clear()
	l := getData()
	mainView, _ := g.View("main")
	_, cy := mainView.Cursor()
	fmt.Fprintf(myView, "Title: %s\n\n", l.Children[cy].Data["title"])
	fmt.Fprintf(myView, "Author: %s\n\n", l.Children[cy].Data["author"])
	fmt.Fprintf(myView, "URL: %s\n\n", l.Children[cy].Data["url"])
	fmt.Fprintf(myView, "%s\n\n\n", l.Children[cy].Data["selftext"])
	fmt.Fprintf(myView, "%s", l.Children[cy].Data)
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	myView, _ := g.View("main")
	cx, cy := myView.Cursor()
	if cy == 0 {
		return nil
	}
	if err := myView.SetCursor(cx, cy-1); err != nil {
		return err
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cy+1 >= len(v.BufferLines())-1 {
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
