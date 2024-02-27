package main

import (
	"log"
	"net/http"
	"os"

	urlshortener "github.com/kuhlman-labs/url-shortener/url-shortener"
	"gopkg.in/yaml.v3"
)

type Config struct {
	TemplatePath string `yaml:"template_path"`
	Port         string `yaml:"port"`
	Domain       string `yaml:"domain"`
}

func main() {
	// Read the config.yaml file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config.yaml: %v", err)
	}

	// Create a new SQL URL repository
	db, err := urlshortener.NewSQLURLRepository()
	if err != nil {
		log.Fatalf("Error creating SQL URL repository: %v", err)
	}

	// Start the URL handler
	err = http.ListenAndServe(":"+config.Port, urlshortener.URLHandler(db, config.TemplatePath))
	if err != nil {
		log.Fatalf("Error starting URL handler: %v", err)
	}
}
