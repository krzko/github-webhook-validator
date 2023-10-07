package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

// Config holds configuration values
type Config struct {
	Secret string
}

// App encapsulates the application logic and dependencies
type App struct {
	config Config
	logger *slog.Logger
}

// NewApp creates a new App instance
func NewApp(config Config) *App {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return &App{config: config, logger: logger}
}

// ValidateSignature validates the GitHub webhook signature
func (app *App) ValidateSignature(body []byte, signatureHeader string) bool {
	if signatureHeader == "" || len(signatureHeader) < 7 {
		app.logger.Error("Invalid signature header")
		return false
	}
	app.logger.Debug("Validating signature...")

	computedHash := hmac.New(sha256.New, []byte(app.config.Secret))
	computedHash.Write(body)
	expectedSig := hex.EncodeToString(computedHash.Sum(nil))

	return hmac.Equal([]byte(expectedSig), []byte(signatureHeader[7:]))
}

// WebhookHandler handles GitHub webhook requests
func (app *App) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("Received webhook request")

	if r.Method != http.MethodPost {
		app.logger.Error("Method Not Allowed")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.logger.Error("Error reading request body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	signatureHeader := r.Header.Get("X-Hub-Signature-256")
	if app.ValidateSignature(body, signatureHeader) {
		app.logger.Info("Payload Validated")
		w.Write([]byte("Payload Validated\n"))
	} else {
		app.logger.Warn("Unauthorized - Signature Mismatch")
		http.Error(w, "Unauthorized - Signature Mismatch", http.StatusUnauthorized)
	}
}

func main() {
	config := Config{Secret: os.Getenv("WEBHOOK_SECRET")}
	app := NewApp(config)

	app.logger.Info("Starting webhook server...")
	http.HandleFunc("/webhook", app.WebhookHandler)
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	app.logger.Info("Server is listening on %s...", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
