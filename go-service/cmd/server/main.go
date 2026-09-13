package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func verifySignature(secret string, payload []byte, signatureHeader string) bool {
	// GitHub sends the header as "sha256=<hex digest>"
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return false
	}
	expectedHex := strings.TrimPrefix(signatureHeader, "sha256=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	computedHex := hex.EncodeToString(mac.Sum(nil))

	// hmac.Equal is a constant-time comparison — prevents timing attacks
	return hmac.Equal([]byte(computedHex), []byte(expectedHex))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read body", http.StatusBadRequest)
		return
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	signature := r.Header.Get("X-Hub-Signature-256")

	if !verifySignature(secret, body, signature) {
		log.Println("Invalid signature — rejecting request")
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	log.Println("Received webhook event:")
	log.Println(string(body))

	w.WriteHeader(http.StatusOK)
}

func main() {
	if os.Getenv("GITHUB_WEBHOOK_SECRET") == "" {
		log.Fatal("GITHUB_WEBHOOK_SECRET is not set")
	}

	http.HandleFunc("/webhook", webhookHandler)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
