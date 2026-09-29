package ntfy

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	retryhttp "github.com/hashicorp/go-retryablehttp"
)

const (
	_attempts     = 4               // attempts amount after first failed
	_minRetryWait = 2 * time.Second // min wait time between retries
	_timeout      = 5 * time.Second // request timeout
)

// SendMsg sends message with given params using NtfyClient connection.
func (c *Client) SendMsg(theme, title, priority string,
	tags []string, text string) ([]byte, error) {

	url := c.addr + theme
	req, err := http.NewRequest("POST", url, strings.NewReader(text))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Title", title)
	req.Header.Set("Priority", priority)
	req.Header.Set("Tags", strings.Join(tags, ","))
	req.Header.Set("Authorization", c.authHeader)

	// http auto-retry set up
	client := retryhttp.NewClient()
	client.HTTPClient = &http.Client{Timeout: _timeout}
	client.RetryWaitMin = _minRetryWait
	client.RetryMax = _attempts

	// wrap request for auto-retry
	retryReq, err := retryhttp.FromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("wrap request for retry: %w", err)
	}
	// send request
	resp, err := client.Do(retryReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 { //nolint:mnd // check 2xx code
		return nil, parseError(resp)
	}

	// read body
	bytesMsg, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return bytesMsg, nil
}

// parseError parses error from response.
func parseError(resp *http.Response) error {
	bytesMsg, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("parse error: read body: %w", err)
	}
	return fmt.Errorf("error: %s", string(bytesMsg))
}
