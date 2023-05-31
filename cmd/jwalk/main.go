package main

import (
	"fmt"
	"net/http"

	"github.com/willmeyers/jwalk"
)

func main() {
	c := jwalk.NewClient(
		&http.Client{},
		"https://api.fastmail.com/.well-known/jmap",
		"AUTH CODE",
	)

	fmt.Println(c.Session.EventSourceURL)

	events, err := jwalk.OpenSSEventConnection(c.Session.EventSourceURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for event := range events {
		fmt.Println(event)
	}
}
