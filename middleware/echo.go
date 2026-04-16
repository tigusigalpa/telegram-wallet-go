//go:build echo
// +build echo

package middleware

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/tigusigalpa/telegram-wallet-go"
)

// EchoWebhookMiddleware creates an Echo middleware for Wallet Pay webhook verification.
// Usage:
//
//	e.POST("/webhook/walletpay", handleWebhook, middleware.EchoWebhookMiddleware(client))
//
//	func handleWebhook(c echo.Context) error {
//	    events := c.Get("walletpay_events").([]walletpay.WebhookEvent)
//	    // Process events...
//	    return c.String(200, "OK")
//	}
func EchoWebhookMiddleware(client *walletpay.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timestamp := c.Request().Header.Get("WalletPay-Timestamp")
			signature := c.Request().Header.Get("WalletPay-Signature")

			if timestamp == "" || signature == "" {
				return c.JSON(401, map[string]string{"error": "missing webhook headers"})
			}

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.JSON(400, map[string]string{"error": "cannot read body"})
			}

			events, err := client.VerifyAndParseWebhook(c.Request().Method, c.Request().URL.Path, timestamp, body, signature)
			if err != nil {
				return c.JSON(401, map[string]string{"error": "invalid signature"})
			}

			c.Set("walletpay_events", events)
			return next(c)
		}
	}
}
