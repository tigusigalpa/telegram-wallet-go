# Telegram Wallet Pay Go SDK

![Telegram Wallet Go SDK](https://i.postimg.cc/FsYbYdhJ/telegram-wallet-go-banner.jpg)

[![Go Version](https://img.shields.io/github/go-mod/go-version/tigusigalpa/telegram-wallet-go)](https://github.com/tigusigalpa/telegram-wallet-go)
[![License](https://img.shields.io/github/license/tigusigalpa/telegram-wallet-go)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/telegram-wallet-go)](https://goreportcard.com/report/github.com/tigusigalpa/telegram-wallet-go)

Accept crypto payments in your Telegram bot with just a few lines of Go code. This SDK wraps
the [Telegram Wallet Pay](https://pay.wallet.tg/) API, letting your users pay with TON, USDT, BTC, and NOT — right
inside Telegram.

## Why This SDK?

If you're building a Telegram bot and want to accept crypto payments, you've come to the right place. We've done the
heavy lifting so you can focus on your product:

- **Get started in minutes** — Simple, clean API that feels natural in Go
- **Battle-tested security** — Webhook signatures verified with HMAC-SHA256
- **Works with your stack** — Native support for net/http, Gin, and Echo
- **No bloat** — Zero external dependencies for the core client
- **Production-ready** — Comprehensive tests and proper error handling

## Installation

```bash
go get github.com/tigusigalpa/telegram-wallet-go
```

### Requirements

- Go 1.21 or higher

## Quick Start

Here's how simple it is to create a payment:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/tigusigalpa/telegram-wallet-go"
)

func main() {
    // Initialize the client with your API key from Wallet Pay
    client := walletpay.NewClient("YOUR_STORE_API_KEY")

    // Create a payment order
    order, err := client.CreateOrder(context.Background(), walletpay.CreateOrderRequest{
        Amount: walletpay.MoneyAmount{
            CurrencyCode: "USD",
            Amount:       "9.99",
        },
        Description:            "Premium subscription for 1 month",
        ExternalID:             "ORDER-12345",  // Your unique order ID
        TimeoutSeconds:         3600,            // 1 hour to pay
        CustomerTelegramUserID: 123456789,       // Who can pay this order
        AutoConversionCurrency: "USDT",          // Receive payment in USDT
        ReturnURL:              "https://t.me/YourBot/YourApp",
        CustomData:             `{"user_id":42}`,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Send this link to your user — they'll pay right in Telegram!
    fmt.Println("Payment URL:", order.DirectPayLink)
}
```

### Handle Webhooks

When a payment succeeds (or fails), Wallet Pay will notify your server. Here's how to handle it:

```go
package main

import (
    "log"
    "net/http"

    "github.com/tigusigalpa/telegram-wallet-go"
    "github.com/tigusigalpa/telegram-wallet-go/middleware"
)

func main() {
    client := walletpay.NewClient("YOUR_STORE_API_KEY")

    // The middleware handles signature verification for you
    http.HandleFunc("/webhook/walletpay", middleware.WalletPayWebhookHandler(
        client,
        func(w http.ResponseWriter, r *http.Request, events []walletpay.WebhookEvent) {
            for _, event := range events {
                if event.Type == walletpay.WebhookEventOrderPaid {
                    // 🎉 Payment successful!
                    log.Printf("Order %d paid: %s", event.Payload.ID, event.Payload.ExternalID)
                    
                    // Now you can:
                    // - Update your database
                    // - Grant premium access
                    // - Send a thank-you message
                }
            }
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("OK"))
        },
    ))

    log.Println("Webhook server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Configuration

The default settings work great for most cases, but you can customize everything:

```go
import (
    "net/http"
    "time"

    "github.com/tigusigalpa/telegram-wallet-go"
)

client := walletpay.NewClient(
    "YOUR_STORE_API_KEY",
    walletpay.WithBaseURL("https://pay.wallet.tg"),
    walletpay.WithTimeout(60 * time.Second),
    walletpay.WithHTTPClient(&http.Client{
        Timeout: 90 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
        },
    }),
)
```

### Options at a Glance

| Option                   | What it does            | Default                 |
|--------------------------|-------------------------|-------------------------|
| `WithBaseURL(url)`       | Set custom base URL     | `https://pay.wallet.tg` |
| `WithTimeout(duration)`  | Set HTTP client timeout | `30s`                   |
| `WithHTTPClient(client)` | Use custom HTTP client  | Standard `http.Client`  |

## API Reference

Here's everything you can do with the client:

| Method                    | Parameters                                 | Returns                 | Description                               |
|---------------------------|--------------------------------------------|-------------------------|-------------------------------------------|
| `CreateOrder()`           | `ctx, CreateOrderRequest`                  | `*OrderPreview, error`  | Create a new payment order                |
| `GetOrderPreview()`       | `ctx, orderID string`                      | `*OrderPreview, error`  | Get order details by ID                   |
| `GetOrderList()`          | `ctx, offset int64, count int32`           | `[]OrderPreview, error` | Get paginated list of orders (max 10,000) |
| `GetOrderAmount()`        | `ctx`                                      | `int64, error`          | Get total count of all orders             |
| `VerifyWebhook()`         | `method, path, timestamp, body, signature` | `error`                 | Verify webhook signature                  |
| `VerifyAndParseWebhook()` | `method, path, timestamp, body, signature` | `[]WebhookEvent, error` | Verify and parse webhook                  |

### Order Request Fields

| Field                    | Type          | Required | What it's for                                   |
|--------------------------|---------------|----------|-------------------------------------------------|
| `Amount`                 | `MoneyAmount` | Yes      | How much to charge (e.g., `{"USD", "9.99"}`)    |
| `Description`            | `string`      | Yes      | What the user sees (5-100 chars)                |
| `ExternalID`             | `string`      | Yes      | Your order ID — use this to match payments      |
| `TimeoutSeconds`         | `int64`       | Yes      | How long the order stays valid (30s to 10 days) |
| `CustomerTelegramUserID` | `int64`       | Yes      | Only this Telegram user can pay                 |
| `AutoConversionCurrency` | `string`      | No       | Convert payment to TON/USDT/BTC/NOT (+1% fee)   |
| `ReturnURL`              | `string`      | No       | Where to send user after payment                |
| `FailReturnURL`          | `string`      | No       | Where to send user if payment fails             |
| `CustomData`             | `string`      | No       | Your metadata — comes back in webhooks          |

### Currencies You Can Use

**Fiat (for pricing):** USD, EUR  
**Crypto (for receiving):** TON, USDT, BTC, NOT

### Order Lifecycle

Orders go through these states:

| Status      | Meaning                  |
|-------------|--------------------------|
| `ACTIVE`    | Waiting for payment      |
| `PAID`      | Payment received! 🎉     |
| `EXPIRED`   | Time ran out             |
| `CANCELLED` | User or system cancelled |

### Webhook Events

You'll receive one of these:

- `ORDER_PAID` — Money's in! Time to deliver.
- `ORDER_FAILED` — Something went wrong (expired, cancelled, etc.)

## Setting Up Webhooks

Webhooks tell you when payments happen. Here's how to set them up properly.

### Step 1: Configure Your URL

In your Wallet Pay store settings, set the webhook URL to something like:

```
https://yourdomain.com/webhook/walletpay
```

**Important:**

- Must be HTTPS with a real SSL certificate (Let's Encrypt works great)
- Self-signed certs won't work
- Always return HTTP 200 to confirm receipt

### Step 2: Allowlist Wallet Pay IPs

If you have a firewall, allow these IPs:

- `188.42.38.156`
- `172.255.249.124`

### Step 3: Handle Duplicate Webhooks

Wallet Pay might send the same webhook multiple times (network issues happen). Use `EventID` to avoid processing
duplicates:

```go
var processedEvents sync.Map // or use Redis, database, etc.

func handleWebhook(events []walletpay.WebhookEvent) {
    for _, event := range events {
        eventKey := fmt.Sprintf("webhook:%d", event.EventID)
        
        if _, exists := processedEvents.LoadOrStore(eventKey, true); exists {
            continue // Already processed
        }
        
        // Process the event
        processPayment(event)
    }
}
```

### How Signature Verification Works

You don't need to implement this yourself (our middleware handles it), but here's what happens under the hood:

```
stringToSign = HTTP_METHOD + "." + URI_PATH + "." + TIMESTAMP + "." + Base64(BODY)
signature    = Base64(HmacSHA256(stringToSign, API_KEY))
```

> **Heads up:** The URI path must match exactly what you configured — including the trailing slash (or lack thereof).

## Using With Your Framework

### Standard Library (net/http)

```go
import (
    "github.com/tigusigalpa/telegram-wallet-go"
    "github.com/tigusigalpa/telegram-wallet-go/middleware"
)

client := walletpay.NewClient("YOUR_API_KEY")

http.HandleFunc("/webhook/walletpay", middleware.WalletPayWebhookHandler(
    client,
    func(w http.ResponseWriter, r *http.Request, events []walletpay.WebhookEvent) {
        // Process events
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    },
))
```

### Gin

Using Gin? Build with the `gin` tag:

```bash
go build -tags gin
```

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/tigusigalpa/telegram-wallet-go"
    "github.com/tigusigalpa/telegram-wallet-go/middleware"
)

router := gin.Default()
client := walletpay.NewClient("YOUR_API_KEY")

router.POST("/webhook/walletpay", middleware.GinWebhookMiddleware(client), func(c *gin.Context) {
    events, _ := c.Get("walletpay_events")
    webhookEvents := events.([]walletpay.WebhookEvent)
    
    for _, event := range webhookEvents {
        // Process event
    }
    
    c.String(200, "OK")
})
```

### Echo

Prefer Echo? Build with the `echo` tag:

```bash
go build -tags echo
```

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/tigusigalpa/telegram-wallet-go"
    "github.com/tigusigalpa/telegram-wallet-go/middleware"
)

e := echo.New()
client := walletpay.NewClient("YOUR_API_KEY")

e.POST("/webhook/walletpay", func(c echo.Context) error {
    events := c.Get("walletpay_events").([]walletpay.WebhookEvent)
    
    for _, event := range events {
        // Process event
    }
    
    return c.String(200, "OK")
}, middleware.EchoWebhookMiddleware(client))
```

## When Things Go Wrong

Errors are typed, so you can handle them gracefully:

| Error Type            | What happened                       |
|-----------------------|-------------------------------------|
| `*RequestError`       | Bad request (check your parameters) |
| `*AuthError`          | Invalid API key                     |
| `*NotFoundError`      | Order doesn't exist                 |
| `*RateLimitError`     | Slow down! Too many requests        |
| `*ServerError`        | Wallet Pay is having issues         |
| `ErrInvalidSignature` | Webhook signature doesn't match     |

### Example

```go
import (
    "errors"
    "log"

    "github.com/tigusigalpa/telegram-wallet-go"
)

order, err := client.GetOrderPreview(ctx, "123456")
if err != nil {
    var notFoundErr *walletpay.NotFoundError
    var rateLimitErr *walletpay.RateLimitError
    
    switch {
    case errors.As(err, &notFoundErr):
        log.Println("Order not found:", notFoundErr.Message)
    case errors.As(err, &rateLimitErr):
        log.Println("Rate limit exceeded. Retry later.")
    default:
        log.Println("Error:", err)
    }
    return
}
```

## Real-World Examples

### Full-Featured Order Creation

```go
order, err := client.CreateOrder(ctx, walletpay.CreateOrderRequest{
    Amount: walletpay.MoneyAmount{
        CurrencyCode: "USD",
        Amount:       "49.99",
    },
    Description:            "Annual Pro subscription",
    ExternalID:             fmt.Sprintf("SUB-%d-%d", userID, time.Now().Unix()),
    TimeoutSeconds:         7200, // 2 hours
    CustomerTelegramUserID: userTelegramID,
    AutoConversionCurrency: "USDT",
    ReturnURL:              "https://t.me/YourBot/YourApp?success=1",
    FailReturnURL:          "https://t.me/YourBot/YourApp?failed=1",
    CustomData:             fmt.Sprintf(`{"user_id":%d,"plan":"pro","period":"annual"}`, userID),
})
```

### Check Order Status

```go
order, err := client.GetOrderPreview(ctx, "2703383946854401")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Order Status: %s\n", order.Status)
if order.CompletedDateTime != nil {
    fmt.Printf("Completed at: %s\n", order.CompletedDateTime.Format(time.RFC3339))
}
```

### Fetch All Orders (with Pagination)

```go
const pageSize = 100
offset := int64(0)

for {
    orders, err := client.GetOrderList(ctx, offset, pageSize)
    if err != nil {
        log.Fatal(err)
    }
    
    if len(orders) == 0 {
        break
    }
    
    for _, order := range orders {
        fmt.Printf("Order %d: %s - %s %s\n", 
            order.ID, order.Status, order.Amount.Amount, order.Amount.CurrencyCode)
    }
    
    offset += int64(len(orders))
}
```

### Production-Ready Webhook Handler

Here's a more complete example with idempotency and proper error handling:

```go
func handleWebhook(w http.ResponseWriter, r *http.Request, events []walletpay.WebhookEvent) {
    for _, event := range events {
        // Skip if we've already processed this event
        if isProcessed(event.EventID) {
            continue
        }
        
        switch event.Type {
        case walletpay.WebhookEventOrderPaid:
            if err := handleSuccessfulPayment(event); err != nil {
                // Log the error, but still acknowledge the webhook
                // (otherwise Wallet Pay will keep retrying)
                log.Printf("Error processing payment: %v", err)
            }
            
        case walletpay.WebhookEventOrderFailed:
            handleFailedPayment(event)
        }
        
        // Remember we've processed this
        markAsProcessed(event.EventID)
    }
    
    // Always return 200 to acknowledge receipt
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func handleSuccessfulPayment(event walletpay.WebhookEvent) error {
    payload := event.Payload
    
    // Get your custom data back
    var customData map[string]interface{}
    if err := json.Unmarshal([]byte(payload.CustomData), &customData); err != nil {
        return fmt.Errorf("invalid custom data: %w", err)
    }
    
    userID := int64(customData["user_id"].(float64))
    
    // Your business logic here:
    // 1. Mark order as paid in your database
    // 2. Grant premium access
    // 3. Send confirmation to user
    
    return nil
}
```

## Things to Know

### Opening the Payment Link

The payment link (`DirectPayLink`) needs to be opened correctly:

✅ **In a Telegram Web App:** Use `Telegram.WebApp.openTelegramLink(url)`  
✅ **In a bot message:** Use it as an Inline Button URL  
❌ **Don't use:** `openLink()` or `MenuButtonWebApp` — payment will fail

### Payment Button Text

Telegram requires specific button text:

- `👛 Wallet Pay`
- `👛 Pay via Wallet`

Yes, the purse emoji is mandatory. 👛

### Preventing Duplicate Orders

Use `ExternalID` as your idempotency key. If you retry with the same ID, you'll get the existing order back instead of
creating a duplicate:

```go
externalID := fmt.Sprintf("ORDER-%d-%d-%d", userID, productID, time.Now().Unix())
```

### Auto-Conversion Fees

Want to receive payments in a specific crypto? Set `AutoConversionCurrency`, but note:

- **1% fee** applies
- Minimum: **$1.30** (or $3 for BTC)

### When Can You Withdraw?

Funds are held for **48 hours** after payment before you can withdraw them. This is a Wallet Pay policy.

### One User Per Order

Only the Telegram user specified in `CustomerTelegramUserID` can pay that order. This prevents payment link sharing.

## Testing

```bash
# Run all tests
go test -v ./...

# With coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Try the Examples

We've included working examples you can run right away:

```bash
export WALLETPAY_API_KEY=your_api_key_here

# Create a test order
go run examples/create_order/main.go

# Start a webhook server
go run examples/webhook_server/main.go
```

Check out the [`examples/`](examples/) directory for the full code.

## What's Next?

This SDK is built to grow. The architecture separates payment functionality from future trading features (spot trading,
tokenized stocks, perpetual futures). When Wallet adds new APIs, we'll add support without breaking your existing code.

## Contributing

Found a bug? Have an idea? PRs are welcome!

1. Fork it
2. Create your branch (`git checkout -b fix/something`)
3. Make your changes
4. Run tests (`go test ./...`) and format (`go fmt ./...`)
5. Open a PR

Please include tests for new features.

## License

MIT — do whatever you want with it.

## Links

- [GitHub](https://github.com/tigusigalpa/telegram-wallet-go)
- [pkg.go.dev](https://pkg.go.dev/github.com/tigusigalpa/telegram-wallet-go)
- [Wallet Pay Docs](https://docs.wallet.tg/pay/)

## Need Help?

Open an [issue on GitHub](https://github.com/tigusigalpa/telegram-wallet-go/issues) — I'll do my best to help.

---

Built by [Igor Sazonov](https://github.com/tigusigalpa)
