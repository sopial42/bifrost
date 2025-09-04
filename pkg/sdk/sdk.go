package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

	errResponse := struct {
		Error *appErrors.AppError `json:"error"`
	}{}
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		return nil, appErrors.NewUnexpected("failed to handle error response", err)
	}

	appErr := errResponse.Error
	if errors.Is(appErr, appErrors.ErrUnknown) {
		appErr.Code = appErrors.CodeErrUnexpected
		appErr.Message = string(respBody)
	}

	return nil, appErr
}

// Patch performs a PATCH request and handles error responses
func (c *Client) Patch(ctx context.Context, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("PATCH", c.baseURL+url, bytes.NewReader(body))
	if err != nil {
		return nil, appErrors.NewUnexpected("sdk unable to create PATCH request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, appErrors.NewUnexpected("sdk unable to PATCH request", err)
	}

	return c.handleResponse(ctx, res)
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
