// Package main provides ...
package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/mitchellh/mapstructure"
)

type thing struct {
	Kind string                 `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

type listing struct {
	Children []thing `json:"children,omitempty"`
}

func (a *app) getData() {
	log.Println("getting data")

	res, err := a.client.Get("https://oauth.reddit.com/hot")
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

	a.listings = &listing{}
	mapstructure.Decode(things.Data, a.listings)
}
