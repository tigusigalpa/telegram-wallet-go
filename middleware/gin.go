//go:build gin
// +build gin

package middleware

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/tigusigalpa/telegram-wallet-go"
)

// GinWebhookMiddleware creates a Gin middleware for Wallet Pay webhook verification.
// Usage:
//
//	router.POST("/webhook/walletpay", middleware.GinWebhookMiddleware(client), func(c *gin.Context) {
//	    events, _ := c.Get("walletpay_events")
//	    webhookEvents := events.([]walletpay.WebhookEvent)
//	    // Process events...
//	    c.String(200, "OK")
//	})
func GinWebhookMiddleware(client *walletpay.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp := c.GetHeader("WalletPay-Timestamp")
		signature := c.GetHeader("WalletPay-Signature")

		if timestamp == "" || signature == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing webhook headers"})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "cannot read body"})
			return
		}

		if err := client.VerifyWebhook(c.Request.Method, c.Request.URL.Path, timestamp, body, signature); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid signature"})
			return
		}
		events, err := walletpay.ParseWebhookEvents(body)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid payload"})
			return
		}

		c.Set("walletpay_events", events)
		c.Next()
	}
}
