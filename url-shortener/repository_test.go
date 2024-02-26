package urlshortener

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateURL(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Auto-migrate the schema
	db.AutoMigrate(&URLSchema{})

	// Create the repository
	repo := &SQLURLRepository{
		db: db,
	}

	// Create the URL schema
	url := &URLSchema{
		Slug:     "testslug",
		ShortUrl: "http://short.com",
		LongUrl:  "http://long.com",
	}

	// Call CreateURL
	err = repo.CreateURL(url)
	assert.NoError(t, err)

	// Retrieve the URL from the database
	var retrievedURL URLSchema
	db.Where("slug = ?", "testslug").First(&retrievedURL)

	// Check that the retrieved URL matches the created URL
	assert.Equal(t, url.Slug, retrievedURL.Slug)
	assert.Equal(t, url.ShortUrl, retrievedURL.ShortUrl)
	assert.Equal(t, url.LongUrl, retrievedURL.LongUrl)
}

func TestReadURL(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	defer db.Close()

	// Create the repository
	repo := &SQLURLRepository{
		db: db,
	}

	// Auto-migrate the schema
	db.AutoMigrate(&URLSchema{})

	// Create the URL schema
	url := &URLSchema{
		Slug:     "abc123",
		ShortUrl: "http://localhost:8080/abc123",
		LongUrl:  "http://example.com",
	}

	// Call CreateURL
	err = repo.CreateURL(url)
	assert.NoError(t, err)

	// Call ReadURL
	url, err = repo.ReadURL("http://example.com")
	assert.NoError(t, err)

	// Check that the URL matches the expected URL
	assert.Equal(t, "abc123", url.Slug)
	assert.Equal(t, "http://example.com", url.LongUrl)
	assert.Equal(t, "http://localhost:8080/abc123", url.ShortUrl)
}

func TestReadURLBySlug(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create the repository
	repo := &SQLURLRepository{
		db: db,
	}

	// Auto-migrate the schema
	db.AutoMigrate(&URLSchema{})

	// Create the URL schema
	url := &URLSchema{
		Slug:     "abc123",
		ShortUrl: "http://localhost:8080/abc123",
		LongUrl:  "http://example.com",
	}

	// Call CreateURL
	err = repo.CreateURL(url)
	assert.NoError(t, err)

	// Call ReadURLBySlug
	url, err = repo.ReadURLBySlug("abc123")
	assert.NoError(t, err)

	// Check that the URL matches the expected URL
	assert.Equal(t, "abc123", url.Slug)
	assert.Equal(t, "http://example.com", url.LongUrl)
	assert.Equal(t, "http://localhost:8080/abc123", url.ShortUrl)
}

func TestUpdateURL(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create the repository
	repo := &SQLURLRepository{
		db: db,
	}

	// Auto-migrate the schema
	db.AutoMigrate(&URLSchema{})

	// Create the URL schema
	url := &URLSchema{
		Slug:     "abc123",
		ShortUrl: "http://localhost:8080/abc123",
		LongUrl:  "http://example.com",
	}

	// Call CreateURL
	err = repo.CreateURL(url)
	assert.NoError(t, err)

	// Call UpdateURL
	err = repo.UpdateURL("http://example.com", "http://example.org")
	assert.NoError(t, err)

	// Retrieve the URL from the database
	var retrievedURL URLSchema
	db.Where("slug = ?", "abc123").First(&retrievedURL)

	// Check that the retrieved URL matches the updated URL
	assert.Equal(t, "abc123", retrievedURL.Slug)
	assert.Equal(t, "http://example.org", retrievedURL.LongUrl)
	assert.Equal(t, "http://localhost:8080/abc123", retrievedURL.ShortUrl)
}

func TestDeleteURL(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create the repository
	repo := &SQLURLRepository{
		db: db,
	}

	// Auto-migrate the schema
	db.AutoMigrate(&URLSchema{})

	// Create the URL schema
	url := &URLSchema{
		Slug:     "abc123",
		ShortUrl: "http://localhost:8080/abc123",
		LongUrl:  "http://example.com",
	}

	// Call CreateURL
	err = repo.CreateURL(url)
	assert.NoError(t, err)

	// Call DeleteURL
	err = repo.DeleteURL("http://example.com")
	assert.NoError(t, err)

	// Retrieve the URL from the database
	var retrievedURL URLSchema
	db.Where("slug = ?", "abc123").First(&retrievedURL)

	// Check that the retrieved URL is nil
	assert.Equal(t, URLSchema{}, retrievedURL)
}
