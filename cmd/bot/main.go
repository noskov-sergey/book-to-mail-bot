package main

import (
	"flag"
	"log"
)

const tgBotHost = "api.telegram.org"

func main() {
	// New Client tg

	// NewClientMail()

	// New() Fetcher - take updates from tg - send mail

	// new() Processor - make logic with files and control comands

	// Consumer - make logic (Fetcher, Processor)
}

func mustToken() string {
	t := flag.String("tg-host", "", "should be tg token for access bot")
	flag.Parse()

	if *t == "" {
		log.Fatal("can't find access token, shutdown")
	}

	return *t
}
