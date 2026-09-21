package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// =============================================================================
// HTTP Client for External API Calls
// =============================================================================

// Client is a reusable HTTP client for external API calls.
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// Options configures the HTTP client.
type Options struct {
	Timeout time.Duration
	Headers map[string]string
}

// New creates a new HTTP client with the given base URL and options.
func New(baseURL string, opts Options) *Client {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		headers: opts.Headers,
	}
}

// Get performs a GET request and decodes the response into dest.
func (c *Client) Get(ctx context.Context, path string, dest any) error {
	return c.do(ctx, http.MethodGet, path, nil, dest)
}

// Post performs a POST request with a JSON body and decodes the response into dest.
func (c *Client) Post(ctx context.Context, path string, body any, dest any) error {
	return c.do(ctx, http.MethodPost, path, body, dest)
}

// Put performs a PUT request with a JSON body and decodes the response into dest.
func (c *Client) Put(ctx context.Context, path string, body any, dest any) error {
	return c.do(ctx, http.MethodPut, path, body, dest)
}

// Delete performs a DELETE request and decodes the response into dest.
func (c *Client) Delete(ctx context.Context, path string, dest any) error {
	return c.do(ctx, http.MethodDelete, path, nil, dest)
}

// do executes the HTTP request.
func (c *Client) do(ctx context.Context, method, path string, body any, dest any) error {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("httpclient: failed to marshal body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("httpclient: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for key, val := range c.headers {
		req.Header.Set(key, val)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpclient: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("httpclient: failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
			URL:        url,
			Method:     method,
		}
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("httpclient: failed to decode response: %w", err)
		}
	}

	return nil
}

// APIError represents a non-2xx response from an external API.
type APIError struct {
	StatusCode int
	Body       string
	URL        string
	Method     string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("httpclient: %s %s returned %d: %s", e.Method, e.URL, e.StatusCode, e.Body)
}
