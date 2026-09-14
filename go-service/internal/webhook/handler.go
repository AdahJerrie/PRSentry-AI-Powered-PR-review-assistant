package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func verifySignature(secret string, payload []byte, signatureHeader string) bool {
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return false
	}
	expectedHex := strings.TrimPrefix(signatureHeader, "sha256=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	computedHex := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(computedHex), []byte(expectedHex))
}

// NewHandler returns an http.HandlerFunc configured with the given webhook secret.
func NewHandler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read body", http.StatusBadRequest)
			return
		}

		signature := r.Header.Get("X-Hub-Signature-256")
		if !verifySignature(secret, body, signature) {
			log.Println("Invalid signature — rejecting request")
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var payload GitHubWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Println("Failed to parse payload:", err)
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		if payload.Action != "opened" && payload.Action != "synchronize" && payload.Action != "reopened" {
			log.Printf("Ignoring action: %s\n", payload.Action)
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("PR #%d on %s needs review (action: %s)\n",
			payload.PullRequest.Number, payload.Repository.FullName, payload.Action)
		log.Printf("Diff URL: %s\n", payload.PullRequest.DiffURL)

		w.WriteHeader(http.StatusOK)
	}
}
