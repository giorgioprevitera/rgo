package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gdamore/tcell"
	"github.com/giorgioprevitera/redditproto"
	"github.com/rivo/tview"
)

type app struct {
	client         *http.Client
	listings       []*redditproto.Link
	postList       *tview.List
	app            *tview.Application
	pages          *tview.Pages
	post           *tview.TextView
	comments       *tview.TreeView
	info           *tview.TextView
	layout         *tview.Flex
	activePost     *redditproto.Link
	activeComments []*redditproto.Comment
}

func newApp() (*app, error) {
	a := &app{}

	a.init()
	a.createGui()

	a.listings = getPosts("https://oauth.reddit.com/rising", a.client)
	a.populatePosts()

	return a, nil
}

func (a *app) init() {
	a.client = getClient()
	a.app = tview.NewApplication()
	a.pages = tview.NewPages()
	a.postList = tview.NewList()
	a.post = tview.NewTextView().SetRegions(true)
	a.info = tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWrap(false)
	a.comments = tview.NewTreeView()

	a.setKeybindings()
}

func (a *app) createGui() {
	a.pages.SetBorder(true)
	a.post.SetBorder(true)
	a.comments.SetBorder(true)
	a.info.SetBorder(false)

	a.post.SetTitleColor(tcell.Color(10))
	a.post.SetWrap(true).SetWordWrap(true)

	a.pages.AddPage("postList", a.postList, true, true)
	a.pages.AddPage("post", a.post, true, false)
	a.pages.AddPage("comments", a.comments, true, false)

	a.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.pages, 0, 1, true).
		AddItem(a.info, 1, 1, false)

	a.postList.SetSelectedFunc(a.postHandler)

	a.post.SetDoneFunc(a.goToPostList)

	fmt.Fprintf(a.info, "%s - %s",
		"r: refresh list",
		"enter: open post",
	)
}

func (a *app) populatePosts() {
	a.postList.Clear()
	for _, c := range a.listings {
		title := fmt.Sprintf("%s", *c.Title)
		secondaryText := fmt.Sprintf("%s - %s - %d - %.f\n\n",
			*c.Subreddit,
			*c.Author,
			*c.Score,
			*c.Created,
		)
		a.postList.AddItem(title, secondaryText, 0, nil)
	}
}

func getPosts(url string, client *http.Client) []*redditproto.Link {
	responseBytes := getParsableBytes(url, client)
	links, _, _, _ := redditproto.ParseListing(responseBytes)
	log.Println("len links", len(links))
	return links
}

func getThread(url string, client *http.Client) (*redditproto.Link, error) {
	responseBytes := getParsableBytes(url, client)
	comments, err := redditproto.ParseThread(responseBytes)
	if err != nil {
		log.Println("getThread", err)
		return nil, errors.New(err.Error())
	}
	log.Println("comments", comments)
	return comments, nil
}

func getParsableBytes(url string, client *http.Client) []byte {
	var buf bytes.Buffer
	var err error

	useCachedResults := false
	responseBytes := buf.Bytes()

	log.Println("getting data")

	if useCachedResults {
		log.Println("using cached results")
		responseBytes, err = ioutil.ReadFile("dump.json")
		if err != nil {
			log.Panic(err)
		}
	} else {
		res, err := getUrl(url, client)
		if err != nil {
			log.Panic("unable to retrieve response", err)
		}
		defer res.Body.Close()

		buf.ReadFrom(res.Body)
		responseBytes = buf.Bytes()
	}
	log.Printf("responseBytes: %s", responseBytes)
	return responseBytes
}
