package main

import (
	"log"
	"net/http"

	urlshortener "github.com/kuhlman-labs/url-shortener/url-shortener"
)

func main() {
	// Create a new SQL URL repository
	db, err := urlshortener.NewSQLURLRepository()
	if err != nil {
		log.Fatalf("Error creating SQL URL repository: %v", err)
	}

	// Start the URL handler
	err = http.ListenAndServe(":8080", urlshortener.URLHandler(db, "templates/"))
	if err != nil {
		log.Fatalf("Error starting URL handler: %v", err)
	}
}
