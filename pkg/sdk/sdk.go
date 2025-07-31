package sdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sopial42/bifrost/pkg/errors"
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
func (c *Client) handleResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.NewUnexpected("failed to read response body", err)
	}

	// If status code indicates success, return the body
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return respBody, nil
	}

	var errResp errors.ErrResponse
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		// If we can't parse the error response, use the raw body as message
		errResp.Error.Message = string(respBody)
		return nil, errors.NewUnexpected("failed to unmarshal error response", err)
	}

	return nil, errors.AppError{
		Code:    errResp.Error.AppCode,
		Message: errResp.Error.Message,
	}
}

// Post performs a POST request and handles error responses
func (c *Client) Post(url string, body []byte) ([]byte, error) {
	res, err := c.httpClient.Post(c.baseURL+url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, errors.NewUnexpected("sdk unable to POST request", err)
	}

	return c.handleResponse(res)
}

func (c *Client) Get(url string) ([]byte, error) {
	res, err := c.httpClient.Get(c.baseURL + url)
	if err != nil {
		return nil, errors.NewUnexpected("sdk unable to GET request", err)
	}

	return c.handleResponse(res)
}
