package urlshortener

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockURLRepository is a mock type for urlshortener.URLRepository
type MockURLRepository struct {
	mock.Mock
}

// ReadURLBySlug is a mock method for URLRepository.ReadURLBySlug
func (m *MockURLRepository) ReadURLBySlug(slug string) (*URLSchema, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*URLSchema), args.Error(1)
}

// CreateURL is a mock method for URLRepository.CreateURL
func (m *MockURLRepository) CreateURL(url *URLSchema) error {
	args := m.Called(url)
	return args.Error(0)
}

// ReadURL is a mock method for URLRepository.ReadURL
func (m *MockURLRepository) ReadURL(longUrl string) (*URLSchema, error) {
	args := m.Called(longUrl)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*URLSchema), args.Error(1)
}

// UpdateURL is a mock method for URLRepository.UpdateURL
func (m *MockURLRepository) UpdateURL(slug, newLongURL string) error {
	args := m.Called(slug, newLongURL)
	return args.Error(0)
}

// DeleteURL is a mock method for URLRepository.DeleteURL
func (m *MockURLRepository) DeleteURL(slug string) error {
	args := m.Called(slug)
	return args.Error(0)
}

func TestRootHandler(t *testing.T) {
	// Create a new mock URL repository
	repo := new(MockURLRepository)

	// Create a new URL handler with the mock URL repository
	handler := rootHandler(repo)

	// Expect a call to ReadURLBySlug with "abc123" and return a URLSchema
	repo.On("ReadURLBySlug", "abc123").Return(&URLSchema{
		Slug:     "abc123",
		LongUrl:  "http://example.com",
		ShortUrl: "http://localhost:8080/abc123",
	}, nil)

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/abc123", nil)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusSeeOther, rr.Code)

	// Check the redirect location
	assert.Equal(t, "http://example.com", rr.Header().Get("Location"))

	// Assert that the expectations were met
	repo.AssertExpectations(t)

}

func TestAppHandler(t *testing.T) {

	// Create a new URL handler with the mock URL repository
	handler := appHandler("../templates/")

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/app", nil)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestShortenHandler(t *testing.T) {
	// Create a new mock URL repository
	repo := new(MockURLRepository)

	// Set up the expectation
	repo.On("CreateURL", mock.Anything).Return(nil)

	// Create a new HTTP request with form data
	form := url.Values{}
	form.Add("url", "http://example.com")
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Create the handler
	handler := shortenHandler(repo, "../templates/")

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the expectations were met
	repo.AssertExpectations(t)

}

func Test_Api_Get(t *testing.T) {
	// Create a new mock URL repository
	repo := new(MockURLRepository)

	// Set up the expectation
	expectedResponse := &URLSchema{LongUrl: "http://example.com", ShortUrl: "http://short.com"}
	repo.On("ReadURL", "http://example.com").Return(expectedResponse, nil)

	// Create a new URLRequest
	urlRequest := URLRequest{
		URL: "http://example.com",
	}

	// Marshal the URLRequest to JSON
	jsonRequest, err := json.Marshal(urlRequest)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling the URLRequest", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/api", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handleGet function
	handleGet(rr, req, repo, urlRequest)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response
	var response URLSchema
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when unmarshaling the response", err)
	}

	// Check the response
	assert.Equal(t, expectedResponse, &response)

	// Assert that the expectations were met
	repo.AssertExpectations(t)

}

func Test_Api_Post(t *testing.T) {
	// Create a new mock URL repository
	repo := new(MockURLRepository)

	// Set up the expectation
	repo.On("CreateURL", mock.Anything).Return(nil)

	// Create a new URLRequest
	urlRequest := URLRequest{
		URL: "http://example.com",
	}

	// Marshal the URLRequest to JSON
	jsonRequest, err := json.Marshal(urlRequest)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling the URLRequest", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest("POST", "/api", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handlePost function
	handlePost(rr, req, repo, urlRequest)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Unmarshal the response
	var response URL
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when unmarshaling the response", err)
	}

	// Check the response
	assert.Equal(t, urlRequest.URL, response.LongURL)

	// Assert that the expectations were met
	repo.AssertExpectations(t)

}

func Test_Api_Put(t *testing.T) {
	// Create a new mock URL repository
	repo := new(MockURLRepository)

	// Set up the expectation
	expectedResponse := &URLSchema{LongUrl: "http://example.com", ShortUrl: "http://short.com"}
	repo.On("ReadURL", "http://example.com").Return(expectedResponse, nil)
	repo.On("UpdateURL", "http://example.com", "http://newexample.com").Return(nil)

	// Create a new URLRequest
	urlRequest := URLRequest{
		URL:    "http://example.com",
		NewURL: "http://newexample.com",
	}

	// Marshal the URLRequest to JSON
	jsonRequest, err := json.Marshal(urlRequest)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling the URLRequest", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest("PUT", "/api", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handlePut function
	handlePut(rr, req, repo, urlRequest)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response
	var response URLSchema
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when unmarshaling the response", err)
	}

	// Check the response
	assert.Equal(t, expectedResponse, &response)

	// Assert that the expectations were met
	repo.AssertExpectations(t)

}

func Test_Api_Delete(t *testing.T) {
	// Create a new mock URL repository
	repo := new(MockURLRepository)

	// Set up the expectation
	expectedResponse := &URLSchema{LongUrl: "http://example.com", ShortUrl: "http://short.com"}
	repo.On("ReadURL", "http://example.com").Return(expectedResponse, nil)
	repo.On("DeleteURL", "http://example.com").Return(nil)

	// Create a new URLRequest
	urlRequest := URLRequest{
		URL: "http://example.com",
	}

	// Marshal the URLRequest to JSON
	jsonRequest, err := json.Marshal(urlRequest)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling the URLRequest", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest("DELETE", "/api", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handleDelete function
	handleDelete(rr, req, repo, urlRequest)

	// Check the status code
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Assert that the expectations were met
	repo.AssertExpectations(t)

}
