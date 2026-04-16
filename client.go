package walletpay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Wpay-Store-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func (c *Client) handleErrorResponse(statusCode int, message string) error {
	switch statusCode {
	case 400:
		return &RequestError{Code: 400, Message: message, StatusCode: statusCode}
	case 401:
		return &AuthError{Message: message, StatusCode: statusCode}
	case 404:
		return &NotFoundError{Message: message, StatusCode: statusCode}
	case 429:
		return &RateLimitError{Message: message, StatusCode: statusCode}
	case 500:
		return &ServerError{Message: message, StatusCode: statusCode}
	default:
		return &APIError{Message: message, StatusCode: statusCode}
	}
}

func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiResp apiResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return c.handleErrorResponse(resp.StatusCode, "unknown error")
		}
		return c.handleErrorResponse(resp.StatusCode, apiResp.Message)
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
