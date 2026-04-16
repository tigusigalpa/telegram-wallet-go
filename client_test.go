package walletpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/wpay/store-api/v1/order", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("Wpay-Store-Api-Key"))

		response := orderResponse{
			Status:  "SUCCESS",
			Message: "",
			Data: OrderPreview{
				ID:                 2703383946854401,
				Status:             OrderStatusActive,
				Number:             "9aeb581c",
				Amount:             MoneyAmount{CurrencyCode: "USD", Amount: "1.00"},
				CreatedDateTime:    time.Now(),
				ExpirationDateTime: time.Now().Add(time.Hour),
				DirectPayLink:      "https://t.me/wallet/start?startapp=wpay_order-orderId__2703383946854401",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	req := CreateOrderRequest{
		Amount:                 MoneyAmount{CurrencyCode: "USD", Amount: "1.00"},
		Description:            "Test order",
		ExternalID:             "TEST-001",
		TimeoutSeconds:         3600,
		CustomerTelegramUserID: 123456789,
	}

	order, err := client.CreateOrder(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, int64(2703383946854401), order.ID)
	assert.Equal(t, OrderStatusActive, order.Status)
	assert.Equal(t, "9aeb581c", order.Number)
	assert.Equal(t, "USD", order.Amount.CurrencyCode)
	assert.Equal(t, "1.00", order.Amount.Amount)
}

func TestGetOrderPreview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/wpay/store-api/v1/order/preview")
		assert.Equal(t, "2703383946854401", r.URL.Query().Get("id"))

		completedTime := time.Now()
		response := orderResponse{
			Status:  "SUCCESS",
			Message: "",
			Data: OrderPreview{
				ID:                 2703383946854401,
				Status:             OrderStatusPaid,
				Number:             "9aeb581c",
				Amount:             MoneyAmount{CurrencyCode: "USD", Amount: "1.00"},
				CreatedDateTime:    time.Now().Add(-2 * time.Hour),
				ExpirationDateTime: time.Now().Add(-time.Hour),
				CompletedDateTime:  &completedTime,
				DirectPayLink:      "https://t.me/wallet/start?startapp=wpay_order-orderId__2703383946854401",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	order, err := client.GetOrderPreview(context.Background(), "2703383946854401")
	require.NoError(t, err)
	assert.Equal(t, int64(2703383946854401), order.ID)
	assert.Equal(t, OrderStatusPaid, order.Status)
	assert.NotNil(t, order.CompletedDateTime)
}

func TestGetOrderList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/wpay/store-api/v1/reconciliation/order-list")
		assert.Equal(t, "0", r.URL.Query().Get("offset"))
		assert.Equal(t, "10", r.URL.Query().Get("count"))

		response := orderListResponse{
			Status:  "SUCCESS",
			Message: "",
			Data: struct {
				Items []OrderPreview `json:"items"`
			}{
				Items: []OrderPreview{
					{
						ID:                 1,
						Status:             OrderStatusActive,
						Number:             "ABC123",
						Amount:             MoneyAmount{CurrencyCode: "USD", Amount: "10.00"},
						CreatedDateTime:    time.Now(),
						ExpirationDateTime: time.Now().Add(time.Hour),
						DirectPayLink:      "https://t.me/wallet/start?startapp=wpay_order-orderId__1",
					},
					{
						ID:                 2,
						Status:             OrderStatusPaid,
						Number:             "DEF456",
						Amount:             MoneyAmount{CurrencyCode: "EUR", Amount: "20.00"},
						CreatedDateTime:    time.Now(),
						ExpirationDateTime: time.Now().Add(time.Hour),
						DirectPayLink:      "https://t.me/wallet/start?startapp=wpay_order-orderId__2",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	orders, err := client.GetOrderList(context.Background(), 0, 10)
	require.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, int64(1), orders[0].ID)
	assert.Equal(t, int64(2), orders[1].ID)
	assert.Equal(t, OrderStatusActive, orders[0].Status)
	assert.Equal(t, OrderStatusPaid, orders[1].Status)
}

func TestGetOrderAmount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/wpay/store-api/v1/reconciliation/order-amount")

		response := orderAmountResponse{
			Status:  "SUCCESS",
			Message: "",
			Data: struct {
				TotalAmount int64 `json:"totalAmount"`
			}{
				TotalAmount: 42,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	count, err := client.GetOrderAmount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(42), count)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
		expectedType   interface{}
	}{
		{
			name:           "400 Bad Request",
			statusCode:     400,
			responseBody:   `{"status":"ERROR","message":"Invalid request"}`,
			expectedErrMsg: "Invalid request",
			expectedType:   &RequestError{},
		},
		{
			name:           "401 Unauthorized",
			statusCode:     401,
			responseBody:   `{"status":"ERROR","message":"Invalid API key"}`,
			expectedErrMsg: "Invalid API key",
			expectedType:   &AuthError{},
		},
		{
			name:           "404 Not Found",
			statusCode:     404,
			responseBody:   `{"status":"ERROR","message":"Order not found"}`,
			expectedErrMsg: "Order not found",
			expectedType:   &NotFoundError{},
		},
		{
			name:           "429 Rate Limit",
			statusCode:     429,
			responseBody:   `{"status":"ERROR","message":"Rate limit exceeded"}`,
			expectedErrMsg: "Rate limit exceeded",
			expectedType:   &RateLimitError{},
		},
		{
			name:           "500 Server Error",
			statusCode:     500,
			responseBody:   `{"status":"ERROR","message":"Internal server error"}`,
			expectedErrMsg: "Internal server error",
			expectedType:   &ServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewClient("test-api-key", WithBaseURL(server.URL))

			_, err := client.GetOrderPreview(context.Background(), "123")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
			assert.IsType(t, tt.expectedType, err)
		})
	}
}

func TestClientOptions(t *testing.T) {
	t.Run("WithBaseURL", func(t *testing.T) {
		client := NewClient("test-key", WithBaseURL("https://custom.url"))
		assert.Equal(t, "https://custom.url", client.baseURL)
	})

	t.Run("WithTimeout", func(t *testing.T) {
		client := NewClient("test-key", WithTimeout(60*time.Second))
		assert.Equal(t, 60*time.Second, client.httpClient.Timeout)
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		customClient := &http.Client{Timeout: 90 * time.Second}
		client := NewClient("test-key", WithHTTPClient(customClient))
		assert.Equal(t, customClient, client.httpClient)
	})
}
