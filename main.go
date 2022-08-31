package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Service struct {
	Store map[string]UniversityCollection
}

// GetHealth is a Server method that provides a handler function for the /health route.
func (s *Service) GetHeath(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s request for %s", r.Method, r.URL)

	// Check that status method is GET
	if r.Method != http.MethodGet {
		// Set status code to 405 MethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)

		// Write error in response body
		resp := []byte("405 - Method not allowed")
		bytes, err := w.Write(resp)
		if err != nil || bytes != len(resp) {
			panic(err)
		}
	}

	// If alive, set status code to 200
	w.WriteHeader(http.StatusOK)

	// Write success response body
	resp := []byte("OK\n")
	bytes, err := w.Write(resp)
	if err != nil || bytes != len(resp) {
		panic(err)
	}
}

// PostEmail is a Server method that provides a handler function for the /emails route.
func (s *Service) PostEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s request for %s", r.Method, r.URL)

	// Check that the request method is POST
	if r.Method != http.MethodPost {
		// Set status code to 405 MethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)

		// Write error in response body
		resp := []byte("405 - Method not allowed")
		bytes, err := w.Write(resp)
		if err != nil {
			panic(err)
		}
		if bytes != len(resp) {
			panic(err)
		}
	}

	// Decode JSON from body into Email struct
	var email Email
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		// Set status code to 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		resp := []byte("500 - Internal Server Error")

		// Write error in response body
		bytes, err := w.Write(resp)
		if err != nil {
			panic(err)
		}
		if bytes != len(resp) {
			panic(err)
		}
	}

	// Check the store for the university, using two-value assignment for the university and boolean result
	emailParts := strings.Split(email.Email, "@")
	if len(emailParts) != 2 {
		// Set status code to 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		resp := []byte("500 - Internal Server Error")

		// Write error in response body
		bytes, err := w.Write(resp)
		if err != nil {
			panic(err)
		}
		if bytes != len(resp) {
			panic(err)
		}
	}

	universities, verified := s.Store[emailParts[1]]
	email.Universities = universities
	email.Verified = verified

	// Write the email struct in body response
	bytes, err := json.Marshal(email)
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}

// LoadUniversities is a Server method that loads a universities JSON file into Server.Store.
func (s *Service) LoadUniversities(filepath string) {
	// Open file.
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	// Defer closing the file, and check for error during close.
	defer func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Read bytes in file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Unmarshal the json into a UniversityCollection struct.
	var universities UniversityCollection
	err = json.Unmarshal(bytes, &universities)
	if err != nil {
		panic(err)
	}

	// For each university, add each domain to the key/value store (a map) where the key is the domain.
	for _, university := range universities {
		for _, domain := range university.Domains {
			existing, exists := s.Store[domain]

			// If the domain is already in the store, append the domain to the UniversityCollection.
			// Else create a new UniversityCollection with just the university to start with.
			if exists {
				s.Store[domain] = append(existing, university)
			} else {
				s.Store[domain] = UniversityCollection{university}
			}
		}
	}
}

func main() {
	router := http.NewServeMux()

	server := Service{
		Store: make(map[string]UniversityCollection),
	}

	server.LoadUniversities("data/universities_clean.json")

	router.HandleFunc("/emails", server.PostEmail)
	router.HandleFunc("/health", server.GetHeath)

	log.Println("Running on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
