package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (a *app) setKeybindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'r' {
				log.Println("detected 'r'")
				if a.postList.HasFocus() {
					a.listings = getPosts("https://oauth.reddit.com/hot", a.client)
					a.populatePosts()
				}
			}
			if event.Rune() == 'p' {
				log.Println("detected 'p'")
				if a.post.HasFocus() {
					a.commentsHandler()
				}
			}
		} else if event.Key() == tcell.KeyEsc {
			log.Println("detected 'esc'")
			if a.comments.HasFocus() {
				a.goToPost()
			}
		}
		return event
	})

}

func (a *app) postHandler(i int, title string, secondaryText string, shortcut rune) {
	log.Println("handling post", i)
	a.activePost = a.listings[i]

	a.post.SetTitle(fmt.Sprintf("%s", *a.activePost.Subreddit))
	a.post.SetText(fmt.Sprintf("%s\n%s - %d - %.f\n%s\n\n%s",
		*a.activePost.Title,
		*a.activePost.Author,
		*a.activePost.Score,
		*a.activePost.Created,
		*a.activePost.Url,
		*a.activePost.Selftext,
	))

	a.goToPost()
}

func (a *app) commentsHandler() {
	log.Println("going to comments")
	node := tview.NewTreeNode("my comment")
	subreddit := *a.activePost.Subreddit
	id := *a.activePost.Id
	_ = getPosts(fmt.Sprintf("https://oauth.reddit.com/%s/comments/%s", subreddit, id), a.client)
	for i := 0; i <= 5; i++ {
		n := tview.NewTreeNode(fmt.Sprintf("%s - this\n\nis\n\nmultiline", strconv.Itoa(i)))
		node.AddChild(n)
	}
	a.comments.SetRoot(node)
	a.goToComments()
}

func (a *app) goToPostList(k tcell.Key) {
	log.Println("going to postList")
	a.pages.SwitchToPage("postList")
}

func (a *app) goToPost() {
	log.Println("going to post")
	a.pages.SwitchToPage("post")
}

func (a *app) goToComments() {
	a.pages.SwitchToPage("comments")
}
