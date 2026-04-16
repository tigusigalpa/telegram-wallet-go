package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tigusigalpa/telegram-wallet-go"
)

func main() {
	apiKey := os.Getenv("WALLETPAY_API_KEY")
	if apiKey == "" {
		log.Fatal("WALLETPAY_API_KEY environment variable is required")
	}

	client := walletpay.NewClient(apiKey)

	// Create a payment order
	order, err := client.CreateOrder(context.Background(), walletpay.CreateOrderRequest{
		Amount: walletpay.MoneyAmount{
			CurrencyCode: "USD",
			Amount:       "9.99",
		},
		Description:            "Premium subscription for 1 month",
		ExternalID:             fmt.Sprintf("ORDER-%d", os.Getpid()),
		TimeoutSeconds:         3600,
		CustomerTelegramUserID: 123456789,
		AutoConversionCurrency: "USDT",
		ReturnURL:              "https://t.me/YourBot/YourApp",
		CustomData:             `{"user_id":42,"plan":"premium"}`,
	})
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	fmt.Println("Order created successfully!")
	fmt.Printf("Order ID: %d\n", order.ID)
	fmt.Printf("Order Number: %s\n", order.Number)
	fmt.Printf("Status: %s\n", order.Status)
	fmt.Printf("Amount: %s %s\n", order.Amount.Amount, order.Amount.CurrencyCode)
	fmt.Printf("Payment URL: %s\n", order.DirectPayLink)
	fmt.Printf("Expires at: %s\n", order.ExpirationDateTime.Format("2006-01-02 15:04:05"))

	// Pretty print the full order
	orderJSON, _ := json.MarshalIndent(order, "", "  ")
	fmt.Printf("\nFull order details:\n%s\n", orderJSON)
}
