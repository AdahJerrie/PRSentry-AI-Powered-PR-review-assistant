package main

import (
	"log"
	"net/http"
	"os"

	"prsentry/go-service/internal/webhook"
)

func main() {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("GITHUB_WEBHOOK_SECRET is not set")
	}

	http.HandleFunc("/webhook", webhook.NewHandler(secret))

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}