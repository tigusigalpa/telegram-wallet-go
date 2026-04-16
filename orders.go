package walletpay

import (
	"context"
	"fmt"
)

// CreateOrder creates a new payment order.
func (c *Client) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderPreview, error) {
	resp, err := c.doRequest(ctx, "POST", "/wpay/store-api/v1/order", req)
	if err != nil {
		return nil, err
	}

	var result orderResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetOrderPreview retrieves the current state of an order by its ID.
func (c *Client) GetOrderPreview(ctx context.Context, orderID string) (*OrderPreview, error) {
	path := fmt.Sprintf("/wpay/store-api/v1/order/preview?id=%s", orderID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result orderResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
