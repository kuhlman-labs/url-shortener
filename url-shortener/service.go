package urlshortener

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
)

const domain string = "http://localhost:8080/"

type URL struct {
	LongURL  string
	ShortURL string
	Slug     string
}

func generateSlug(n int, u *URL) (*URL, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error generating random bytes: %w", err)
	}

	u.Slug = base64.RawURLEncoding.EncodeToString(b)[:n] // use only the first n characters
	return u, nil
}

func (*URL) GenerateShortURL(longURL string) (*URL, error) {
	u := &URL{
		LongURL: longURL,
	}

	log.Printf("Generating short URL for: %s", u.LongURL)

	err := validateURL(u.LongURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	slug, err := generateSlug(6, u)
	if err != nil {
		return nil, fmt.Errorf("error generating slug: %w", err)
	}

	u.ShortURL = domain + slug.Slug
	u.Slug = slug.Slug
	return u, nil
}

func validateURL(u string) error {
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must have http or https scheme")
	}

	// Prevent SSRF attacks by disallowing URLs that point to localhost or private network addresses
	if strings.HasPrefix(parsedURL.Hostname(), "localhost") || strings.HasPrefix(parsedURL.Hostname(), "127.0.0.1") {
		return errors.New("invalid url host")
	}

	return nil
}
