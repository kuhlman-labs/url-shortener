package main

import (
	"log"

	urlshortener "github.com/kuhlman-labs/url-shortener/url-shortener"
)

func main() {
	// Create a new SQL URL repository
	db, err := urlshortener.NewSQLURLRepository()
	if err != nil {
		log.Fatalf("Error creating SQL URL repository: %v", err)
	}

	// Start the URL handler
	urlshortener.URLHandler(db)
}
