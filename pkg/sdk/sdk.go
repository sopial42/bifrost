package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	appErrors "github.com/sopial42/bifrost/pkg/errors"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewSDKClient(baseURL string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// handleResponse reads the response body and converts any error responses to AppErrors
func (c *Client) handleResponse(ctx context.Context, res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, appErrors.NewUnexpected("failed to read response body", err)
	}

	// If status code indicates success, return the body
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return respBody, nil
	}

	var appErr appErrors.AppError
	if err := json.Unmarshal(respBody, &appErr); err != nil {
		return nil, appErrors.NewUnexpected("failed to handle error response", err)
	}

	if appErr.Code == appErrors.ErrUnknown {
		appErr.Code = appErrors.ErrUnexpected
		appErr.Message = string(respBody)
	}

	return nil, appErr
}

// Post performs a POST request and handles error responses
func (c *Client) Post(ctx context.Context, url string, body []byte) ([]byte, error) {
	res, err := c.httpClient.Post(c.baseURL+url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, appErrors.NewUnexpected("sdk unable to POST request", err)
	}

	return c.handleResponse(ctx, res)
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	res, err := c.httpClient.Get(c.baseURL + url)
	if err != nil {
		return nil, appErrors.NewUnexpected("sdk unable to GET request", err)
	}

	return c.handleResponse(ctx, res)
}
