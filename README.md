# Go URL Shortener

This project is a URL shortener written in Go. It provides an API to create short URLs from long URLs and redirects requests from the short URL to the original long URL.

## Features

- **URL Validation**: Checks if the URL is valid and has the correct scheme (http or https).
- **Slug Generation**: Generates a unique slug for each URL.
- **Redirection**: Redirects requests from the short URL to the original long URL.
- **API**: Provides an API with CRUD operations to create a short url from a given long URL.
- **Database**: Uses a SQLite database to store the long URL and slug.
- **Web Interface**: Provides a simple web interface to create short URLs from long URLs.
- **Tests**: Includes unit tests for the service, handler, and database.

## Running Locally

To run this project locally, follow these steps:

1. Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/urlshortener.git
```

2. Navigate to the project directory:

```bash
cd urlshortener
```

3. Run the project:

```bash
go run main.go
```

The URL shortener will start on port 8080. You can access the APIs at `http://localhost:8080/api`. You can also access the web interface at `http://localhost:8080/app`.

## Usage

### Using the API

#### GET
```bash
curl -X GET "http://localhost:8080/api" -H "Content-Type: application/json" -d '{"url": "http://example.com"}'
```

#### POST
```bash
curl -X POST "http://localhost:8080/api" -H "Content-Type: application/json" -d '{"url": "http://example.com"}'
```

#### PUT
```bash
curl -X PUT "http://localhost:8080/api" -H "Content-Type: application/json" -d '{"url": "http://example.com", "new_url": "http://example2.com"}'
```    

#### DELETE
```bash    
curl -X DELETE "http://localhost:8080/api" -H "Content-Type: application/json" -d '{"url": "http://example.com"}'
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)