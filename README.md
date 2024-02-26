# URL Shortener

This is a simple URL shortener written in Go. It provides an API to create short URLs from long URLs and redirects requests from the short URL to the original long URL.

## Features

- URL validation: Checks if the URL is valid and has the correct scheme (http or https).
- Slug generation: Generates a unique slug for each URL.
- Redirection: Redirects requests from the short URL to the original long URL.

## Usage

### Generate Short URL

To generate a short URL, call the `GenerateShortURL` function with the long URL as an argument. This function will validate the URL, generate a unique slug, and return a `URL` struct with the long URL, short URL, and slug.

```go
u, err := GenerateShortURL("http://example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Println(u.ShortURL) // Outputs: http://localhost:8080/abc123
```

### Redirection

To redirect requests from the short URL to the original long URL, use the `URLHandler` function. This function will look up the slug in the URL, find the corresponding long URL in the database, and redirect to the long URL.

```go
http.HandleFunc("/", URLHandler)
http.ListenAndServe(":8080", nil)
```

## Installation

To install this package, run the following command:

```bash
go get github.com/kuhlman-labs/urlshortener
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)