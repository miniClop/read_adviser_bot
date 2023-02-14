package main

import (
	"flag"
	"log"

	tgClient "read_adviser_bot/clients/telegram"
	"read_adviser_bot/consumer/event-consumer"
	"read_adviser_bot/events/telegram"
	"read_adviser_bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Printf("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

	//fetcher = fetcher.New()

	//processor = processor.New()

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
