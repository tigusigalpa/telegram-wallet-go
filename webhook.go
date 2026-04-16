package walletpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// VerifyWebhook verifies the HMAC-SHA256 signature of an incoming webhook.
// httpMethod should be "POST", uriPath is the request path (e.g., "/webhook/"),
// timestamp comes from the WalletPay-Timestamp header,
// body is the raw request body bytes,
// signature comes from the WalletPay-Signature header.
func (c *Client) VerifyWebhook(httpMethod, uriPath, timestamp string, body []byte, signature string) error {
	base64Body := base64.StdEncoding.EncodeToString(body)
	stringToSign := httpMethod + "." + uriPath + "." + timestamp + "." + base64Body
	mac := hmac.New(sha256.New, []byte(c.apiKey))
	mac.Write([]byte(stringToSign))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrInvalidSignature
	}

	return nil
}

// ParseWebhookEvents parses the webhook payload into a slice of WebhookEvent.
func ParseWebhookEvents(body []byte) ([]WebhookEvent, error) {
	var events []WebhookEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to parse webhook events: %w", err)
	}
	return events, nil
}

// VerifyAndParseWebhook verifies the webhook signature and parses the events.
func (c *Client) VerifyAndParseWebhook(httpMethod, uriPath, timestamp string, body []byte, signature string) ([]WebhookEvent, error) {
	if err := c.VerifyWebhook(httpMethod, uriPath, timestamp, body, signature); err != nil {
		return nil, err
	}
	return ParseWebhookEvents(body)
}
