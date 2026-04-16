package walletpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyWebhook(t *testing.T) {
	apiKey := "test-api-key-12345"
	client := NewClient(apiKey)

	tests := []struct {
		name       string
		httpMethod string
		uriPath    string
		timestamp  string
		body       string
		signature  string
		wantErr    bool
	}{
		{
			name:       "valid signature",
			httpMethod: "POST",
			uriPath:    "/webhook/walletpay",
			timestamp:  "1234567890123456789",
			body:       `{"test":"data"}`,
			signature:  computeSignature(apiKey, "POST", "/webhook/walletpay", "1234567890123456789", `{"test":"data"}`),
			wantErr:    false,
		},
		{
			name:       "valid signature with trailing slash",
			httpMethod: "POST",
			uriPath:    "/webhook/walletpay/",
			timestamp:  "1234567890123456789",
			body:       `{"test":"data"}`,
			signature:  computeSignature(apiKey, "POST", "/webhook/walletpay/", "1234567890123456789", `{"test":"data"}`),
			wantErr:    false,
		},
		{
			name:       "invalid signature",
			httpMethod: "POST",
			uriPath:    "/webhook/walletpay",
			timestamp:  "1234567890123456789",
			body:       `{"test":"data"}`,
			signature:  "invalid-signature",
			wantErr:    true,
		},
		{
			name:       "tampered body",
			httpMethod: "POST",
			uriPath:    "/webhook/walletpay",
			timestamp:  "1234567890123456789",
			body:       `{"test":"tampered"}`,
			signature:  computeSignature(apiKey, "POST", "/webhook/walletpay", "1234567890123456789", `{"test":"data"}`),
			wantErr:    true,
		},
		{
			name:       "wrong timestamp",
			httpMethod: "POST",
			uriPath:    "/webhook/walletpay",
			timestamp:  "9999999999999999999",
			body:       `{"test":"data"}`,
			signature:  computeSignature(apiKey, "POST", "/webhook/walletpay", "1234567890123456789", `{"test":"data"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.VerifyWebhook(tt.httpMethod, tt.uriPath, tt.timestamp, []byte(tt.body), tt.signature)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidSignature, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseWebhookEvents(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantLen int
		wantErr bool
	}{
		{
			name: "valid ORDER_PAID event",
			body: `[{
				"eventDateTime": "2023-07-25T16:47:06.383352Z",
				"eventId": 9906750163000780,
				"type": "ORDER_PAID",
				"payload": {
					"id": 10030467668508673,
					"number": "XYTNJP2O",
					"customData": "client_ref=4E89",
					"externalId": "JDF23NN",
					"orderAmount": {
						"amount": "0.100000340",
						"currencyCode": "TON"
					},
					"selectedPaymentOption": {
						"amount": {"amount": "0.132653", "currencyCode": "USDT"},
						"amountFee": {"amount": "0.001327", "currencyCode": "USDT"},
						"amountNet": {"amount": "0.131326", "currencyCode": "USDT"},
						"exchangeRate": "1.3265247467314987"
					},
					"orderCompletedDateTime": "2023-07-28T10:20:17.628946Z"
				}
			}]`,
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "valid ORDER_FAILED event",
			body: `[{
				"eventDateTime": "2023-07-25T16:47:06.383352Z",
				"eventId": 9906750163000781,
				"type": "ORDER_FAILED",
				"payload": {
					"id": 10030467668508674,
					"number": "XYTNJP2P",
					"externalId": "JDF23NO",
					"orderAmount": {
						"amount": "0.100000340",
						"currencyCode": "TON"
					},
					"status": "EXPIRED"
				}
			}]`,
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "multiple events",
			body: `[
				{
					"eventDateTime": "2023-07-25T16:47:06.383352Z",
					"eventId": 1,
					"type": "ORDER_PAID",
					"payload": {
						"id": 1,
						"number": "A",
						"externalId": "E1",
						"orderAmount": {"amount": "1.00", "currencyCode": "USD"}
					}
				},
				{
					"eventDateTime": "2023-07-25T16:47:07.383352Z",
					"eventId": 2,
					"type": "ORDER_FAILED",
					"payload": {
						"id": 2,
						"number": "B",
						"externalId": "E2",
						"orderAmount": {"amount": "2.00", "currencyCode": "EUR"}
					}
				}
			]`,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "invalid json",
			body:    `{invalid json}`,
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := ParseWebhookEvents([]byte(tt.body))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, events, tt.wantLen)
				if tt.wantLen > 0 {
					assert.NotZero(t, events[0].EventID)
					assert.NotEmpty(t, events[0].Type)
					assert.NotZero(t, events[0].Payload.ID)
				}
			}
		})
	}
}

func TestVerifyAndParseWebhook(t *testing.T) {
	apiKey := "test-api-key-12345"
	client := NewClient(apiKey)

	httpMethod := "POST"
	uriPath := "/webhook/walletpay"
	timestamp := "1234567890123456789"
	body := `[{
		"eventDateTime": "2023-07-25T16:47:06.383352Z",
		"eventId": 9906750163000780,
		"type": "ORDER_PAID",
		"payload": {
			"id": 10030467668508673,
			"number": "XYTNJP2O",
			"externalId": "JDF23NN",
			"orderAmount": {
				"amount": "0.100000340",
				"currencyCode": "TON"
			}
		}
	}]`

	signature := computeSignature(apiKey, httpMethod, uriPath, timestamp, body)

	events, err := client.VerifyAndParseWebhook(httpMethod, uriPath, timestamp, []byte(body), signature)
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, int64(9906750163000780), events[0].EventID)
	assert.Equal(t, WebhookEventOrderPaid, events[0].Type)
	assert.Equal(t, int64(10030467668508673), events[0].Payload.ID)
}

func TestVerifyAndParseWebhook_InvalidSignature(t *testing.T) {
	apiKey := "test-api-key-12345"
	client := NewClient(apiKey)

	httpMethod := "POST"
	uriPath := "/webhook/walletpay"
	timestamp := "1234567890123456789"
	body := `[{"eventDateTime": "2023-07-25T16:47:06.383352Z", "eventId": 1, "type": "ORDER_PAID", "payload": {"id": 1, "number": "A", "externalId": "E1", "orderAmount": {"amount": "1.00", "currencyCode": "USD"}}}]`
	signature := "invalid-signature"

	events, err := client.VerifyAndParseWebhook(httpMethod, uriPath, timestamp, []byte(body), signature)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidSignature, err)
	assert.Nil(t, events)
}

// Helper function to compute HMAC-SHA256 signature
func computeSignature(apiKey, httpMethod, uriPath, timestamp, body string) string {
	base64Body := base64.StdEncoding.EncodeToString([]byte(body))
	stringToSign := httpMethod + "." + uriPath + "." + timestamp + "." + base64Body
	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
