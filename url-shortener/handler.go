package urlshortener

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type URLRequest struct {
	URL    string `json:"url"`
	NewURL string `json:"new_url"`
}

func URLHandler(db URLRepository, templatePath string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler(db))
	mux.HandleFunc("/app", appHandler(templatePath))
	mux.HandleFunc("/shorten", shortenHandler(db, templatePath))
	mux.HandleFunc("/api", apiHandler(db))

	return mux
}

func rootHandler(db URLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path[1:]
		log.Printf("Looking up slug: %s", slug)

		query, err := db.ReadURLBySlug(slug)
		if err != nil {
			http.Error(w, "Error reading URL", http.StatusInternalServerError)
			return
		}

		if query == nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, query.LongUrl, http.StatusSeeOther)
	}
}

func appHandler(templatePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(templatePath + "form.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}

func shortenHandler(db URLRepository, templatePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/app", http.StatusSeeOther)
			return
		}

		longURL := r.FormValue("url")
		url := &URL{
			LongURL: longURL,
		}

		shortURL, err := url.GenerateShortURL(longURL)
		if err != nil {
			http.Error(w, "Error generating short URL", http.StatusInternalServerError)
			return
		}

		u := &URLSchema{
			LongUrl:  longURL,
			ShortUrl: shortURL.ShortURL,
			Slug:     shortURL.Slug,
		}

		err = db.CreateURL(u)
		if err != nil {
			http.Error(w, "Error creating URL", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles(templatePath + "result.html"))
		if err != nil {
			http.Error(w, "Error loading the template", http.StatusInternalServerError)
			return
		}

		log.Printf("Shortened URL: %s", u.ShortUrl)
		err = tmpl.Execute(w, u.ShortUrl)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}

func apiHandler(db URLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var urlRequest URLRequest
		err := json.NewDecoder(r.Body).Decode(&urlRequest)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGet(w, r, db, urlRequest)
		case http.MethodPost:
			handlePost(w, r, db, urlRequest)
		case http.MethodPut:
			handlePut(w, r, db, urlRequest)
		case http.MethodDelete:
			handleDelete(w, r, db, urlRequest)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, db URLRepository, urlRequest URLRequest) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&urlRequest)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	log.Printf("GET request received for: %s", urlRequest.URL)

	response, err := db.ReadURL(urlRequest.URL)
	if err != nil {
		http.Error(w, "Error reading URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, db URLRepository, urlRequest URLRequest) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&urlRequest)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	log.Printf("POST request received for: %s", urlRequest.URL)

	u := &URL{}
	u, err = u.GenerateShortURL(urlRequest.URL)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	url := &URLSchema{
		Slug:     u.Slug,
		ShortUrl: u.ShortURL,
		LongUrl:  u.LongURL,
	}

	err = db.CreateURL(url)
	if err != nil {
		http.Error(w, "Error creating URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handlePut(w http.ResponseWriter, r *http.Request, db URLRepository, urlRequest URLRequest) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&urlRequest)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	log.Printf("PUT request received for: %s", urlRequest.URL)

	response, err := db.ReadURL(urlRequest.URL)
	if err != nil {
		http.Error(w, "Error reading URL", http.StatusInternalServerError)
		return
	}

	err = db.UpdateURL(urlRequest.URL, urlRequest.NewURL)
	if err != nil {
		http.Error(w, "Error updating URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request, db URLRepository, urlRequest URLRequest) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&urlRequest)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	log.Printf("DELETE request received for: %s", urlRequest.URL)

	response, err := db.ReadURL(urlRequest.URL)
	if err != nil {
		http.Error(w, "Error reading URL", http.StatusInternalServerError)
		return
	}

	if response == nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	err = db.DeleteURL(response.LongUrl)
	if err != nil {
		http.Error(w, "Error deleting URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
