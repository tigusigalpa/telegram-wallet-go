//go:build echo
// +build echo

package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	walletpay "github.com/tigusigalpa/telegram-wallet-go"
)

func TestEchoWebhookMiddleware(t *testing.T) {
	const (
		apiKey    = "test-api-key"
		timestamp = "1234567890"
		path      = "/webhook/walletpay"
	)

	t.Run("passes verified events to the handler", func(t *testing.T) {
		e := echo.New()
		e.POST(path, func(c echo.Context) error {
			events := c.Get("walletpay_events").([]walletpay.WebhookEvent)
			assert.Len(t, events, 1)
			return c.NoContent(http.StatusNoContent)
		}, EchoWebhookMiddleware(walletpay.NewClient(apiKey)))

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(webhookPayload))
		req.Header.Set("WalletPay-Timestamp", timestamp)
		req.Header.Set("WalletPay-Signature", webhookSignature(apiKey, req.Method, path, timestamp, []byte(webhookPayload)))
		recorder := httptest.NewRecorder()

		e.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("rejects a signed but invalid payload", func(t *testing.T) {
		const payload = "not json"
		e := echo.New()
		e.POST(path, func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		}, EchoWebhookMiddleware(walletpay.NewClient(apiKey)))

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(payload))
		req.Header.Set("WalletPay-Timestamp", timestamp)
		req.Header.Set("WalletPay-Signature", webhookSignature(apiKey, req.Method, path, timestamp, []byte(payload)))
		recorder := httptest.NewRecorder()

		e.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"invalid payload"}`, recorder.Body.String())
	})
}
