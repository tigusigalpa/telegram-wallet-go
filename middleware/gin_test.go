//go:build gin
// +build gin

package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	walletpay "github.com/tigusigalpa/telegram-wallet-go"
)

func TestGinWebhookMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const (
		apiKey    = "test-api-key"
		timestamp = "1234567890"
		path      = "/webhook/walletpay"
	)

	t.Run("passes verified events to the handler", func(t *testing.T) {
		router := gin.New()
		router.POST(path, GinWebhookMiddleware(walletpay.NewClient(apiKey)), func(c *gin.Context) {
			events := c.MustGet("walletpay_events").([]walletpay.WebhookEvent)
			assert.Len(t, events, 1)
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(webhookPayload))
		req.Header.Set("WalletPay-Timestamp", timestamp)
		req.Header.Set("WalletPay-Signature", webhookSignature(apiKey, req.Method, path, timestamp, []byte(webhookPayload)))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("rejects a signed but invalid payload", func(t *testing.T) {
		const payload = "not json"
		router := gin.New()
		router.POST(path, GinWebhookMiddleware(walletpay.NewClient(apiKey)), func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(payload))
		req.Header.Set("WalletPay-Timestamp", timestamp)
		req.Header.Set("WalletPay-Signature", webhookSignature(apiKey, req.Method, path, timestamp, []byte(payload)))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"invalid payload"}`, recorder.Body.String())
	})
}
