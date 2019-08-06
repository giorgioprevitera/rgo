package main

import (
	"context"
	"net/http"

	"github.com/giorgioprevitera/oauth2"
)

type transport struct {
	http.RoundTripper
	useragent string
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", t.useragent)
	return t.RoundTripper.RoundTrip(req)
}

// getContext injects a custom User Agent, required to authenticate against the Reddit APIs
func getContext() context.Context {
	c := &http.Client{}
	c.Transport = &transport{http.DefaultTransport, "rgo is an Open Source CLI client for Reddit git@github.com:giorgioprevitera/rgo.git"}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)
	return ctx
}

func getClient() *http.Client {
	o := newOauthSession()
	ctx := getContext()
	token := o.getToken(ctx)
	client := o.config.Client(ctx, token)
	return client
}
