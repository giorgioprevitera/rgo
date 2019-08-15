package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gdamore/tcell"
	"github.com/giorgioprevitera/redditproto"
	"github.com/rivo/tview"
)

type app struct {
	client     *http.Client
	listings   []*redditproto.Link
	postList   *tview.List
	app        *tview.Application
	pages      *tview.Pages
	post       *tview.TextView
	comments   *tview.TreeView
	info       *tview.TextView
	layout     *tview.Flex
	activePost *redditproto.Link
}

func newApp() (*app, error) {
	a := &app{}

	a.init()
	a.createGui()

	a.listings = getPosts("https://oauth.reddit.com/hot", a.client)
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
	for i, c := range a.listings {
		title := fmt.Sprintf("%s", *c.Title)
		secondaryText := fmt.Sprintf("%s - %s - %d - %.f\n\n",
			*c.Subreddit,
			*c.Author,
			*c.Score,
			*c.Created,
		)
		a.postList.AddItem(title, secondaryText, rune(i+97), nil)
	}
}

func getPosts(url string, client *http.Client) []*redditproto.Link {
	responseBytes := getParsableBytes(url, client)
	link, _, _, _ := redditproto.ParseListing(responseBytes)
	log.Println("link", link)
	log.Println("len link", len(link))
	return link
}

func getParsableBytes(url string, client *http.Client) []byte {
	var buf bytes.Buffer

	useCachedResults := false
	responseBytes := buf.Bytes()

	log.Println("getting data")

	if useCachedResults {
		log.Println("using cached results")
		responseBytes, _ = ioutil.ReadFile("dump.json")
	} else {
		res, err := getUrl(url, client)
		if err != nil {
			log.Panic("unable to retrieve response", err)
		}
		defer res.Body.Close()

		buf.ReadFrom(res.Body)
		responseBytes = buf.Bytes()
	}
	return responseBytes
}
