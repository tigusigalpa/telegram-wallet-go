package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tigusigalpa/telegram-wallet-go"
)

// WalletPayWebhookHandler wraps an http.HandlerFunc with Wallet Pay webhook
// signature verification. On success, it parses the events and passes them
// to the provided handler function.
func WalletPayWebhookHandler(
	client *walletpay.Client,
	handler func(w http.ResponseWriter, r *http.Request, events []walletpay.WebhookEvent),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := r.Header.Get("WalletPay-Timestamp")
		signature := r.Header.Get("WalletPay-Signature")

		if timestamp == "" || signature == "" {
			http.Error(w, "missing webhook headers", http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		if err := client.VerifyWebhook(r.Method, r.URL.Path, timestamp, body, signature); err != nil {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var events []walletpay.WebhookEvent
		if err := json.Unmarshal(body, &events); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		handler(w, r, events)
	}
}
