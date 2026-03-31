package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := map[string]string{
			"status":  "ok",
			"service": "go-api-service",
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

	log.Printf("API started on :%s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}