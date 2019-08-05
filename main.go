package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/mitchellh/mapstructure"
)

type app struct {
	gui    *gocui.Gui
	client *http.Client
}

func getData() *listing {
	log.Println("getting data")
	client := getClient()

	res, err := client.Get("https://oauth.reddit.com/hot")
	if err != nil {
		log.Panic("unable to retrieve response", err)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)

	var things thing

	// res, _ := ioutil.ReadFile("dump.json")
	// json.Unmarshal(res, &things)

	json.Unmarshal(buf.Bytes(), &things)

	l := &listing{}
	mapstructure.Decode(things.Data, l)
	return l
}

func main() {
	f, err := os.OpenFile("rgo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	var a app

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

func populateView(g *gocui.Gui, v *gocui.View) error {
	log.Println("populating main view")
	v.Clear()
	green := color.New(color.FgGreen)
	l := getData()
	for _, c := range l.Children {
		green.Fprintf(v, "%s \t %s\n", c.Data["subreddit_name_prefixed"], c.Data["title"])
	}
	return nil
}

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
