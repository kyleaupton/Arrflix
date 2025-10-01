package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Basic health route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	// Dummy route to show it’s alive
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Snaggle dummy API is running on %s\n", port)
	})

	// Background ticker just to prove logs are flowing
	go func() {
		for {
			log.Printf("[snaggle] still alive at %s", time.Now().Format(time.RFC3339))
			time.Sleep(10 * time.Second)
		}
	}()

	log.Printf("Starting dummy API on :%s ...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
