package main

import (
	"log"
	"os"
)

var a *app

func main() {

	// Logging
	f, err := os.OpenFile("rgo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("--------Start---------")

	a, err := newApp()

	if err := a.app.SetRoot(a.layout, true).Run(); err != nil {
		panic(err)
	}
}
