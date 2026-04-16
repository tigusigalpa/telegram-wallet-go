package main

import (
	"log"
	"net/http"
	"os"

	walletpay "github.com/tigusigalpa/telegram-wallet-go"
	"github.com/tigusigalpa/telegram-wallet-go/middleware"
)

func main() {
	apiKey := os.Getenv("WALLETPAY_API_KEY")
	if apiKey == "" {
		log.Fatal("WALLETPAY_API_KEY environment variable is required")
	}

	client := walletpay.NewClient(apiKey)

	// Create webhook handler with signature verification
	http.HandleFunc("/webhook/walletpay", middleware.WalletPayWebhookHandler(
		client,
		func(w http.ResponseWriter, r *http.Request, events []walletpay.WebhookEvent) {
			log.Printf("Received %d webhook event(s)", len(events))

			for _, event := range events {
				log.Printf("Event ID: %d, Type: %s", event.EventID, event.Type)

				switch event.Type {
				case walletpay.WebhookEventOrderPaid:
					handleOrderPaid(event)
				case walletpay.WebhookEventOrderFailed:
					handleOrderFailed(event)
				default:
					log.Printf("Unknown event type: %s", event.Type)
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting webhook server on port %s...", port)
	log.Printf("Webhook endpoint: http://localhost:%s/webhook/walletpay", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleOrderPaid(event walletpay.WebhookEvent) {
	payload := event.Payload
	log.Printf("✅ Order PAID: ID=%d, Number=%s, ExternalID=%s",
		payload.ID, payload.Number, payload.ExternalID)

	if payload.SelectedPaymentOption != nil {
		log.Printf("   Payment: %s %s (fee: %s %s, net: %s %s)",
			payload.SelectedPaymentOption.Amount.Amount,
			payload.SelectedPaymentOption.Amount.CurrencyCode,
			payload.SelectedPaymentOption.AmountFee.Amount,
			payload.SelectedPaymentOption.AmountFee.CurrencyCode,
			payload.SelectedPaymentOption.AmountNet.Amount,
			payload.SelectedPaymentOption.AmountNet.CurrencyCode,
		)
		log.Printf("   Exchange rate: %s", payload.SelectedPaymentOption.ExchangeRate)
	}

	if payload.CustomData != "" {
		log.Printf("   Custom data: %s", payload.CustomData)
	}

	// TODO: Update your database, grant access, send notification, etc.
	// Example:
	// - Update payment status in database
	// - Grant premium access to user
	// - Send confirmation email/notification
	// - Trigger fulfillment process
}

func handleOrderFailed(event walletpay.WebhookEvent) {
	payload := event.Payload
	status := "UNKNOWN"
	if payload.Status != nil {
		status = string(*payload.Status)
	}

	log.Printf("❌ Order FAILED: ID=%d, Number=%s, ExternalID=%s, Status=%s",
		payload.ID, payload.Number, payload.ExternalID, status)

	if payload.CustomData != "" {
		log.Printf("   Custom data: %s", payload.CustomData)
	}

	// TODO: Handle failed payment
	// Example:
	// - Update payment status in database
	// - Send notification to user
	// - Log for analytics
}
