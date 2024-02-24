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

func URLHandler(db URLRepository) {
	query := db

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path[1:]
		log.Printf("Looking up slug: %s", slug)

		// Look up the long URL in the database.
		query, err := query.ReadURLBySlug(slug)
		if err != nil {
			http.Error(w, "Error reading URL", http.StatusInternalServerError)
			return
		}

		if query == nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		// Redirect to the long URL.
		http.Redirect(w, r, query.LongUrl, http.StatusSeeOther)

	})

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving the form")
		// Serve the form.
		tmpl := template.Must(template.ParseFiles("templates/form.html"))
		// Execute the template, writing the result to the http.ResponseWriter.
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/app", http.StatusSeeOther)
			return
		}

		longURL := r.FormValue("url")
		url := &URL{
			LongURL: longURL,
		}

		// Generate the short URL.
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

		err = query.CreateURL(u)
		if err != nil {
			http.Error(w, "Error creating URL", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/result.html"))
		if err != nil {
			http.Error(w, "Error loading the template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, u.ShortUrl)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		var urlRequest URLRequest

		switch r.Method {
		case "GET":
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&urlRequest)
			if err != nil {
				http.Error(w, "Error decoding request body", http.StatusBadRequest)
				return
			}

			log.Printf("GET request received for: %s", urlRequest.URL)

			response, err := query.ReadURL(urlRequest.URL)
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

		case "POST":
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

			err = query.CreateURL(url)
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
		case "PUT":
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&urlRequest)
			if err != nil {
				http.Error(w, "Error decoding request body", http.StatusBadRequest)
				return
			}

			log.Printf("PUT request received for: %s", urlRequest.URL)

			response, err := query.ReadURL(urlRequest.URL)
			if err != nil {
				http.Error(w, "Error reading URL", http.StatusInternalServerError)
				return
			}

			err = query.UpdateURL(urlRequest.URL, urlRequest.NewURL)
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

		case "DELETE":
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&urlRequest)
			if err != nil {
				http.Error(w, "Error decoding request body", http.StatusBadRequest)
				return
			}

			log.Printf("DELETE request received for: %s", urlRequest.URL)

			response, err := query.ReadURL(urlRequest.URL)
			if err != nil {
				http.Error(w, "Error reading URL", http.StatusInternalServerError)
				return
			}

			if response == nil {
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}

			err = query.DeleteURL(response.LongUrl)
			if err != nil {
				http.Error(w, "Error deleting URL", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
