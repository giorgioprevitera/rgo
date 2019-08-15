package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/giorgioprevitera/oauth2"
	"github.com/pkg/browser"
)

type oauthSession struct {
	config *oauth2.Config
	state  string
	code   string
	server http.Server
	done   chan bool
}

func newOauthSession() *oauthSession {
	return &oauthSession{
		config: &oauth2.Config{
			ClientID:     "OskDlD8hMZeV5w",
			ClientSecret: "mr9UbAo0eyK5v4dRY2H-Q0wdS2A",
			Scopes:       []string{"identity", "read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.reddit.com/api/v1/authorize",
				TokenURL: "https://www.reddit.com/api/v1/access_token",
			},
			RedirectURL: "http://127.0.0.1:65010",
			Duration:    "permanent",
		},
		state: "4UfyGdH5W8DEM78F60mqZEavvvI",
	}
}

func saveToken(token *oauth2.Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Panicln("unable to decode token")
	}
	if ioutil.WriteFile(".token", data, 0600) != nil {
		log.Panicln("unable to write token to disk")

	}
	return nil
}

func getCachedToken() (*oauth2.Token, error) {
	var token *oauth2.Token
	bs, err := ioutil.ReadFile(".token")
	if err != nil {
		log.Println("no cached token")
		return token, errors.New("no cached token")
	}
	if json.Unmarshal(bs, &token) != nil {
		log.Panic("unable to decode cached token from disk")
	}
	log.Println("got cached token")
	return token, nil
}

func (o *oauthSession) getToken(ctx context.Context) *oauth2.Token {
	var token *oauth2.Token
	var err error

	token, err = getCachedToken()
	if err != nil {
		o.retrieveCode()
		token, err = o.config.Exchange(ctx, o.code)
		if err != nil {
			log.Fatal(err)
		}
		saveToken(token)
	}
	log.Println("getToken:", token)
	return token
}

func (o *oauthSession) callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != o.state {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	o.code = r.FormValue("code")
	if o.code == "" {
		http.Error(w, "unable to fetch token", 403)
	}
	fmt.Fprintln(w, "Succesfully authenticated, you can now close this page")
	go o.server.Shutdown(context.Background())
}

func (o *oauthSession) retrieveCode() {
	url := o.config.AuthCodeURL(o.state, oauth2.AccessTypeOffline)
	log.Println("retrieving code from", url)

	browser.OpenURL(url)

	m := http.NewServeMux()
	o.server = http.Server{Addr: ":65010", Handler: m}
	m.HandleFunc("/", o.callbackHandler)
	o.server.ListenAndServe()
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func getUrl(url string, client *http.Client) (*http.Response, error) {
	res, err := client.Get(url)
	return res, err
}
