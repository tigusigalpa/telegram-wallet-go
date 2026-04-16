package walletpay

import "time"

// MoneyAmount represents a monetary amount with currency.
type MoneyAmount struct {
	CurrencyCode string `json:"currencyCode"`
	Amount       string `json:"amount"`
}

// CreateOrderRequest is the payload for creating a new order.
type CreateOrderRequest struct {
	Amount                 MoneyAmount `json:"amount"`
	Description            string      `json:"description"`
	ExternalID             string      `json:"externalId"`
	TimeoutSeconds         int64       `json:"timeoutSeconds"`
	CustomerTelegramUserID int64       `json:"customerTelegramUserId"`
	AutoConversionCurrency string      `json:"autoConversionCurrency,omitempty"`
	ReturnURL              string      `json:"returnUrl,omitempty"`
	FailReturnURL          string      `json:"failReturnUrl,omitempty"`
	CustomData             string      `json:"customData,omitempty"`
}

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusActive    OrderStatus = "ACTIVE"
	OrderStatusExpired   OrderStatus = "EXPIRED"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderPreview contains the details of an order.
type OrderPreview struct {
	ID                     int64       `json:"id"`
	Status                 OrderStatus `json:"status"`
	Number                 string      `json:"number"`
	Amount                 MoneyAmount `json:"amount"`
	AutoConversionCurrency string      `json:"autoConversionCurrency,omitempty"`
	CreatedDateTime        time.Time   `json:"createdDateTime"`
	ExpirationDateTime     time.Time   `json:"expirationDateTime"`
	CompletedDateTime      *time.Time  `json:"completedDateTime,omitempty"`
	DirectPayLink          string      `json:"directPayLink"`
	PayLink                string      `json:"payLink,omitempty"`
}

// WebhookEventType represents the type of a webhook event.
type WebhookEventType string

const (
	WebhookEventOrderPaid   WebhookEventType = "ORDER_PAID"
	WebhookEventOrderFailed WebhookEventType = "ORDER_FAILED"
)

// WebhookEvent is a single event delivered by Wallet Pay.
type WebhookEvent struct {
	EventDateTime time.Time        `json:"eventDateTime"`
	EventID       int64            `json:"eventId"`
	Type          WebhookEventType `json:"type"`
	Payload       WebhookPayload   `json:"payload"`
}

// WebhookPayload contains the order data from a webhook event.
type WebhookPayload struct {
	ID                     int64           `json:"id"`
	Number                 string          `json:"number"`
	CustomData             string          `json:"customData,omitempty"`
	ExternalID             string          `json:"externalId"`
	OrderAmount            MoneyAmount     `json:"orderAmount"`
	SelectedPaymentOption  *PaymentOption  `json:"selectedPaymentOption,omitempty"`
	OrderCompletedDateTime *time.Time      `json:"orderCompletedDateTime,omitempty"`
	Status                 *OrderStatus    `json:"status,omitempty"`
}

// PaymentOption describes how the payer chose to pay.
type PaymentOption struct {
	Amount       MoneyAmount `json:"amount"`
	AmountFee    MoneyAmount `json:"amountFee"`
	AmountNet    MoneyAmount `json:"amountNet"`
	ExchangeRate string      `json:"exchangeRate"`
}

// apiResponse is the generic API response wrapper.
type apiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// orderResponse is the response for order creation and preview.
type orderResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    OrderPreview `json:"data"`
}

// orderListResponse is the response for order list.
type orderListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Items []OrderPreview `json:"items"`
	} `json:"data"`
}

// orderAmountResponse is the response for order count.
type orderAmountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TotalAmount int64 `json:"totalAmount"`
	} `json:"data"`
}
