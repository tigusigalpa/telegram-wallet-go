package walletpay

import (
	"context"
	"fmt"
)

// GetOrderList returns a paginated list of orders sorted by creation time (ascending).
// offset: number of orders to skip (pagination)
// count: number of orders to return (0-10000)
func (c *Client) GetOrderList(ctx context.Context, offset int64, count int32) ([]OrderPreview, error) {
	path := fmt.Sprintf("/wpay/store-api/v1/reconciliation/order-list?offset=%d&count=%d", offset, count)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result orderListResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data.Items, nil
}

// GetOrderAmount returns the total count of all orders in the store.
func (c *Client) GetOrderAmount(ctx context.Context) (int64, error) {
	resp, err := c.doRequest(ctx, "GET", "/wpay/store-api/v1/reconciliation/order-amount", nil)
	if err != nil {
		return 0, err
	}

	var result orderAmountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Data.TotalAmount, nil
}
