package urlshortener

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/url"
)

const domain string = "https://kuhlman-labs.io/"

type URL struct {
	LongURL  string
	ShortURL string
	Slug     string
}

func generateSlug(n int, u *URL) (*URL, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	u.Slug = base64.RawURLEncoding.EncodeToString(b)
	return u, nil
}

func (*URL) GenerateShortURL(longURL string) (*URL, error) {

	u := &URL{
		LongURL: longURL,
	}

	log.Printf("Generating short URL for: %s", u.LongURL)

	err := validateURL(u.LongURL)
	if err != nil {
		return nil, err
	}

	slug, err := generateSlug(6, u)
	if err != nil {
		return nil, err
	}

	u.ShortURL = domain + slug.Slug
	u.Slug = slug.Slug
	return u, nil
}

func validateURL(u string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}

	return nil
}
