package main

import(
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/studentlayer", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving %s request for %s", r.Method, r.URL)
		_, err := fmt.Fprintf(w, "Hello world")
		if err != nil {
			panic(err)
		}
	})

	log.Println("Running on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
