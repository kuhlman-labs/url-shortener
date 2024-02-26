package urlshortener

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http url",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https url",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "localhost url",
			url:     "http://localhost",
			wantErr: true,
		},
		{
			name:    "127.0.0.1 url",
			url:     "http://127.0.0.1",
			wantErr: true,
		},
		{
			name:    "invalid url",
			url:     "not a url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateShortURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http url",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https url",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "localhost url",
			url:     "http://localhost",
			wantErr: true,
		},
		{
			name:    "127.0.0.1 url",
			url:     "http://127.0.0.1",
			wantErr: true,
		},
		{
			name:    "invalid url",
			url:     "not a url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := (&URL{}).GenerateShortURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateShortURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
