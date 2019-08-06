package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("post", maxX/2, 0, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		log.Println("init post view")
		v.Title = "post"
		v.Highlight = true
		v.Wrap = true
		v.Editable = true

	}
	v, err := g.SetView("main", 0, 0, maxX/2, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		log.Println("init main view")
		v.Title = "main"
		v.Highlight = true
		v.Wrap = false
		v.Editable = true
		populateView(g, v)

	}
	g.SetCurrentView("main")
	return nil
}

func openPost(g *gocui.Gui, v *gocui.View) error {
	log.Println("opening post")
	myView, _ := g.View("post")
	myView.Clear()
	mainView, _ := g.View("main")
	_, cy := mainView.Cursor()
	fmt.Fprintf(myView, "Title: %s\n\n", a.listings.Children[cy].Data["title"])
	fmt.Fprintf(myView, "Author: %s\n\n", a.listings.Children[cy].Data["author"])
	fmt.Fprintf(myView, "URL: %s\n\n", a.listings.Children[cy].Data["url"])
	fmt.Fprintf(myView, "%s\n\n\n", a.listings.Children[cy].Data["selftext"])
	fmt.Fprintf(myView, "%s", a.listings.Children[cy].Data)
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

func populateView(g *gocui.Gui, v *gocui.View) error {
	log.Println("populating main view")
	v.Clear()
	green := color.New(color.FgGreen)
	for _, c := range a.listings.Children {
		green.Fprintf(v, "%s \t %s\n", c.Data["subreddit_name_prefixed"], c.Data["title"])
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
