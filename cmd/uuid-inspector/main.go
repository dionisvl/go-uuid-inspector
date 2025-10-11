package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dionisvl/go-uuid-inspector/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/inspect", handlers.Inspect)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("UUID Inspector running on http://localhost:%s\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
