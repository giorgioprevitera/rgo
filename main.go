package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type oauthSession struct {
	config *oauth2.Config
	state  string
	code   string
	server http.Server
	done   chan bool
}

func newOauthSession() *oauthSession {
	o := &oauthSession{
		config: &oauth2.Config{
			ClientID:     "OskDlD8hMZeV5w",
			ClientSecret: "mr9UbAo0eyK5v4dRY2H-Q0wdS2A",
			Scopes:       []string{"identity", "read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.reddit.com/api/v1/authorize",
				TokenURL: "https://www.reddit.com/api/v1/access_token",
			},
			RedirectURL: "http://127.0.0.1:65010",
		},
		state: "myrandomstatestring",
	}
	return o
}

func shutdown(w http.ResponseWriter, r *http.Request) {

}

func (o *oauthSession) callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != o.state {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	o.code = r.FormValue("code")
	o.server.Shutdown(context.Background())
}

func (o *oauthSession) retrieveCode() {
	url := o.config.AuthCodeURL(o.state, oauth2.AccessTypeOffline)
	log.Println("retrieving code from", url)

	m := http.NewServeMux()
	o.server = http.Server{Addr: ":65010", Handler: m}
	m.HandleFunc("/", o.callbackHandler)
	o.server.ListenAndServe()
}

type transport struct {
	http.RoundTripper
	useragent string
}

// Any request headers can be modified here.
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", t.useragent)
	return t.RoundTripper.RoundTrip(req)
}

func main() {
	o := newOauthSession()
	o.retrieveCode()
	c := &http.Client{}
	c.Transport = &transport{http.DefaultTransport, "goreddit cli client"}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)

	tok, err := o.config.Exchange(ctx, o.code)
	if err != nil {
		log.Fatal(err)
	}

	client := o.config.Client(ctx, tok)
	res, _ := client.Get("https://oauth.reddit.com/api/v1/me")
	io.Copy(os.Stdout, res.Body)
}
