package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/gocui"
)

var a app

type app struct {
	gui      *gocui.Gui
	client   *http.Client
	listings *listing
}

func main() {
	f, err := os.OpenFile("rgo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	a.client = getClient()
	a.getData()

	a.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer a.gui.Close()

	a.gui.SetManagerFunc(layout)
	setKeybindings(a.gui)

	if err := a.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
