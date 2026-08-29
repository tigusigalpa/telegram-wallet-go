package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	walletpay "github.com/tigusigalpa/telegram-wallet-go"
)

const webhookPayload = `[{"eventDateTime":"2023-07-25T16:47:06Z","eventId":1,"type":"ORDER_PAID","payload":{"id":1,"number":"ORDER-1","externalId":"external-1","orderAmount":{"amount":"1.00","currencyCode":"USD"}}}]`

func TestWalletPayWebhookHandler(t *testing.T) {
	const (
		apiKey    = "test-api-key"
		timestamp = "1234567890"
		path      = "/webhook/walletpay"
	)

	client := walletpay.NewClient(apiKey)
	called := false
	handler := WalletPayWebhookHandler(client, func(w http.ResponseWriter, r *http.Request, events []walletpay.WebhookEvent) {
		called = true
		assert.Len(t, events, 1)
		assert.Equal(t, int64(1), events[0].EventID)
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("accepts a valid webhook", func(t *testing.T) {
		called = false
		req := signedRequest(t, apiKey, timestamp, path, webhookPayload)
		recorder := httptest.NewRecorder()

		handler(recorder, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("rejects missing headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(webhookPayload))
		recorder := httptest.NewRecorder()

		handler(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Equal(t, "missing webhook headers\n", recorder.Body.String())
	})

	t.Run("rejects an invalid signature", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(webhookPayload))
		req.Header.Set("WalletPay-Timestamp", timestamp)
		req.Header.Set("WalletPay-Signature", "invalid")
		recorder := httptest.NewRecorder()

		handler(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Equal(t, "invalid signature\n", recorder.Body.String())
	})

	t.Run("rejects an invalid payload", func(t *testing.T) {
		const invalidPayload = "not json"
		req := signedRequest(t, apiKey, timestamp, path, invalidPayload)
		recorder := httptest.NewRecorder()

		handler(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "invalid payload\n", recorder.Body.String())
	})

	t.Run("rejects an unreadable body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.Body = io.NopCloser(errorReader{})
		req.Header.Set("WalletPay-Timestamp", timestamp)
		req.Header.Set("WalletPay-Signature", "unused")
		recorder := httptest.NewRecorder()

		handler(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "cannot read body\n", recorder.Body.String())
	})
}

func signedRequest(t *testing.T, apiKey, timestamp, path, body string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("WalletPay-Timestamp", timestamp)
	req.Header.Set("WalletPay-Signature", webhookSignature(apiKey, req.Method, path, timestamp, []byte(body)))
	return req
}

func webhookSignature(apiKey, method, path, timestamp string, body []byte) string {
	encodedBody := base64.StdEncoding.EncodeToString(body)
	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(method + "." + path + "." + timestamp + "." + encodedBody))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

var _ io.Reader = errorReader{}
